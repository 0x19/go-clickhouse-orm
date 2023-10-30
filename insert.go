package gchm

import "context"

type InsertBuilder struct {
	ctx context.Context
}

func NewInsert() *InsertBuilder {
	return &InsertBuilder{}
}
