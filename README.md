# regress

Сравнение двух потоков json файлов

![](./assets/images/testing.png)

## Usage

1. Создать директорию `before` и положи туда JSON файлы
1. Создай директорию `after` и положи туда **одноименные** JSON файлы
1. Запусти скрипт командой

```bash
go run regress.go [--before=<путь>] [--after=<путь>] [--all] [--workers=N]
```

## Build

```bash
go build -o regress regress.go
./regress
```

## Benchmarks
```
go test -bench=. -benchmem -benchtime 5s -count=3
goos: darwin
goarch: arm64
pkg: regress
cpu: Apple M4
BenchmarkFindDifferences/simple-10                      14819276               404.5 ns/op           208 B/op          3 allocs/op
BenchmarkFindDifferences/simple-10                      14647569               413.0 ns/op           208 B/op          3 allocs/op
BenchmarkFindDifferences/simple-10                      14491422               412.0 ns/op           208 B/op          3 allocs/op
BenchmarkFindDifferences/nested-10                       4347871              1387 ns/op             872 B/op         19 allocs/op
BenchmarkFindDifferences/nested-10                       4317838              1387 ns/op             872 B/op         19 allocs/op
BenchmarkFindDifferences/nested-10                       4309172              1391 ns/op             872 B/op         19 allocs/op
BenchmarkFindDifferences/large_slice-10                  1754826              3419 ns/op              48 B/op          1 allocs/op
BenchmarkFindDifferences/large_slice-10                  1754550              3418 ns/op              48 B/op          1 allocs/op
BenchmarkFindDifferences/large_slice-10                  1747910              3419 ns/op              48 B/op          1 allocs/op
BenchmarkFindDifferences/identical-10                   14158467               426.1 ns/op           202 B/op          3 allocs/op
BenchmarkFindDifferences/identical-10                   14015692               426.4 ns/op           202 B/op          3 allocs/op
BenchmarkFindDifferences/identical-10                   13983896               428.0 ns/op           202 B/op          3 allocs/op
BenchmarkFindDifferences/empty-10                       357706972               16.88 ns/op            0 B/op          0 allocs/op
BenchmarkFindDifferences/empty-10                       355446472               16.87 ns/op            0 B/op          0 allocs/op
BenchmarkFindDifferences/empty-10                       353741350               16.76 ns/op            0 B/op          0 allocs/op
BenchmarkFindDifferencesWithSimpleSlices-10             21178575               282.1 ns/op           192 B/op          3 allocs/op
BenchmarkFindDifferencesWithSimpleSlices-10             21260169               285.2 ns/op           192 B/op          3 allocs/op
BenchmarkFindDifferencesWithSimpleSlices-10             21185482               284.5 ns/op           192 B/op          3 allocs/op
BenchmarkFindDifferencesWithComplexSlices-10             6882717               866.5 ns/op           248 B/op         10 allocs/op
BenchmarkFindDifferencesWithComplexSlices-10             6961256               865.1 ns/op           248 B/op         10 allocs/op
BenchmarkFindDifferencesWithComplexSlices-10             7018177               870.5 ns/op           248 B/op         10 allocs/op
BenchmarkReadDynamicJSON-10                               614155              9972 ns/op            2320 B/op         41 allocs/op
BenchmarkReadDynamicJSON-10                               588736             10125 ns/op            2320 B/op         41 allocs/op
BenchmarkReadDynamicJSON-10                               596290              9939 ns/op            2320 B/op         41 allocs/op
BenchmarkCompareFile-10                                   102414             58049 ns/op           39210 B/op        935 allocs/op
BenchmarkCompareFile-10                                   103400             57712 ns/op           39210 B/op        935 allocs/op
BenchmarkCompareFile-10                                   103934             57664 ns/op           39210 B/op        935 allocs/op
BenchmarkCompareFileWithDifferences-10                    273188             20298 ns/op            4122 B/op         71 allocs/op
BenchmarkCompareFileWithDifferences-10                    293323             20317 ns/op            4122 B/op         71 allocs/op
BenchmarkCompareFileWithDifferences-10                    289054             20243 ns/op            4122 B/op         71 allocs/op
```