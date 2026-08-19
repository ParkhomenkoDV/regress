package filter

import (
	"regress/internal/shared"
	"testing"
)

func BenchmarkDiff(b *testing.B) {
	// Создаём слайс из 1000 сравнений, половина с изменениями, половина без, и часть с пропущенными файлами
	const size = 1000
	comparisons := make([]shared.Comparison, size)
	for i := 0; i < size; i++ {
		existsBoth := i%3 != 0 // некоторые отсутствуют
		hasDiff := i%2 == 0
		var diffs []shared.Difference
		if hasDiff {
			diffs = []shared.Difference{{Field: "f", Before: i, After: i + 1}}
		}
		comparisons[i] = shared.Comparison{
			FileName:     string(rune('a' + i%26)),
			ExistsBefore: existsBoth || (i%4 == 1), // для разнообразия
			ExistsAfter:  existsBoth || (i%4 == 2),
			Differences:  diffs,
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Diff(comparisons)
	}
}
