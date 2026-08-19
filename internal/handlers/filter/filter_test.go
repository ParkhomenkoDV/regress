package filter

import (
	"regress/internal/shared"
	"testing"
)

func TestDiff(t *testing.T) {
	tests := []struct {
		name          string
		comparisons   []shared.Comparison
		wantCount     int      // ожидаемое количество после фильтрации
		wantFileNames []string // для проверки конкретных имён
	}{
		{
			name:          "empty",
			comparisons:   []shared.Comparison{},
			wantCount:     0,
			wantFileNames: []string{},
		},
		{
			name: "no differences and exists both",
			comparisons: []shared.Comparison{
				{
					FileName:     "file1",
					ExistsBefore: true,
					ExistsAfter:  true,
					Differences:  []shared.Difference{}, // пусто
				},
			},
			wantCount:     0,
			wantFileNames: []string{},
		},
		{
			name: "has differences and exists both",
			comparisons: []shared.Comparison{
				{
					FileName:     "file2",
					ExistsBefore: true,
					ExistsAfter:  true,
					Differences:  []shared.Difference{{Field: "a", Before: 1, After: 2}}, // хотя бы один
				},
			},
			wantCount:     1,
			wantFileNames: []string{"file2"},
		},
		{
			name: "missing in before",
			comparisons: []shared.Comparison{
				{
					FileName:     "file3",
					ExistsBefore: false,
					ExistsAfter:  true,
					Differences:  []shared.Difference{},
				},
			},
			wantCount:     1,
			wantFileNames: []string{"file3"},
		},
		{
			name: "missing in after",
			comparisons: []shared.Comparison{
				{
					FileName:     "file4",
					ExistsBefore: true,
					ExistsAfter:  false,
					Differences:  []shared.Difference{},
				},
			},
			wantCount:     1,
			wantFileNames: []string{"file4"},
		},
		{
			name: "mixed",
			comparisons: []shared.Comparison{
				{FileName: "keep1", ExistsBefore: true, ExistsAfter: true, Differences: []shared.Difference{{}}},
				{FileName: "skip1", ExistsBefore: true, ExistsAfter: true, Differences: []shared.Difference{}},
				{FileName: "keep2", ExistsBefore: false, ExistsAfter: true, Differences: []shared.Difference{}},
				{FileName: "keep3", ExistsBefore: true, ExistsAfter: false, Differences: []shared.Difference{}},
			},
			wantCount:     3,
			wantFileNames: []string{"keep1", "keep2", "keep3"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Diff(tt.comparisons)
			if len(got) != tt.wantCount {
				t.Errorf("Diff() count = %d, want %d", len(got), tt.wantCount)
			}
			// Проверяем имена файлов
			if tt.wantFileNames != nil {
				gotNames := make([]string, len(got))
				for i, c := range got {
					gotNames[i] = c.FileName
				}
				// Порядок может сохраняться, но для надёжности сравниваем множества
				if !equalStringSlices(gotNames, tt.wantFileNames) {
					t.Errorf("Diff() file names = %v, want %v", gotNames, tt.wantFileNames)
				}
			}
		})
	}
}

func equalStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	m := make(map[string]int)
	for _, s := range a {
		m[s]++
	}
	for _, s := range b {
		m[s]--
		if m[s] < 0 {
			return false
		}
	}
	return true
}
