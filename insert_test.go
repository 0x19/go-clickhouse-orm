package gchm

import (
	"context"
	"testing"

	"github.com/0x19/go-clickhouse-model/models"
	"github.com/stretchr/testify/assert"
)

type DummyModel struct {
	models.Model
}

func (d *DummyModel) TableName() string {
	return "dummy_model"
}

func TestInsertBuilder(t *testing.T) {
	tests := []struct {
		name          string
		ctx           context.Context
		ormConfig     *Config
		model         models.Model
		wantOrmErr    bool
		wantInsertErr bool
	}{
		{
			name: "Basic Insert With No Model",
			ctx:  context.TODO(),
			ormConfig: &Config{
				Host:     "localhost",
				Port:     9000,
				Username: "default",
				Password: "local12345",
				Database: "unpack",
				Insecure: true,
			},
			wantInsertErr: true,
		},
		{
			name: "Basic Insert With Model",
			ctx:  context.TODO(),
			ormConfig: &Config{
				Host:     "localhost",
				Port:     9000,
				Username: "default",
				Password: "local12345",
				Database: "unpack",
				Insecure: true,
			},
			model:         &DummyModel{},
			wantInsertErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tAssert := assert.New(t)
			orm, err := NewORM(tt.ctx, tt.ormConfig)
			if tt.wantOrmErr {
				tAssert.Error(err)
				return
			}

			tAssert.NoError(err)
			tAssert.NotNil(orm)

			insertBuilder, err := NewInsert(tt.ctx, tt.model)
			if tt.wantInsertErr {
				tAssert.Error(err)
				return
			}

			tAssert.NoError(err)
			tAssert.NotNil(insertBuilder)

		})
	}
}
