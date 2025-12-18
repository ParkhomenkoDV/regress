package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regress/internal/config"
	"regress/internal/storage"
	"regress/pkg/utils"
	"sync"
	"time"

	excel "github.com/xuri/excelize/v2"
)

// Comparison содержит сравнение одного файла
type Comparison struct {
	FileName     string
	Before       storage.DB
	After        storage.DB
	ExistsInBoth bool
	Differences  []Difference
}

// Difference описывает одно различие
type Difference struct {
	Field  string      `doc:"Поле"`
	Before interface{} `doc:"Значение до"`
	After  interface{} `doc:"Значение после"`
}

func main() {
	cfg, err := config.New()
	if err != nil {
		panic(err)
	}

	start := time.Now()
	comparisons, err := compareJSONs(cfg.BeforeDir, cfg.AfterDir, cfg.Workers)
	if err != nil {
		panic(err)
	}

	// Фильтруем если нужно показывать только изменения
	var comparison []Comparison
	for _, comp := range comparisons {
		if cfg.ShowAll || len(comp.Differences) > 0 || !comp.ExistsInBoth {
			comparison = append(comparison, comp)
		}
	}

	err = ExportToExcel(comparison, "comparison.xlsx")
	if err != nil {
		panic(err)
	}

	fmt.Printf("\nСтатистика:\n")
	fmt.Printf("\tВсего файлов: %d\n", len(comparisons))
	fmt.Printf("\tС изменениями: %d\n", countChanged(comparisons))
	fmt.Printf("\tВремя обработки: %v\n", time.Since(start))
}

func compareJSONs(beforeDir, afterDir string, workers int) ([]Comparison, error) {
	// Получаем списки файлов
	beforeFiles, err := utils.GetJSONFiles(beforeDir)
	if err != nil {
		return nil, err
	}
	afterFiles, err := utils.GetJSONFiles(afterDir)
	if err != nil {
		return nil, err
	}

	// Создаем мапу для быстрого поиска
	afterMap := make(map[string]bool)
	for _, f := range afterFiles {
		afterMap[f] = true
	}

	jobs := make(chan string, len(beforeFiles))
	results := make(chan Comparison, len(beforeFiles))
	errors := make(chan error, len(beforeFiles))

	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go worker(jobs, results, errors, beforeDir, afterDir, afterMap, &wg)
	}

	// Отправляем задания
	for _, file := range beforeFiles {
		jobs <- file
	}
	close(jobs)

	// Ждем завершения
	go func() {
		wg.Wait()
		close(results)
		close(errors)
	}()

	// Собираем результаты
	var comparisons []Comparison
	for result := range results {
		comparisons = append(comparisons, result)
	}

	// Проверяем ошибки
	select {
	case err := <-errors:
		if err != nil {
			return nil, err
		}
	default:
	}

	return comparisons, nil
}

func worker(jobs <-chan string, results chan<- Comparison, errors chan<- error,
	beforeDir, afterDir string, afterMap map[string]bool, wg *sync.WaitGroup) {
	defer wg.Done()

	for filename := range jobs {
		comp, err := compareFile(filename, beforeDir, afterDir, afterMap[filename])
		if err != nil {
			errors <- fmt.Errorf("файл %s: %v", filename, err)
			continue
		}
		results <- comp
	}
}

func compareFile(filename, beforeDir, afterDir string, existsInAfter bool) (Comparison, error) {
	var comparison Comparison
	comparison.FileName = filename
	comparison.ExistsInBoth = existsInAfter

	// Читаем файл "до"
	beforePath := filepath.Join(beforeDir, filename)
	beforeData, err := readDynamicJSON(beforePath)
	if err != nil {
		return comparison, err
	}
	comparison.Before = beforeData

	// Читаем файл "после" если существует
	if existsInAfter {
		afterPath := filepath.Join(afterDir, filename)
		afterData, err := readDynamicJSON(afterPath)
		if err != nil {
			return comparison, err
		}
		comparison.After = afterData

		// Находим различия
		comparison.Differences = findDifferences(beforeData, afterData, "")
	}

	return comparison, nil
}

func readDynamicJSON(path string) (storage.DB, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var db storage.DB
	if err := json.Unmarshal(data, &db); err != nil {
		return nil, err
	}

	return db, nil
}
func findDifferences(before, after storage.DB, prefix string) []Difference {
	var diffs []Difference

	// Все уникальные ключи
	allKeys := make(map[string]bool)
	for k := range before {
		allKeys[k] = true
	}
	for k := range after {
		allKeys[k] = true
	}

	for key := range allKeys {
		fullKey := key
		if prefix != "" {
			fullKey = prefix + "." + key
		}

		beforeVal, beforeExists := before[key]
		afterVal, afterExists := after[key]

		switch {
		case beforeExists && !afterExists:
			diffs = append(diffs, Difference{
				Field:  fullKey,
				Before: beforeVal,
				After:  nil,
			})

		case !beforeExists && afterExists:
			diffs = append(diffs, Difference{
				Field:  fullKey,
				Before: nil,
				After:  afterVal,
			})

		default:
			if utils.IsMap(beforeVal) && utils.IsMap(afterVal) {
				// Если оба значения - словари, сравниваем рекурсивно
				beforeMap, _ := beforeVal.(map[string]interface{})
				afterMap, _ := afterVal.(map[string]interface{})
				nestedDiffs := findDifferences(beforeMap, afterMap, fullKey)
				diffs = append(diffs, nestedDiffs...)
			} else if utils.IsSlice(beforeVal) && utils.IsSlice(afterVal) {
				// Если оба значения - срезы/массивы, сравниваем поэлементно
				beforeSlice, _ := beforeVal.([]interface{})
				afterSlice, _ := afterVal.([]interface{})

				// Сравниваем по максимальной длине
				maxLen := len(beforeSlice)
				if len(afterSlice) > maxLen {
					maxLen = len(afterSlice)
				}

				for i := 0; i < maxLen; i++ {
					elementKey := fmt.Sprintf("%s[%d]", fullKey, i)

					var beforeElem, afterElem interface{}
					var beforeElemExists, afterElemExists bool

					if i < len(beforeSlice) {
						beforeElem = beforeSlice[i]
						beforeElemExists = true
					}

					if i < len(afterSlice) {
						afterElem = afterSlice[i]
						afterElemExists = true
					}

					if beforeElemExists && !afterElemExists {
						diffs = append(diffs, Difference{
							Field:  elementKey,
							Before: beforeElem,
							After:  nil,
						})
					} else if !beforeElemExists && afterElemExists {
						diffs = append(diffs, Difference{
							Field:  elementKey,
							Before: nil,
							After:  afterElem,
						})
					} else if !utils.IsEqual(beforeElem, afterElem) {
						// Рекурсивно сравниваем элементы если они сложные
						if utils.IsMap(beforeElem) && utils.IsMap(afterElem) {
							beforeMap, _ := beforeElem.(map[string]interface{})
							afterMap, _ := afterElem.(map[string]interface{})
							nestedDiffs := findDifferences(beforeMap, afterMap, elementKey)
							diffs = append(diffs, nestedDiffs...)
						} else if utils.IsSlice(beforeElem) && utils.IsSlice(afterElem) {
							// Для вложенных срезов тоже рекурсивно сравниваем
							beforeMap := map[string]interface{}{"": beforeElem}
							afterMap := map[string]interface{}{"": afterElem}
							nestedDiffs := findDifferences(beforeMap, afterMap, elementKey)
							diffs = append(diffs, nestedDiffs...)
						} else {
							diffs = append(diffs, Difference{
								Field:  elementKey,
								Before: beforeElem,
								After:  afterElem,
							})
						}
					}
				}
			} else if !utils.IsEqual(beforeVal, afterVal) {
				// Для всех остальных типов используем isEqual
				diffs = append(diffs, Difference{
					Field:  fullKey,
					Before: beforeVal,
					After:  afterVal,
				})
			}
		}
	}

	return diffs
}
func countChanged(comparisons []Comparison) int {
	count := 0
	for _, comp := range comparisons {
		if comp.ExistsInBoth && len(comp.Differences) > 0 {
			count++
		}
	}
	return count
}

func ExportToExcel(comparisons []Comparison, filename string) error {
	f := excel.NewFile() // Создаем новый Excel файл

	// Получаем все уникальные названия полей для создания заголовков
	fieldSet := make(map[string]bool)
	for _, comp := range comparisons {
		for _, diff := range comp.Differences {
			fieldSet[diff.Field] = true
		}
	}

	// Преобразуем в упорядоченный список полей
	fields := make([]string, 0, len(fieldSet))
	for field := range fieldSet {
		fields = append(fields, field)
	}

	// Создаем заголовки
	headers := []string{"FileName", "ExistsInBoth"}
	for _, field := range fields {
		headers = append(headers, field+"_before", field+"_after")
	}

	// Записываем заголовки в первую строку
	for col, header := range headers {
		cell, _ := excel.CoordinatesToCellName(col+1, 1)
		f.SetCellValue("Sheet1", cell, header)
	}

	// Записываем данные
	for row, comp := range comparisons {
		// Преобразуем различия в карту для быстрого доступа
		diffMap := make(map[string]Difference)
		for _, diff := range comp.Differences {
			diffMap[diff.Field] = diff
		}

		// Заполняем ячейки
		col := 1

		// FileName
		cell, _ := excel.CoordinatesToCellName(col, row+2)
		f.SetCellValue("Sheet1", cell, comp.FileName)
		col++

		// ExistsInBoth
		cell, _ = excel.CoordinatesToCellName(col, row+2)
		f.SetCellValue("Sheet1", cell, comp.ExistsInBoth)
		col++

		// Данные для каждого поля
		for _, field := range fields {
			if diff, exists := diffMap[field]; exists {
				// Before значение
				cell, _ = excel.CoordinatesToCellName(col, row+2)
				f.SetCellValue("Sheet1", cell, diff.Before)
				col++

				// After значение
				cell, _ = excel.CoordinatesToCellName(col, row+2)
				f.SetCellValue("Sheet1", cell, diff.After)
				col++
			} else {
				// Если поля нет в различиях, оставляем пустые ячейки
				col += 2
			}
		}
	}
	// Настраиваем ширину колонок для лучшего отображения
	for i := 1; i <= len(headers); i++ {
		f.SetColWidth("Sheet1", string(rune('A'+i-1)), string(rune('A'+i-1)), 20)
	}

	// Записываем форматирование для заголовков
	style, _ := f.NewStyle(&excel.Style{
		Font: &excel.Font{
			Bold: true,
		},
		Fill: excel.Fill{
			Type:    "pattern",
			Color:   []string{"#E0E0E0"},
			Pattern: 1,
		},
	})

	// Применяем стиль к заголовкам
	for i := 1; i <= len(headers); i++ {
		cell, _ := excel.CoordinatesToCellName(i, 1)
		f.SetCellStyle("Sheet1", cell, cell, style)
	}

	// Сохраняем файл
	if err := f.SaveAs(filename); err != nil {
		return err
	}

	return nil
}
