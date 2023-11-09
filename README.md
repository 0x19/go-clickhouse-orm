# Golang Clickhouse ORM

Playground, do not use this code yet... It works partially.

## Design Choices

- As minimal reflection as possible. It looks very sweet to use go struct tags to define model without additional
hussle, however that includes that for each insert/update/delete/select we need to deal with reflection which will slow down operations quite a lot. 

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