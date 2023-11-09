package chorm

import (
	"context"
	"fmt"
	"reflect"

	"github.com/0x19/go-clickhouse-orm/models"
	"github.com/0x19/go-clickhouse-orm/sql"
	"github.com/vahid-sohrabloo/chconn/v3"
	"github.com/vahid-sohrabloo/chconn/v3/chpool"
	"github.com/vahid-sohrabloo/chconn/v3/column"
)

type SelectBuilder[T models.Model] struct {
	*sql.SelectBuilder
	ctx          context.Context
	queryOptions *chconn.QueryOptions
	orm          *ORM
	serverInfo   *chconn.ServerInfo
}

func (s *SelectBuilder[T]) GetBuilder() *sql.SelectBuilder {
	return s.SelectBuilder
}

func (s *SelectBuilder[T]) Build() (string, error) {
	return s.SelectBuilder.Build()
}

func (s *SelectBuilder[T]) ExecContext(ctx context.Context, queryOptions *chconn.QueryOptions, columns ...column.ColumnBasic) error {
	return s.orm.GetConn().InsertWithOption(ctx, s.SQL(), queryOptions, columns...)
}

func (s *SelectBuilder[T]) SQL() string {
	return s.SelectBuilder.String()
}

func (s *SelectBuilder[T]) One(ctx context.Context, queryOptions *chconn.QueryOptions, record T) error {
	if queryOptions != nil {
		s.queryOptions = queryOptions
	}

	return nil
}

func (s *SelectBuilder[T]) Scan(ctx context.Context, queryOptions *chconn.QueryOptions) ([]T, error) {
	var toReturn []T

	stmt, err := s.orm.GetConn().SelectWithOption(ctx, s.SQL(), queryOptions, s.GetModel().GetDeclaration().GetPreparedColumns()...)
	if err != nil {
		return toReturn, err
	}
	defer stmt.Close()

	for stmt.Next() {
		elm := reflect.TypeOf((*T)(nil)).Elem()
		record := reflect.New(elm.Elem()).Interface().(T)
		if err := record.ScanRow(stmt); err != nil {
			return toReturn, fmt.Errorf(
				"error scanning row into model `%T`: %w",
				elm.Name(),
				err,
			)
		}

		toReturn = append(toReturn, record)
	}

	if stmt.Err() != nil {
		return toReturn, stmt.Err()
	}

	return toReturn, nil
}

func NewSelect[T models.Model](ctx context.Context, orm *ORM, queryOptions *chconn.QueryOptions) (*SelectBuilder[T], error) {
	stmtBuilder := sql.NewSelectBuilder()
	stmtBuilder.Database(orm.GetDatabaseName())
	stmtBuilder.Select("*") // Initially select all... this can be overridden later.

	tType := reflect.TypeOf((*T)(nil)).Elem()
	instancePtr := reflect.New(tType.Elem()).Interface()

	// Assert that the created instance satisfies the models.Model interface.
	modelInstance, ok := instancePtr.(models.Model)
	if !ok {
		return nil, fmt.Errorf("instance does not satisfy the models.Model interface")
	}

	stmtBuilder.Model(modelInstance)

	builder := &SelectBuilder[T]{
		ctx:           ctx,
		orm:           orm,
		SelectBuilder: stmtBuilder,
		queryOptions:  queryOptions,
	}

	orm.GetConn().AcquireFunc(ctx, func(conn chpool.Conn) error {
		builder.serverInfo = conn.Conn().ServerInfo()
		return nil
	})

	return builder, nil
}
