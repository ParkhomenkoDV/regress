package read

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

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
		JSON(filePath)
	}
}
