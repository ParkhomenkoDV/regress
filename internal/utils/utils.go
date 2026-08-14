package utils

import (
	"fmt"
)

// IsMap проверяет, является ли значение картой
func IsMap(v any) bool {
	_, ok := v.(map[string]any)
	return ok
}

// IsSlice проверяет, является ли значение срезом/массивом
func IsSlice(v any) bool {
	_, ok := v.([]any)
	return ok
}

// IsSimpleSlice проверяет, содержит ли срез только простые типы
func IsSimpleSlice(slice []any) bool {
	for _, v := range slice {
		switch v.(type) {
		case int, int8, int16, int32, int64,
			uint, uint8, uint16, uint32, uint64,
			float32, float64,
			string,
			bool,
			nil:
		// Простой тип - продолжаем проверку
		default: // Нашли сложный тип
			return false
		}
	}
	return true
}

// IsEqualSimpleSlices сравнивает два среза простых типов
func IsEqualSimpleSlices(a, b []any) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// IsEqual проверяет равенство двух значений без использования reflect.DeepEqual
func IsEqual(a, b any) bool {
	// Проверяем указатели
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
	case []any:
		if bVal, ok := b.([]any); ok {
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
	case map[string]any:
		if bVal, ok := b.(map[string]any); ok {
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
	default: // Для остальных типов используем строгое сравнение
		return a == b
	}
}
