package read

import (
	"encoding/json"
	"os"
	"regress/internal/storage"
	"strings"
)

func JSONFiles(dir string) ([]string, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	result := make([]string, 0, len(files))
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(strings.ToLower(file.Name()), ".json") {
			result = append(result, file.Name())
		}
	}
	return result, nil
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
