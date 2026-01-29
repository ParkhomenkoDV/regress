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
BenchmarkFindDifferences/simple-10                      13514348               430.4 ns/op           160 B/op          4 allocs/op
BenchmarkFindDifferences/simple-10                      13916923               431.0 ns/op           160 B/op          4 allocs/op
BenchmarkFindDifferences/simple-10                      13852612               432.1 ns/op           160 B/op          4 allocs/op
BenchmarkFindDifferences/nested-10                       4255280              1413 ns/op             707 B/op         20 allocs/op
BenchmarkFindDifferences/nested-10                       4162827              1416 ns/op             707 B/op         20 allocs/op
BenchmarkFindDifferences/nested-10                       4230715              1422 ns/op             707 B/op         20 allocs/op
BenchmarkFindDifferences/large_slice-10                  1749848              3443 ns/op               0 B/op          0 allocs/op
BenchmarkFindDifferences/large_slice-10                  1747146              3435 ns/op               0 B/op          0 allocs/op
BenchmarkFindDifferences/large_slice-10                  1746916              3438 ns/op               0 B/op          0 allocs/op
BenchmarkFindDifferences/identical-10                   15123063               394.3 ns/op            10 B/op          2 allocs/op
BenchmarkFindDifferences/identical-10                   15101938               396.2 ns/op            10 B/op          2 allocs/op
BenchmarkFindDifferences/identical-10                   15115576               395.5 ns/op            10 B/op          2 allocs/op
BenchmarkFindDifferences/empty-10                       375870226               15.92 ns/op            0 B/op          0 allocs/op
BenchmarkFindDifferences/empty-10                       375919171               15.90 ns/op            0 B/op          0 allocs/op
BenchmarkFindDifferences/empty-10                       376927483               15.90 ns/op            0 B/op          0 allocs/op
BenchmarkFindDifferencesWithSimpleSlices-10             21837944               273.5 ns/op            96 B/op          3 allocs/op
BenchmarkFindDifferencesWithSimpleSlices-10             21913785               272.3 ns/op            96 B/op          3 allocs/op
BenchmarkFindDifferencesWithSimpleSlices-10             21954864               272.1 ns/op            96 B/op          3 allocs/op
BenchmarkFindDifferencesWithComplexSlices-10             6910350               870.8 ns/op           200 B/op         10 allocs/op
BenchmarkFindDifferencesWithComplexSlices-10             6846963               870.5 ns/op           200 B/op         10 allocs/op
BenchmarkFindDifferencesWithComplexSlices-10             6942277               866.2 ns/op           200 B/op         10 allocs/op
BenchmarkReadDynamicJSON-10                               588855              9931 ns/op            2320 B/op         41 allocs/op
BenchmarkReadDynamicJSON-10                               591286              9945 ns/op            2320 B/op         41 allocs/op
BenchmarkReadDynamicJSON-10                               582547              9937 ns/op            2320 B/op         41 allocs/op
BenchmarkCompareFile-10                                   103023             57937 ns/op           34199 B/op        933 allocs/op
BenchmarkCompareFile-10                                    99825             58378 ns/op           34199 B/op        933 allocs/op
BenchmarkCompareFile-10                                   102177             58246 ns/op           34199 B/op        933 allocs/op
BenchmarkCompareFileWithDifferences-10                    290126             20532 ns/op            4218 B/op         73 allocs/op
BenchmarkCompareFileWithDifferences-10                    288631             20505 ns/op            4266 B/op         73 allocs/op
BenchmarkCompareFileWithDifferences-10                    289946             20974 ns/op            4266 B/op         73 allocs/op
```