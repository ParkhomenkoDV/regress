package compare

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"regress/internal/storage"
	"testing"
)

// BenchmarkJSONs измеряет производительность функции JSONs с различными размерами данных и количеством воркеров.
func BenchmarkJSONs(b *testing.B) {
	// Подготовка тестовых данных
	testCases := []struct {
		name          string
		before, after storage.DB
	}{
		{
			name: "simple",
			before: storage.DB{
				"id":     1,
				"name":   "test",
				"price":  99.99,
				"active": true,
			},
			after: storage.DB{
				"id":     1,
				"name":   "test_modified",
				"price":  109.99,
				"active": true,
			},
		},
		{
			name: "nested",
			before: storage.DB{
				"user": map[string]interface{}{
					"name": "John",
					"age":  30,
					"address": map[string]interface{}{
						"city":    "Moscow",
						"street":  "Tverskaya",
						"zipcode": "123456",
					},
				},
				"items": []interface{}{
					map[string]interface{}{"id": 1, "name": "item1"},
					map[string]interface{}{"id": 2, "name": "item2"},
				},
			},
			after: storage.DB{
				"user": map[string]interface{}{
					"name": "John",
					"age":  31,
					"address": map[string]interface{}{
						"city":    "Moscow",
						"street":  "Tverskaya",
						"zipcode": "123457",
					},
				},
				"items": []interface{}{
					map[string]interface{}{"id": 1, "name": "item1"},
					map[string]interface{}{"id": 2, "name": "item2_modified"},
				},
			},
		},
		{
			name: "large_slice",
			before: storage.DB{
				"data": generateLargeSlice(1000),
			},
			after: storage.DB{
				"data": generateLargeSlice(1000),
			},
		},
		{
			name:   "identical",
			before: storage.DB{"key": "value", "nested": map[string]interface{}{"a": 1}},
			after:  storage.DB{"key": "value", "nested": map[string]interface{}{"a": 1}},
		},
		{
			name:   "empty",
			before: storage.DB{},
			after:  storage.DB{},
		},
	}

	// Количество воркеров для тестирования
	workerCounts := []int{1, 4, 8, 16}

	// Создаём временную директорию для файлов
	dir := b.TempDir()

	// Для каждого набора данных создаём пару файлов
	for _, tc := range testCases {
		// Сериализуем данные в JSON
		beforeData, err := json.MarshalIndent(tc.before, "", "  ")
		if err != nil {
			b.Fatalf("failed to marshal before: %v", err)
		}
		afterData, err := json.MarshalIndent(tc.after, "", "  ")
		if err != nil {
			b.Fatalf("failed to marshal after: %v", err)
		}

		// Имена файлов
		beforeFile := filepath.Join(dir, tc.name+"_before.json")
		afterFile := filepath.Join(dir, tc.name+"_after.json")

		// Записываем файлы
		if err := os.WriteFile(beforeFile, beforeData, 0644); err != nil {
			b.Fatalf("failed to write before file: %v", err)
		}
		if err := os.WriteFile(afterFile, afterData, 0644); err != nil {
			b.Fatalf("failed to write after file: %v", err)
		}
	}

	// Карты файлов для каждого набора
	filesBefore := make(map[string]string)
	filesAfter := make(map[string]string)
	for _, tc := range testCases {
		filesBefore[tc.name] = filepath.Join(dir, tc.name+"_before.json")
		filesAfter[tc.name] = filepath.Join(dir, tc.name+"_after.json")
	}

	// Перенаправляем stdout в /dev/null, чтобы прогресс-бар и вывод статистики не мешали
	// Сохраняем оригинальный stdout
	oldStdout := os.Stdout
	devNull, err := os.OpenFile(os.DevNull, os.O_WRONLY, 0)
	if err != nil {
		b.Fatalf("failed to open /dev/null: %v", err)
	}
	os.Stdout = devNull
	defer func() {
		os.Stdout = oldStdout
		devNull.Close()
	}()

	ctx := context.Background()

	// Запускаем бенчмарк для каждой комбинации количества воркеров
	for _, workers := range workerCounts {
		b.Run("workers_"+string(rune(workers+'0')), func(b *testing.B) {
			// Для точности измеряем только время выполнения JSONs
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				// Запускаем JSONs (каждый раз читаем одни и те же файлы)
				_ = JSONs(ctx, filesBefore, filesAfter, workers)
			}
		})
	}
}

// Вспомогательная функция для генерации большого среза
func generateLargeSlice(size int) []interface{} {
	slice := make([]interface{}, size)
	for i := 0; i < size; i++ {
		slice[i] = i
	}
	return slice
}

func BenchmarkFindDifferences(b *testing.B) {
	// Тестовые данные

	equal := storage.DB{
		"user": map[string]interface{}{
			"name": "John",
			"age":  31, // Изменение
			"address": map[string]interface{}{
				"city":    "Moscow",
				"street":  "Tverskaya", // Без изменений
				"zipcode": "123457",    // Изменение
			},
		},
		"items": []interface{}{
			map[string]interface{}{"id": 1, "name": "item1"},
			map[string]interface{}{"id": 2, "name": "item2_modified"}, // Изменение
		},
	}

	largeSliceBefore := storage.DB{
		"data": generateLargeSlice(1000),
	}
	largeSliceAfter := storage.DB{
		"data": generateLargeSlice(1000),
	}

	tests := []struct {
		name          string
		before, after storage.DB
	}{
		{
			name: "simple",
			before: storage.DB{
				"id":     1,
				"name":   "test",
				"price":  99.99,
				"active": true,
			},
			after: storage.DB{
				"id":     1,
				"name":   "test_modified",
				"price":  109.99,
				"active": true,
			},
		},
		{
			name: "nested",
			before: storage.DB{
				"user": map[string]interface{}{
					"name": "John",
					"age":  30,
					"address": map[string]interface{}{
						"city":    "Moscow",
						"street":  "Tverskaya",
						"zipcode": "123456",
					},
				},
				"items": []interface{}{
					map[string]interface{}{"id": 1, "name": "item1"},
					map[string]interface{}{"id": 2, "name": "item2"},
				},
			},
			after: storage.DB{
				"user": map[string]interface{}{
					"name": "John",
					"age":  31, // Изменение
					"address": map[string]interface{}{
						"city":    "Moscow",
						"street":  "Tverskaya", // Без изменений
						"zipcode": "123457",    // Изменение
					},
				},
				"items": []interface{}{
					map[string]interface{}{"id": 1, "name": "item1"},
					map[string]interface{}{"id": 2, "name": "item2_modified"}, // Изменение
				},
			},
		},
		{
			name: "large_slice", before: largeSliceBefore, after: largeSliceAfter},
		{
			name:   "identical",
			before: equal,
			after:  equal,
		},
		{
			name:   "empty",
			before: storage.DB{},
			after:  storage.DB{}},
	}

	for _, test := range tests {
		b.Run(test.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				findDifferences(test.before, test.after, "")
			}
		})
	}
}
