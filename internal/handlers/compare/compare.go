package compare

import (
	"fmt"
	"path/filepath"
	"regress/internal/handlers/read"
	"regress/internal/shared"
	"regress/internal/storage"
	"regress/internal/utils"
	"sync"
	"time"
)

type result struct {
	fileName string
	duration time.Duration
	err      error
}

func JSONs(beforeDir, afterDir string, workers int) ([]shared.Comparison, error) {
	// Получаем списки файлов
	beforeFiles, err := read.JSONFiles(beforeDir)
	if err != nil {
		return nil, err
	}
	afterFiles, err := read.JSONFiles(afterDir)
	if err != nil {
		return nil, err
	}

	// Создаем мапу для быстрого поиска
	afterMap := make(map[string]bool)
	for _, f := range afterFiles {
		afterMap[f] = true
	}

	jobs := make(chan string, len(beforeFiles))
	results := make(chan shared.Comparison, len(beforeFiles))
	errors := make(chan error, len(beforeFiles))

	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go work(jobs, results, errors, beforeDir, afterDir, afterMap, &wg)
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
	var comparisons []shared.Comparison
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

func work(
	jobs <-chan string, results chan<- shared.Comparison, errors chan<- error,
	beforeDir, afterDir string,
	afterMap map[string]bool,
	wg *sync.WaitGroup,
) {
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

func compareFile(filename, beforeDir, afterDir string, existsInAfter bool) (shared.Comparison, error) {
	comparison := shared.Comparison{
		FileName:     filename,
		ExistsInBoth: existsInAfter,
	}

	// Читаем файл "до"
	beforePath := filepath.Join(beforeDir, filename)
	beforeData, err := read.JSON(beforePath)
	if err != nil {
		return comparison, err
	}
	comparison.Before = beforeData

	// Читаем файл "после" если существует
	if existsInAfter {
		afterPath := filepath.Join(afterDir, filename)
		afterData, err := read.JSON(afterPath)
		if err != nil {
			return comparison, err
		}
		comparison.After = afterData

		// Находим различия
		comparison.Differences = findDifferences(beforeData, afterData, "")
	}

	return comparison, nil
}

func findDifferences(before, after storage.DB, prefix string) []shared.Difference {
	diffs := make([]shared.Difference, 0, max(len(before), len(after)))

	// Все уникальные ключи
	allKeys := make(map[string]struct{})
	for k := range before {
		allKeys[k] = struct{}{}
	}
	for k := range after {
		allKeys[k] = struct{}{}
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
			diffs = append(diffs, shared.Difference{
				Field:  fullKey,
				Before: beforeVal,
				After:  nil,
			})

		case !beforeExists && afterExists:
			diffs = append(diffs, shared.Difference{
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
				// Если оба значения - срезы/массивы
				beforeSlice, _ := beforeVal.([]interface{})
				afterSlice, _ := afterVal.([]interface{})

				// Проверяем, содержит ли срез только простые типы
				if utils.IsSimpleSlice(beforeSlice) && utils.IsSimpleSlice(afterSlice) {
					// Если оба среза содержат только простые типы, сравниваем целиком
					if !utils.IsEqualSimpleSlices(beforeSlice, afterSlice) {
						diffs = append(diffs, shared.Difference{
							Field:  fullKey,
							Before: beforeSlice,
							After:  afterSlice,
						})
					}
				} else {
					// Если срезы содержат сложные типы, сравниваем поэлементно
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
							diffs = append(diffs, shared.Difference{
								Field:  elementKey,
								Before: beforeElem,
								After:  nil,
							})
						} else if !beforeElemExists && afterElemExists {
							diffs = append(diffs, shared.Difference{
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
								diffs = append(diffs, shared.Difference{
									Field:  elementKey,
									Before: beforeElem,
									After:  afterElem,
								})
							}
						}
					}
				}
			} else if !utils.IsEqual(beforeVal, afterVal) {
				// Для всех остальных типов используем isEqual
				diffs = append(diffs, shared.Difference{
					Field:  fullKey,
					Before: beforeVal,
					After:  afterVal,
				})
			}
		}
	}

	return diffs
}
