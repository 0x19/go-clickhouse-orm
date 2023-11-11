# Golang Clickhouse ORM

Playground, do not use this code yet... It works partially.

## Design Choices

- Reflection caching. Reflection is slow, so we cache it by building reflection manager that is initialized only once and then used for all operations. On this way we can avoid reflection on every insert, select, update, delete operation while still keeping the code clean and simple.

## Getting Started

### Installation

Todo...

## Benchmarks

### Insert Benchmarking
```
go test -benchmem -run=^$ -bench ^BenchmarkNewInsert$ github.com/0x19/go-clickhouse-orm -v

goos: linux
goarch: amd64
pkg: github.com/0x19/go-clickhouse-orm
cpu: AMD Ryzen Threadripper 3960X 24-Core Processor 
BenchmarkNewInsert
BenchmarkNewInsert-48               1194            942314 ns/op            6493 B/op         61 allocs/op
PASS
ok      github.com/0x19/go-clickhouse-orm     1.249s
```