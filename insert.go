package gchm

import (
	"context"
	"fmt"
	"reflect"

	"github.com/0x19/go-clickhouse-model/dml"
	"github.com/0x19/go-clickhouse-model/models"
)

type InsertBuilder[T models.Model] struct {
	ctx     context.Context
	orm     *ORM
	model   T
	builder *dml.InsertBuilder
}

func (b *InsertBuilder[T]) Build() (string, error) {
	return b.builder.Build()
}

func (b *InsertBuilder[T]) Exec() error {
	return nil
}

func (b *InsertBuilder[T]) ExecContext(ctx context.Context) error {
	return nil
}

func (b *InsertBuilder[T]) SQL() string {
	return b.builder.String()
}

func NewInsert[T models.Model](ctx context.Context, orm *ORM, model T) (T, *InsertBuilder[T], error) {
	// Check if the underlying value of the interface is nil. Unfortunately, it is a T and we cannot
	// directly check if it's nil due to type missmatch.
	{
		modelValue := reflect.ValueOf(model)

		if !modelValue.IsValid() {
			return model, nil, fmt.Errorf("underlying value of model cannot be nil")
		}

		if modelValue.Kind() == reflect.Ptr && modelValue.IsNil() {
			return model, nil, fmt.Errorf("underlying value of model cannot be nil")
		}
	}

	stmtBuilder := dml.NewInsertBuilder()
	stmtBuilder.Model(model)

	builder := &InsertBuilder[T]{
		ctx:     ctx,
		orm:     orm,
		model:   model,
		builder: stmtBuilder,
	}

	sql, err := stmtBuilder.Build()
	if err != nil {
		return model, builder, err
	}

	_ = sql

	return model, builder, nil
}
