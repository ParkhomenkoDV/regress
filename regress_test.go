package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"regress/internal/storage"
	"testing"
)

func BenchmarkFindDifferences(b *testing.B) {
	// Тестовые данные
	simpleBefore := storage.DB{
		"id":     1,
		"name":   "test",
		"price":  99.99,
		"active": true,
	}

	simpleAfter := storage.DB{
		"id":     1,
		"name":   "test_modified",
		"price":  109.99,
		"active": true,
	}

	nestedBefore := storage.DB{
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
	}

	nestedAfter := storage.DB{
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

	testCases := []struct {
		name          string
		before, after storage.DB
	}{
		{"simple", simpleBefore, simpleAfter},
		{"nested", nestedBefore, nestedAfter},
		{"large_slice", largeSliceBefore, largeSliceAfter},
		{"identical", simpleBefore, simpleBefore},
		{"empty", storage.DB{}, storage.DB{}},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				findDifferences(tc.before, tc.after, "")
			}
		})
	}
}

func BenchmarkFindDifferencesWithSimpleSlices(b *testing.B) {
	// Тестирование срезами простых типов
	before := storage.DB{
		"tags":  []interface{}{"go", "test", "benchmark"},
		"ids":   []interface{}{1, 2, 3, 4, 5},
		"flags": []interface{}{true, false, true},
	}

	after := storage.DB{
		"tags":  []interface{}{"go", "test", "performance"}, // Одно изменение
		"ids":   []interface{}{1, 2, 3, 4, 5},               // Без изменений
		"flags": []interface{}{true, false, true},           // Без изменений
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		findDifferences(before, after, "")
	}
}

func BenchmarkFindDifferencesWithComplexSlices(b *testing.B) {
	// Тестирование срезами сложных типов
	before := storage.DB{
		"users": []interface{}{
			map[string]interface{}{"id": 1, "name": "Alice"},
			map[string]interface{}{"id": 2, "name": "Bob"},
			map[string]interface{}{"id": 3, "name": "Charlie"},
		},
	}

	after := storage.DB{
		"users": []interface{}{
			map[string]interface{}{"id": 1, "name": "Alice"},
			map[string]interface{}{"id": 2, "name": "Bob Modified"},
			map[string]interface{}{"id": 3, "name": "Charlie"},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		findDifferences(before, after, "")
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

func BenchmarkReadDynamicJSON(b *testing.B) {
	// Создаем временный файл для тестирования
	tempDir := b.TempDir()
	testData := map[string]interface{}{
		"id":   1,
		"name": "Test Document",
		"metadata": map[string]interface{}{
			"created": "2024-01-01",
			"updated": "2024-01-02",
		},
		"items": []interface{}{1, 2, 3, 4, 5},
	}

	data, _ := json.Marshal(testData)
	filePath := filepath.Join(tempDir, "test.json")
	os.WriteFile(filePath, data, 0644)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		readDynamicJSON(filePath)
	}
}

func BenchmarkCompareFile(b *testing.B) {
	// Создаем тестовые директории и файлы
	tempDir := b.TempDir()
	beforeDir := filepath.Join(tempDir, "before")
	afterDir := filepath.Join(tempDir, "after")
	os.MkdirAll(beforeDir, 0755)
	os.MkdirAll(afterDir, 0755)

	// Создаем идентичные файлы
	testData := map[string]interface{}{
		"id":    1,
		"title": "Test Document",
		"data":  generateTestData(100), // 100 элементов
	}

	data, _ := json.Marshal(testData)
	beforeFile := filepath.Join(beforeDir, "test.json")
	afterFile := filepath.Join(afterDir, "test.json")
	os.WriteFile(beforeFile, data, 0644)
	os.WriteFile(afterFile, data, 0644)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		compareFile("test.json", beforeDir, afterDir, true)
	}
}

func BenchmarkCompareFileWithDifferences(b *testing.B) {
	tempDir := b.TempDir()
	beforeDir := filepath.Join(tempDir, "before")
	afterDir := filepath.Join(tempDir, "after")
	os.MkdirAll(beforeDir, 0755)
	os.MkdirAll(afterDir, 0755)

	// Разные данные
	beforeData := map[string]interface{}{
		"id":    1,
		"title": "Original Title",
		"count": 100,
		"tags":  []interface{}{"a", "b", "c"},
	}

	afterData := map[string]interface{}{
		"id":    1,
		"title": "Modified Title",             // Изменение
		"count": 200,                          // Изменение
		"tags":  []interface{}{"a", "b", "d"}, // Изменение
	}

	beforeJSON, _ := json.Marshal(beforeData)
	afterJSON, _ := json.Marshal(afterData)

	beforeFile := filepath.Join(beforeDir, "test.json")
	afterFile := filepath.Join(afterDir, "test.json")

	os.WriteFile(beforeFile, beforeJSON, 0644)
	os.WriteFile(afterFile, afterJSON, 0644)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		compareFile("test.json", beforeDir, afterDir, true)
	}
}

// Вспомогательная функция для генерации тестовых данных
func generateTestData(count int) map[string]interface{} {
	data := make(map[string]interface{})
	for i := 0; i < count; i++ {
		key := string(rune('a'+i%26)) + string(rune('a'+(i/26)%26))
		data[key] = i
	}
	return data
}
