package chorm

import (
	"context"
	"fmt"
	"reflect"

	"github.com/0x19/go-clickhouse-orm/models"
	"github.com/0x19/go-clickhouse-orm/sql"
	"github.com/vahid-sohrabloo/chconn/v3"
	"github.com/vahid-sohrabloo/chconn/v3/column"
)

type UpdateBuilder[T models.Model] struct {
	ctx     context.Context
	orm     *ORM
	model   T
	builder *sql.UpdateBuilder
}

func (b *UpdateBuilder[T]) Build() (string, error) {
	return b.builder.Build()
}

func (b *UpdateBuilder[T]) ExecContext(ctx context.Context, queryOptions *chconn.QueryOptions, columns ...column.ColumnBasic) error {
	return b.orm.GetConn().InsertWithOption(ctx, b.SQL(), queryOptions, columns...)
}

func (b *UpdateBuilder[T]) SQL() string {
	return b.builder.String()
}

func NewUpdate[T models.Model](ctx context.Context, orm *ORM, model T, queryOptions *chconn.QueryOptions) (T, *UpdateBuilder[T], error) {
	// Check if the underlying value of the interface is nil. Unfortunately, it is a T and we cannot
	// directly check if it's nil due to type missmatch.
	{
		modelValue := reflect.ValueOf(model)

		if !modelValue.IsValid() {
			return model, nil, fmt.Errorf("model cannot be nil")
		}

		if modelValue.Kind() == reflect.Ptr && modelValue.IsNil() {
			return model, nil, fmt.Errorf("model cannot be nil")
		}
	}

	stmtBuilder := sql.NewUpdateBuilder()
	stmtBuilder.Model(model)
	declaration := model.GetDeclaration()

	if declaration == nil {
		return model, nil, fmt.Errorf("model declaration cannot be nil")
	}

	if declaration.DatabaseName != "" {
		stmtBuilder.Database(declaration.DatabaseName)
	} else {
		stmtBuilder.Database(orm.GetDatabaseName())
	}

	stmtBuilder.Fields(declaration.GetFieldNames()...)

	builder := &UpdateBuilder[T]{
		ctx:     ctx,
		orm:     orm,
		model:   model,
		builder: stmtBuilder,
	}

	if err := builder.ExecContext(ctx, queryOptions, declaration.GetPreparedColumns()...); err != nil {
		return model, builder, err
	}

	return model, builder, nil
}
