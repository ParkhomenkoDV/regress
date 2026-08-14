package read

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func BenchmarkJSONFiles(b *testing.B) {
	tmpDir := b.TempDir()
	// Создаём 1000 файлов, из которых половина .json
	for i := 0; i < 1000; i++ {
		ext := ".json"
		if i%2 == 0 {
			ext = ".txt"
		}
		name := filepath.Join(tmpDir, fmt.Sprintf("file%d%s", i, ext))
		if err := os.WriteFile(name, []byte("data"), 0644); err != nil {
			b.Fatal(err)
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = JSONFiles(tmpDir)
	}
}

func BenchmarkJSON(b *testing.B) {
	// Создаем временный файл для тестирования
	tempDir := b.TempDir()
	testData := map[string]any{
		"id":   1,
		"name": "Test Document",
		"metadata": map[string]any{
			"created": "2024-01-01",
			"updated": "2024-01-02",
		},
		"items": []any{1, 2, 3, 4, 5},
	}

	data, _ := json.Marshal(testData)
	filePath := filepath.Join(tempDir, "test.json")
	os.WriteFile(filePath, data, 0644)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		JSON(filePath)
	}
}
