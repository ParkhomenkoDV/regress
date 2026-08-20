# Benchmarks
```
goos: darwin
goarch: arm64
pkg: regress/internal/handlers/compare
cpu: Apple M4
BenchmarkJSONs/workers_1-10         5618            211005 ns/op          143401 B/op       4283 allocs/op
BenchmarkJSONs/workers_4-10         7636            149298 ns/op          147777 B/op       4291 allocs/op
BenchmarkJSONs/workers_8-10         7560            156674 ns/op          148832 B/op       4299 allocs/op
BenchmarkJSONs/workers_@-10         7413            158179 ns/op          149976 B/op       4315 allocs/op
BenchmarkFindDifferences/simple-10               2846188               417.7 ns/op           208 B/op          3 allocs/op
BenchmarkFindDifferences/nested-10                816177              1464 ns/op             872 B/op         19 allocs/op
BenchmarkFindDifferences/large_slice-10           343484              3381 ns/op              48 B/op          1 allocs/op
BenchmarkFindDifferences/identical-10            1000000              1146 ns/op             552 B/op         15 allocs/op
BenchmarkFindDifferences/empty-10               69029817                17.16 ns/op            0 B/op          0 allocs/op
```