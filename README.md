# Regress = инструмент сравнения JSON файлов

**Regress** — мощная утилита для сравнения JSON файлов, которая помогает находить различия между двумя версиями данных и экспортировать результаты в удобном формате Excel.

![](./assets/images/testing.png)

## Requirements

- Go 1.21 или выше

## Usage

1. Создать директорию `before` и положи туда JSON файлы
1. Создай директорию `after` и положи туда **одноименные** JSON файлы
1. Запусти скрипт командой

```bash
go run regress.go [-before=S] [-after=S] [-all] [-workers=N]
```

Флаг       | Описание                                      | По умолчанию
-----------|-----------------------------------------------|-------------
`-before`  | Директория с JSON-файлами **ДО** изменений    | `before`
`-after`   | Директория с JSON-файлами **ПОСЛЕ** изменений | `after`
`-all`     | Показывать все поля даже без изменений        | `false`
`-workers` | Количество параллельных воркеров обработки    | количетсво ядер CPU


## Build

```bash
go build -o regress regress.go
./regress
```

## Install

```bash
go install
```

## Schema
```
before/      after/
| 1.json -   | 1.json -
| 2.json  |  | 2.json  |
| 3.json  |  | 3.json  |
| ...     |  | ...     |
          |            |  
          v            v
         ----------------
        | i.json  i.json |
        |    |       |   |
        |    v       v   |
        |    UnMarshal   | gorutine
        |      |   |     |
        |      v   v     |
        |     Compare    |
         ----------------
                 |
                 v
         comparison.xlsx
```

## Benchmarks
```
go test -bench=. -benchmem -benchtime 3s -count=1
goos: darwin
goarch: arm64
pkg: regress
cpu: Apple M4
BenchmarkFindDifferences/simple-10               8528337               409.6 ns/op           208 B/op          3 allocs/op
BenchmarkFindDifferences/nested-10               2595163              1395 ns/op             872 B/op         19 allocs/op
BenchmarkFindDifferences/large_slice-10          1000000              3468 ns/op              48 B/op          1 allocs/op
BenchmarkFindDifferences/identical-10            3231038              1122 ns/op             552 B/op         15 allocs/op
BenchmarkFindDifferences/empty-10               209907532               17.09 ns/op            0 B/op          0 allocs/op
BenchmarkReadDynamicJSON-10                       361873              9769 ns/op            2320 B/op         41 allocs/op
BenchmarkCompareFileWithDifferences-10            176138             20186 ns/op            4122 B/op         71 allocs/op
```