package chorm

import (
	"context"
	"fmt"
	"reflect"

	"github.com/0x19/go-clickhouse-model/models"
	"github.com/0x19/go-clickhouse-model/sql"
	"github.com/vahid-sohrabloo/chconn/v2"
	"github.com/vahid-sohrabloo/chconn/v2/column"
)

type TableBuilder[T models.Model] struct {
	ctx     context.Context
	orm     *ORM
	model   T
	builder *sql.TableBuilder
}

func (b *TableBuilder[T]) Build() (string, error) {
	return b.builder.Build()
}

func (b *TableBuilder[T]) ExecContext(ctx context.Context, queryOptions *chconn.QueryOptions, columns ...column.ColumnBasic) error {
	return b.orm.GetConn().ExecWithOption(ctx, b.SQL(), queryOptions)
}

func (b *TableBuilder[T]) Exec(ctx context.Context) error {
	return b.orm.GetConn().Exec(ctx, b.SQL())
}

func (b *TableBuilder[T]) GetBuilder() *sql.TableBuilder {
	return b.builder
}

func (b *TableBuilder[T]) SQL() string {
	return b.builder.String()
}

func NewCreateTable[T models.Model](ctx context.Context, orm *ORM, model T, queryOptions *chconn.QueryOptions) (T, *TableBuilder[T], error) {
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

	stmtBuilder := sql.NewCreateTableBuilder()
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

	fields := make([]string, 0)
	for _, field := range declaration.GetFields() {
		fields = append(fields, field.GetDDL())
	}
	stmtBuilder.Fields(fields...)

	stmtBuilder.Engine(declaration.Engine)
	stmtBuilder.PartitionBy(declaration.PartitionBy)

	pkFields := make([]string, 0)
	for _, field := range declaration.GetPKFields() {
		pkFields = append(pkFields, field.Name)
	}
	stmtBuilder.PrimaryKeys(pkFields...)

	stmtBuilder.OrderBy(declaration.OrderBy...)
	stmtBuilder.Settings(declaration.Settings...)

	builder := &TableBuilder[T]{
		ctx:     ctx,
		orm:     orm,
		model:   model,
		builder: stmtBuilder,
	}

	if err := builder.Exec(ctx); err != nil {
		return model, builder, err
	}

	return model, builder, nil
}

func NewDropTable[T models.Model](ctx context.Context, orm *ORM, model T, queryOptions *chconn.QueryOptions) (*TableBuilder[T], error) {
	// Check if the underlying value of the interface is nil. Unfortunately, it is a T and we cannot
	// directly check if it's nil due to type missmatch.
	{
		modelValue := reflect.ValueOf(model)

		if !modelValue.IsValid() {
			return nil, fmt.Errorf("model cannot be nil")
		}

		if modelValue.Kind() == reflect.Ptr && modelValue.IsNil() {
			return nil, fmt.Errorf("model cannot be nil")
		}
	}

	stmtBuilder := sql.NewDropTableBuilder()
	stmtBuilder.Model(model)
	declaration := model.GetDeclaration()

	if declaration == nil {
		return nil, fmt.Errorf("model declaration cannot be nil")
	}

	if declaration.DatabaseName != "" {
		stmtBuilder.Database(declaration.DatabaseName)
	} else {
		stmtBuilder.Database(orm.GetDatabaseName())
	}

	builder := &TableBuilder[T]{
		ctx:     ctx,
		orm:     orm,
		model:   model,
		builder: stmtBuilder,
	}

	if err := builder.Exec(ctx); err != nil {
		return builder, err
	}

	return builder, nil
}
