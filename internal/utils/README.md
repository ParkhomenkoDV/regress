# Benchmarks
```
goos: darwin
goarch: arm64
pkg: regress/internal/utils
cpu: Apple M4
BenchmarkIsMap-10                       1000000000               0.2283 ns/op          0 B/op          0 allocs/op
BenchmarkIsSlice-10                     1000000000               0.2259 ns/op          0 B/op          0 allocs/op
BenchmarkIsSimpleSlice/simple_int-10            445735425                2.720 ns/op           0 B/op          0 allocs/op
BenchmarkIsSimpleSlice/simple_string-10         590063455                2.054 ns/op           0 B/op          0 allocs/op
BenchmarkIsSimpleSlice/mixed-10                 565827439                1.904 ns/op           0 B/op          0 allocs/op
BenchmarkIsSimpleSlice/large_simple-10           2578384               471.9 ns/op             0 B/op          0 allocs/op
BenchmarkIsEqual/int_equal-10                   20880549                59.15 ns/op            4 B/op          2 allocs/op
BenchmarkIsEqual/int_not_equal-10               20520756                59.10 ns/op            4 B/op          2 allocs/op
BenchmarkIsEqual/string_equal-10                418125932                2.936 ns/op           0 B/op          0 allocs/op
BenchmarkIsEqual/string_not_equal-10            355084038                3.442 ns/op           0 B/op          0 allocs/op
BenchmarkIsEqual/slice_equal-10                  7543701               142.5 ns/op             0 B/op          0 allocs/op
BenchmarkIsEqual/slice_not_equal-10              8554220               142.3 ns/op             0 B/op          0 allocs/op
BenchmarkIsEqual/map_equal-10                   32434149                35.75 ns/op            0 B/op          0 allocs/op
BenchmarkIsEqual/map_not_equal-10               38682121                31.07 ns/op            0 B/op          0 allocs/op
BenchmarkIsEqualSimpleSlices/small_equal-10     220589316                5.453 ns/op           0 B/op          0 allocs/op
BenchmarkIsEqualSimpleSlices/small_not_equal-10                 196896096                5.703 ns/op           0 B/op          0 allocs/op
BenchmarkIsEqualSimpleSlices/large_equal-10                       792720              1498 ns/op               0 B/op          0 allocs/op
BenchmarkIsEqualSimpleSlices/large_not_equal-10                  1343343               879.3 ns/op             0 B/op          0 allocs/op
```