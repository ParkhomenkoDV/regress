package shared

import "regress/internal/storage"

// Difference описывает одно различие
type Difference struct {
	Field  string `doc:"Поле"`
	Before any    `doc:"Значение до"`
	After  any    `doc:"Значение после"`
}

// Comparison содержит сравнение одного файла
type Comparison struct {
	FileName     string       `doc:"Имя файла"`
	Before       storage.DB   `doc:"До"`
	After        storage.DB   `doc:"После"`
	ExistsInBoth bool         `doc:"Флаг существования в обоих файлах"`
	Differences  []Difference `doc:"Различия"`
}
