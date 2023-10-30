# Golang Clickhouse ORM
Golang CRUD support for ClickhouseDB


## Benchmarks

### Insert Benchmarking
```
go test -benchmem -run=^$ -bench ^BenchmarkNewInsert$ github.com/0x19/go-clickhouse-model -v

goos: linux
goarch: amd64
pkg: github.com/0x19/go-clickhouse-model
cpu: AMD Ryzen Threadripper 3960X 24-Core Processor 
BenchmarkNewInsert
BenchmarkNewInsert-48            7057206               161.4 ns/op           176 B/op          4 allocs/op
PASS
ok      github.com/0x19/go-clickhouse-model     1.335s
```