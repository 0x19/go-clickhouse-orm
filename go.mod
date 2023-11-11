module github.com/0x19/go-clickhouse-orm

go 1.19

require (
	github.com/davecgh/go-spew v1.1.1
	github.com/google/uuid v1.3.1
	github.com/stretchr/testify v1.8.4
	github.com/vahid-sohrabloo/chconn/v3 v3.0.0-00010101000000-000000000000
	go.uber.org/zap v1.26.0
)

replace github.com/vahid-sohrabloo/chconn/v3 => ../chconn

require (
	github.com/go-faster/city v1.0.1 // indirect
	github.com/jackc/puddle/v2 v2.2.1 // indirect
	github.com/klauspost/compress v1.17.0 // indirect
	github.com/pierrec/lz4/v4 v4.1.18 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	go.uber.org/multierr v1.10.0 // indirect
	golang.org/x/sync v0.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
