package filter

import "regress/internal/shared"

// Diff фильтрует сравнения
func Diff(comparisons []shared.Comparison) (filtered []shared.Comparison) {
	filtered = make([]shared.Comparison, 0, len(comparisons))
	for _, comparison := range comparisons {
		if len(comparison.Differences) > 0 || !comparison.ExistsInBoth() {
			filtered = append(filtered, comparison)
		}
	}
	return filtered
}
