package utils

import (
	"os"
	"strings"
)

// Вспомогательные функции
func GetJSONFiles(dir string) ([]string, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var result []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(strings.ToLower(file.Name()), ".json") {
			result = append(result, file.Name())
		}
	}
	return result, nil
}

func IsMap(v interface{}) bool {
	_, ok := v.(map[string]interface{})
	return ok
}
