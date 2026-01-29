// [file name]: utils_bench_test.go
package utils_test

import (
	"regress/pkg/utils"
	"testing"
)

func BenchmarkIsMap(b *testing.B) {
	m := map[string]interface{}{"key": "value"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		utils.IsMap(m)
	}
}

func BenchmarkIsSlice(b *testing.B) {
	slice := []interface{}{1, 2, 3, "test"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		utils.IsSlice(slice)
	}
}

func BenchmarkIsSimpleSlice(b *testing.B) {
	testCases := []struct {
		name  string
		slice []interface{}
	}{
		{"simple_int", []interface{}{1, 2, 3, 4}},
		{"simple_string", []interface{}{"a", "b", "c"}},
		{"mixed", []interface{}{1, "test", map[string]interface{}{"key": "value"}}},
		{"large_simple", make([]interface{}, 1000)},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				utils.IsSimpleSlice(tc.slice)
			}
		})
	}
}

func BenchmarkIsEqual(b *testing.B) {
	testCases := []struct {
		name string
		a, b interface{}
	}{
		{"int_equal", 42, 42},
		{"int_not_equal", 42, 43},
		{"string_equal", "hello", "hello"},
		{"string_not_equal", "hello", "world"},
		{"slice_equal", []interface{}{1, 2, 3}, []interface{}{1, 2, 3}},
		{"slice_not_equal", []interface{}{1, 2, 3}, []interface{}{1, 2, 4}},
		{"map_equal", map[string]interface{}{"key": "value"}, map[string]interface{}{"key": "value"}},
		{"map_not_equal", map[string]interface{}{"key": "value"}, map[string]interface{}{"key": "value2"}},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				utils.IsEqual(tc.a, tc.b)
			}
		})
	}
}

func BenchmarkIsEqualSimpleSlices(b *testing.B) {
	small1 := []interface{}{1, 2, 3}
	small2 := []interface{}{1, 2, 4}
	large1 := make([]interface{}, 1000)
	large2 := make([]interface{}, 1000)
	for i := 0; i < 1000; i++ {
		large1[i] = i
		large2[i] = i
	}
	large2[500] = 999 // Одно различие

	b.Run("small_equal", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			utils.IsEqualSimpleSlices(small1, small1)
		}
	})

	b.Run("small_not_equal", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			utils.IsEqualSimpleSlices(small1, small2)
		}
	})

	b.Run("large_equal", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			utils.IsEqualSimpleSlices(large1, large1)
		}
	})

	b.Run("large_not_equal", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			utils.IsEqualSimpleSlices(large1, large2)
		}
	})
}
