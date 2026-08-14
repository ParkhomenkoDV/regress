package utils

import (
	"testing"
)

func TestIsMap(t *testing.T) {
	tests := []struct {
		name string
		v    any
		want bool
	}{
		{"nil", nil, false},
		{"map[string]any пустая", map[string]any{}, true},
		{"map[string]any с данными", map[string]any{"a": 1}, true},
		{"map[string]string", map[string]string{"a": "b"}, false},
		{"map[int]any", map[int]any{1: "a"}, false},
		{"slice", []any{}, false},
		{"int", 42, false},
		{"string", "hello", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsMap(tt.v); got != tt.want {
				t.Errorf("IsMap(%v) = %v, want %v", tt.v, got, tt.want)
			}
		})
	}
}

func TestIsSlice(t *testing.T) {
	tests := []struct {
		name string
		v    any
		want bool
	}{
		{"nil", nil, false},
		{"[]any пустой", []any{}, true},
		{"[]any с данными", []any{1, "a"}, true},
		{"[]int", []int{1, 2}, false},
		{"[]string", []string{"a"}, false},
		{"map", map[string]any{}, false},
		{"int", 10, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsSlice(tt.v); got != tt.want {
				t.Errorf("IsSlice(%v) = %v, want %v", tt.v, got, tt.want)
			}
		})
	}
}

func TestIsSimpleSlice(t *testing.T) {
	tests := []struct {
		name  string
		slice []any
		want  bool
	}{
		{"nil слайс", nil, true},
		{"пустой слайс", []any{}, true},
		{"только простые типы", []any{1, "a", true, 3.14, nil}, true},
		{"содержит map", []any{1, map[string]any{"a": 1}}, false},
		{"содержит slice", []any{1, []any{2}}, false},
		{"содержит структуру", []any{1, struct{}{}}, false},
		{"все числа разных типов", []any{int8(1), uint16(2), float32(3.0)}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsSimpleSlice(tt.slice); got != tt.want {
				t.Errorf("IsSimpleSlice(%v) = %v, want %v", tt.slice, got, tt.want)
			}
		})
	}
}

func TestIsEqualSimpleSlices(t *testing.T) {
	tests := []struct {
		name string
		a    []any
		b    []any
		want bool
	}{
		{"оба nil", nil, nil, true},
		{"оба пустые", []any{}, []any{}, true},
		{"равные значения", []any{1, "a", true}, []any{1, "a", true}, true},
		{"разная длина", []any{1, 2}, []any{1}, false},
		{"разные значения", []any{1, "a"}, []any{2, "a"}, false},
		{"разные типы (числа)", []any{1, 2}, []any{1.0, 2.0}, false}, // 1.0 != 1 (int vs float) – false, т.к. используется ==
		{"nil vs пустой", nil, []any{}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsEqualSimpleSlices(tt.a, tt.b); got != tt.want {
				t.Errorf("IsEqualSimpleSlices(%v, %v) = %v, want %v", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestIsEqual(t *testing.T) {
	tests := []struct {
		name string
		a    any
		b    any
		want bool
	}{
		// nil случаи
		{"nil, nil", nil, nil, true},
		{"nil, not nil", nil, 1, false},
		{"not nil, nil", 1, nil, false},

		// простые типы
		{"string равные", "abc", "abc", true},
		{"string разные", "abc", "def", false},
		{"bool равные", true, true, true},
		{"bool разные", true, false, false},

		// числа: сравнение через fmt.Sprintf, поэтому int и float с одинаковым значением равны
		{"int равные", 5, 5, true},
		{"int разные", 5, 6, false},
		{"int vs float равные", 5, 5.0, true},
		{"int vs float разные", 5, 5.1, false},
		{"int vs string", 5, "5", false}, // не числа, а строка – false // TODO

		// слайсы
		{"slice пустые", []any{}, []any{}, true},
		{"slice nil vs пустой", []any(nil), []any{}, true}, // nil slice считается равным пустому? По логике: len(nil)=0, len(empty)=0, цикл не идет => true
		{"slice одинаковые простые", []any{1, "a", true}, []any{1, "a", true}, true},
		{"slice разные значения", []any{1, 2}, []any{1, 3}, false},
		{"slice разная длина", []any{1}, []any{1, 2}, false},
		{"slice вложенные", []any{1, []any{2, 3}}, []any{1, []any{2, 3}}, true},
		{"slice вложенные разные", []any{1, []any{2, 4}}, []any{1, []any{2, 3}}, false},

		// карты
		{"map пустые", map[string]any{}, map[string]any{}, true},
		{"map nil vs пустая", map[string]any(nil), map[string]any{}, true}, // аналогично
		{"map одинаковые", map[string]any{"a": 1, "b": "x"}, map[string]any{"a": 1, "b": "x"}, true},
		{"map разные ключи", map[string]any{"a": 1}, map[string]any{"b": 1}, false},
		{"map разные значения", map[string]any{"a": 1}, map[string]any{"a": 2}, false},
		{"map вложенные", map[string]any{"a": map[string]any{"b": 1}}, map[string]any{"a": map[string]any{"b": 1}}, true},
		{"map вложенные разные", map[string]any{"a": map[string]any{"b": 1}}, map[string]any{"a": map[string]any{"b": 2}}, false},

		// смешанные типы
		{"map vs slice", map[string]any{}, []any{}, false},
		{"int vs float равные", 42, 42.0, true},
		{"uint vs int", uint(10), 10, true},
		{"int8 vs int64", int8(3), int64(3), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsEqual(tt.a, tt.b); got != tt.want {
				t.Errorf("IsEqual(%v, %v) = %v, want %v", tt.a, tt.b, got, tt.want)
			}
		})
	}
}
