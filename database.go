package chorm

import (
	"context"
	"errors"

	"github.com/0x19/go-clickhouse-model/sql"
	"github.com/vahid-sohrabloo/chconn/v2"
)

type DatabaseBuilder struct {
	ctx     context.Context
	orm     *ORM
	builder *sql.DatabaseBuilder
}

func (b *DatabaseBuilder) Build() (string, error) {
	return b.builder.Build()
}

func (b *DatabaseBuilder) ExecContext(ctx context.Context, queryOptions *chconn.QueryOptions) error {
	return b.orm.GetConn().ExecWithOption(ctx, b.SQL(), queryOptions)
}

func (b *DatabaseBuilder) SQL() string {
	return b.builder.String()
}

func NewCreateDatabase(ctx context.Context, orm *ORM, dbName string, queryOptions *chconn.QueryOptions) (*DatabaseBuilder, error) {
	if dbName == "" {
		return nil, errors.New("database name cannot be empty")
	}

	stmtBuilder := sql.NewCreateDatabaseBuilder()
	stmtBuilder.Database(dbName)

	builder := &DatabaseBuilder{
		ctx:     ctx,
		orm:     orm,
		builder: stmtBuilder,
	}

	if err := builder.ExecContext(ctx, queryOptions); err != nil {
		return builder, err
	}

	return builder, nil
}

func NewDropDatabase(ctx context.Context, orm *ORM, dbName string, queryOptions *chconn.QueryOptions) (*DatabaseBuilder, error) {
	if dbName == "" {
		return nil, errors.New("database name cannot be empty")
	}

	stmtBuilder := sql.NewDropDatabaseBuilder()
	stmtBuilder.Database(dbName)

	builder := &DatabaseBuilder{
		ctx:     ctx,
		orm:     orm,
		builder: stmtBuilder,
	}

	if err := builder.ExecContext(ctx, queryOptions); err != nil {
		return builder, err
	}

	return builder, nil
}
