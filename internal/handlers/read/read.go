package read

import (
	"os"
	"regress/internal/storage"
	"strings"

	json "github.com/goccy/go-json"
)

func JSONFiles(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	files := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		// регистронезависимая проверка расширения .json
		if strings.HasSuffix(strings.ToLower(name), ".json") {
			files = append(files, name)
		}
	}
	return files, nil
}

func JSON(filePath string) (db storage.DB, err error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(data, &db); err != nil {
		return nil, err
	}

	return db, nil
}
