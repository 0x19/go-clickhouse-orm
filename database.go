package chorm

import (
	"context"
	"errors"
	"fmt"

	"github.com/0x19/go-clickhouse-orm/sql"
	"github.com/vahid-sohrabloo/chconn/v3"
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

func NewCreateDatabase(ctx context.Context, orm *ORM, dbName string, useDb bool, queryOptions *chconn.QueryOptions) (*DatabaseBuilder, error) {
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

	// Database wont be automatically selected when created. In case that the user wants to use the database
	// after creating it, we need to select it.
	if useDb {
		if err := orm.GetConn().ExecWithOption(ctx, fmt.Sprintf("USE %s;", dbName), queryOptions); err != nil {
			return builder, err
		}
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
