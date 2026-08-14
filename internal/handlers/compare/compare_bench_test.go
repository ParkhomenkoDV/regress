package compare

import (
	"encoding/json"
	"os"
	"path/filepath"
	"regress/internal/storage"
	"testing"
)

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
