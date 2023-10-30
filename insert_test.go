package gchm

import (
	"context"
	"fmt"
	"testing"

	"github.com/0x19/go-clickhouse-model/models"
	"github.com/stretchr/testify/assert"
	"github.com/vahid-sohrabloo/chconn/v2"
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
		queryOptions  *chconn.QueryOptions
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
			model:         nil,
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
			model: &TestModel{
				Name: "test",
			},
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

			record, builder, err := NewInsert(tt.ctx, orm, tt.model, tt.queryOptions)
			if tt.wantInsertErr {
				tAssert.Error(err)
				return
			}

			tAssert.NoError(err)
			tAssert.NotNil(record)
			tAssert.NotNil(builder)

			fmt.Println("SQL: ", builder.SQL())
			fmt.Printf("response: %+v \n", record)

			record.ToMap()
		})
	}
}

func BenchmarkNewInsert(b *testing.B) {
	// Setup
	ctx := context.TODO()
	ormConfig := &Config{
		Host:     "localhost",
		Port:     9000,
		Username: "default",
		Password: "local12345",
		Database: "unpack",
		Insecure: true,
	}
	model := &TestModel{
		Name: "test",
	}

	orm, err := NewORM(ctx, ormConfig)
	if err != nil {
		b.Fatalf("Failed to create ORM: %v", err)
	}

	// Benchmark loop
	b.ResetTimer() // Reset the timer to exclude setup time
	for i := 0; i < b.N; i++ {
		_, _, err := NewInsert(ctx, orm, model, nil)
		if err != nil {
			b.Fatalf("Failed to insert: %v", err)
		}
	}
}
