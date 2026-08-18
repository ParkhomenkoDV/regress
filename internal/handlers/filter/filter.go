package filter

import "regress/internal/shared"

// Diff фильтрует сравнения
func Diff(comparisons []shared.Comparison) (filtered []shared.Comparison) {
	filtered = make([]shared.Comparison, 0, len(comparisons))
	for _, comp := range comparisons {
		if len(comp.Differences) > 0 || !comp.ExistsInBoth() {
			filtered = append(filtered, comp)
		}
	}
	return filtered
}
