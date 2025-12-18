package utils

import (
	"os"
	"reflect"
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

// IsMap проверяет, является ли значение словарем
func IsMap(v interface{}) bool {
	_, ok := v.(map[string]interface{})
	return ok
}

// IsSlice проверяет, является ли значение срезом/массивом
func IsSlice(v interface{}) bool {
	if v == nil {
		return false
	}
	return reflect.TypeOf(v).Kind() == reflect.Slice ||
		reflect.TypeOf(v).Kind() == reflect.Array
}
