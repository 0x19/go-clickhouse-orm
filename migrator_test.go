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

	migrator := orm.GetMigrator()
	tAssert.NotNil(migrator)

	err = migrator.Setup(ctx, nil)
	tAssert.NoError(err)

	err = migrator.RegisterMigration(
		"some-name",
		func(ctx context.Context, orm *ORM, migrator *Migrator) error {
			return nil
		},
		func(ctx context.Context, orm *ORM, migrator *Migrator) error {
			return nil
		},
	)
	tAssert.NoError(err)

	tAssert.Equal(1, len(migrator.GetMigrations()))

	err = migrator.Migrate(ctx, nil)
	tAssert.NoError(err)

	//err = migrator.Destroy(ctx, nil)
	//tAssert.NoError(err)
}
