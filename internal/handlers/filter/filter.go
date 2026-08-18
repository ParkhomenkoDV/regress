package filter

import "regress/internal/shared"

// ShowAll фильтрует сравнения по настройкам
func ShowAll(comparisons []shared.Comparison, showAll bool) []shared.Comparison {
	if showAll {
		return comparisons
	}

	filtered := make([]shared.Comparison, 0, len(comparisons))
	for _, comp := range comparisons {
		if len(comp.Differences) > 0 || !comp.ExistsInBoth() {
			filtered = append(filtered, comp)
		}
	}
	return filtered
}
