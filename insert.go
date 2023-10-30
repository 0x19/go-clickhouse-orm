package gchm

import (
	"context"
	"fmt"
	"reflect"

	"github.com/0x19/go-clickhouse-model/models"
)

type Modeler interface {
	models.Model
}

type InsertBuilder[T models.Model] struct {
	ctx   context.Context
	orm   *ORM
	model T
}

func NewInsert[T models.Model](ctx context.Context, orm *ORM, model T) (T, error) {
	// Check if the underlying value of the interface is nil. Unfortunately, it is a T and we cannot
	// directly check if it's nil due to type missmatch.
	if !reflect.ValueOf(model).IsValid() {
		return model, fmt.Errorf("underlying value of model cannot be nil")
	}

	fmt.Printf("INSERT INTO %s \n - %T \n", model.TableName(), model)

	builder := InsertBuilder[T]{ctx: ctx, orm: orm, model: model}

	_ = builder

	return model, nil
}
