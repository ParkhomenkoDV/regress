package utils

import (
	"fmt"
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

	// Проверяем явно на тип []interface{}
	switch v.(type) {
	case []interface{}:
		return true
	default:
		return false
	}
}

// IsEqual проверяет равенство двух значений без использования reflect.DeepEqual
func IsEqual(a, b interface{}) bool {
	if a == nil && b == nil {
		return true
	} else if a == nil || b == nil {
		return false
	}

	// Проверяем простые типы
	switch aVal := a.(type) {
	case string:
		if bVal, ok := b.(string); ok {
			return aVal == bVal
		}
		return false
	case bool:
		if bVal, ok := b.(bool); ok {
			return aVal == bVal
		}
		return false
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		// Для чисел используем форматирование в строку для сравнения
		return fmt.Sprintf("%v", a) == fmt.Sprintf("%v", b)
	case []interface{}:
		if bVal, ok := b.([]interface{}); ok {
			if len(aVal) != len(bVal) {
				return false
			}
			for i := range aVal {
				if !IsEqual(aVal[i], bVal[i]) {
					return false
				}
			}
			return true
		}
		return false
	case map[string]interface{}:
		if bVal, ok := b.(map[string]interface{}); ok {
			if len(aVal) != len(bVal) {
				return false
			}
			for key := range aVal {
				if !IsEqual(aVal[key], bVal[key]) {
					return false
				}
			}
			return true
		}
		return false
	default:
		// Для остальных типов используем строгое сравнение
		return a == b
	}
}
