package gchm

import (
	"context"
	"fmt"
	"reflect"

	"github.com/0x19/go-clickhouse-model/models"
)

type InsertBuilder[T models.Model] struct {
	ctx   context.Context
	model T
}

func NewInsert[T models.Model](ctx context.Context, model T) (*InsertBuilder[T], error) {
	// Check if the underlying value of the interface is nil
	if reflect.ValueOf(model).IsValid() == false {
		return nil, fmt.Errorf("underlying value of model cannot be nil")
	}

	fmt.Printf("INSERT INTO %s \n", model.TableName())
	return &InsertBuilder[T]{
		ctx:   ctx,
		model: model,
	}, nil
}
