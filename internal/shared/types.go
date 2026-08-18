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
	ExistsBefore bool         `doc:"Флаг существования в BEFORE"`
	ExistsAfter  bool         `doc:"Флаг существования в AFTER"`
	Differences  []Difference `doc:"Различия"`
}

// Существует в обоих файлах: before и after.
func (c *Comparison) ExistsInBoth() bool {
	return c.ExistsBefore && c.ExistsAfter
}
