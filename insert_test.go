package gchm

import (
	"context"
	"fmt"
	"testing"

	"github.com/0x19/go-clickhouse-model/models"
	"github.com/stretchr/testify/assert"
)

type TestModel struct {
	models.Model

	Name string `gchm:"cn: name"`
}

func (d *TestModel) TableName() string {
	return "dummy_model"
}

func (d *TestModel) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"name": d.Name,
	}
}

func TestInsertBuilder(t *testing.T) {
	tests := []struct {
		name          string
		ctx           context.Context
		ormConfig     *Config
		model         *TestModel
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
			model:         &TestModel{},
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

			record, err := NewInsert(tt.ctx, orm, tt.model)
			if tt.wantInsertErr {
				tAssert.Error(err)
				return
			}

			tAssert.NoError(err)
			tAssert.NotNil(record)

			result, err := orm.Insert(tt.ctx, tt.model)
			if tt.wantInsertErr {
				tAssert.Error(err)
				return
			}

			dbRecord := result.(*TestModel)

			tAssert.NoError(err)
			tAssert.NotNil(result)

			fmt.Printf("response: %T \n", record.Name)
			fmt.Printf("response orm: %T \n", dbRecord.Name)

			record.ToMap()
		})
	}
}
