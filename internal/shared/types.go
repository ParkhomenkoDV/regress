package shared

// Difference описывает одно различие
type Difference struct {
	Field  string `doc:"Поле"`
	Before any    `doc:"Значение до"`
	After  any    `doc:"Значение после"`
}

// Comparison содержит сравнение одного файла
type Comparison struct {
	FileName     string       `doc:"Имя файла"`
	ExistsInBoth bool         `doc:"Флаг существования обоих файлов"`
	Differences  []Difference `doc:"Различия"`
}
