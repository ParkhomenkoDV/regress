package compare

import (
	"context"
	"fmt"
	"regress/internal/handlers/read"
	"regress/internal/progress"
	"regress/internal/shared"
	"regress/internal/storage"
	"regress/internal/utils"
	"sync"
	"sync/atomic"
	"time"
)

type task struct {
	fileName       string `doc:"Имя файла"`
	filePathBefore string `doc:"Путь до файла ДО"`
	filePathAfter  string `doc:"Путь до файла ПОСЛЕ"`
}

type result struct {
	shared.Comparison `doc:"Сравнение"`
	duration          time.Duration `doc:"Время обработки"`
	err               error         `doc:"Ошибка"`
}

func JSONs(ctx context.Context, filesBefore, filesAfter map[string]string, workers int) []shared.Comparison {
	startTime := time.Now()

	allFiles := make(map[string]task, len(filesBefore)+len(filesAfter))
	for name, dir := range filesBefore {
		allFiles[name] = task{fileName: name, filePathBefore: dir, filePathAfter: ""}
	}
	for name, dir := range filesAfter {
		if t, ok := allFiles[name]; ok {
			t.filePathAfter = dir
			allFiles[name] = t
		} else {
			allFiles[name] = task{fileName: name, filePathBefore: "", filePathAfter: dir}
		}
	}

	tasks := make(chan task, len(allFiles))     // канал задач
	results := make(chan result, len(allFiles)) // канал результатов

	// Заполняем очередь задач
	for _, t := range allFiles {
		tasks <- t
	}
	close(tasks) // Закрываем смену

	var wg sync.WaitGroup // Счётчик рабочих

	// Атомарная статистика
	var (
		total        uint64
		totalSuccess uint64
		totalErrors  uint64
	)

	bar := progress.New(uint64(len(allFiles)), time.Second, true, false, true)
	go bar.Show(ctx, &total, &totalSuccess, &totalErrors)

	// Запускаем рабочих
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go work(i, &wg,
			ctx,
			tasks, results,
			&total, &totalSuccess, &totalErrors)
	}

	// Ждём окончания смены
	go func() {
		wg.Wait()
		close(results)
	}()

	// Собираем результаты
	result := make([]shared.Comparison, 0, len(allFiles))
	successCount, errorCount := 0, 0
	for r := range results {
		result = append(result, r.Comparison)
		if r.err != nil {
			errorCount++
		} else {
			successCount++
		}
	}

	totalDuration := time.Since(startTime)

	fmt.Println()
	fmt.Printf("📊 ИТОГО: Успешно %d | Ошибок %d | Всего %d\n", successCount, errorCount, len(allFiles))
	fmt.Printf("⏱️  Общее время: %v\n", totalDuration.Round(time.Millisecond))
	fmt.Printf("⚡ Скорость: %.2f файлов/сек\n", float64(len(allFiles))/totalDuration.Seconds())
	fmt.Println()

	return result
}

func work(
	i int, wg *sync.WaitGroup,
	ctx context.Context,
	tasks <-chan task, results chan<- result,
	total, totalSuccess, totalErrors *uint64,
) {
	defer wg.Done()

	for t := range tasks {
		select {
		case <-ctx.Done():
			results <- result{
				Comparison: shared.Comparison{
					FileName: t.fileName,
				},
				err: ctx.Err(),
			}
			atomic.AddUint64(total, 1)
			atomic.AddUint64(totalErrors, 1)
			continue
		default:
		}

		startTime := time.Now()

		// Читаем файлы
		jsonBefore, err := read.JSON(t.filePathBefore)
		if err != nil {
			results <- result{
				Comparison: shared.Comparison{
					FileName: t.fileName,
				},
				duration: time.Since(startTime),
				err:      err,
			}
			atomic.AddUint64(total, 1)
			atomic.AddUint64(totalErrors, 1)
			continue
		}
		jsonAfter, err := read.JSON(t.filePathAfter)
		if err != nil {
			results <- result{
				Comparison: shared.Comparison{
					FileName: t.fileName,
				},
				duration: time.Since(startTime),
				err:      err,
			}
			atomic.AddUint64(total, 1)
			atomic.AddUint64(totalErrors, 1)
			continue
		}

		difference := findDifferences(jsonBefore, jsonAfter, "")

		results <- result{
			Comparison: shared.Comparison{
				FileName:     t.fileName,
				ExistsInBoth: t.filePathBefore != "" && t.filePathAfter != "",
				Differences:  difference,
			},
			duration: time.Since(startTime),
			err:      nil,
		}
		atomic.AddUint64(total, 1)
		atomic.AddUint64(totalSuccess, 1)
	}
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
