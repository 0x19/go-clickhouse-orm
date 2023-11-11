package chorm

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMigrator(t *testing.T) {
	tAssert := assert.New(t)

	ctx, ctxCancel := context.WithCancel(context.Background())
	defer ctxCancel()

	orm, err := NewORM(ctx, &Config{
		Host:     "localhost",
		Port:     9000,
		Username: "default",
		Password: "local12345",
		Database: "unpack",
		Insecure: true,
	})
	defer orm.Close()

	tAssert.NoError(err)
	tAssert.NotNil(orm)
}
