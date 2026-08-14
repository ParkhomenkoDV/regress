# Benchmarks
```
goos: darwin
goarch: arm64
pkg: regress/internal/handlers/compare
cpu: Apple M4
BenchmarkJSONs/workers_1-10         5355            225777 ns/op          145511 B/op       4287 allocs/op
BenchmarkJSONs/workers_4-10         8121            155864 ns/op          145798 B/op       4294 allocs/op
BenchmarkJSONs/workers_8-10         7918            162510 ns/op          146150 B/op       4307 allocs/op
BenchmarkJSONs/workers_@-10         7990            162976 ns/op          146874 B/op       4322 allocs/op
BenchmarkFindDifferences/simple-10               2433264               479.7 ns/op           210 B/op          3 allocs/op
BenchmarkFindDifferences/nested-10                738666              1660 ns/op             880 B/op         19 allocs/op
BenchmarkFindDifferences/large_slice-10           316468              3792 ns/op              67 B/op          1 allocs/op
BenchmarkFindDifferences/identical-10             899026              1329 ns/op             559 B/op         15 allocs/op
BenchmarkFindDifferences/empty-10               62306240                19.50 ns/op            0 B/op          0 allocs/op
```