package chorm

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vahid-sohrabloo/chconn/v2"
)

func TestInsertBuilder(t *testing.T) {
	tests := []struct {
		name          string
		ctx           context.Context
		ormConfig     *Config
		queryOptions  *chconn.QueryOptions
		dbName        string
		model         *TestModel
		wantOrmErr    bool
		wantCreateErr bool
		wantDropErr   bool
		wantInsertErr bool
		wantTblErr    bool
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
			wantCreateErr: true,
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
				Name:      "test",
				CreatedAt: time.Now().UTC(),
				UpdatedAt: time.Now().UTC(),
			},
			dbName:        "chorm",
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

			dbBuilder, err := NewCreateDatabase(tt.ctx, orm, tt.dbName, true, tt.queryOptions)
			if tt.wantCreateErr {
				tAssert.Error(err)
				return
			}

			tAssert.NoError(err)
			tAssert.NotNil(dbBuilder)

			tblBuilder, err := NewCreateTable(tt.ctx, orm, tt.model, tt.queryOptions)
			if tt.wantTblErr {
				tAssert.Error(err)
				return
			}

			tAssert.NoError(err)
			tAssert.NotNil(tblBuilder)

			record, builder, err := NewInsert(tt.ctx, orm, tt.model, tt.queryOptions)
			if tt.wantInsertErr {
				tAssert.Error(err)
				return
			}

			tAssert.NoError(err)
			tAssert.NotNil(record)
			tAssert.NotNil(builder)

			t.Logf("Insert SQL: %s", builder.SQL())
			fmt.Printf("response: %+v \n", record)

			/* 			dbBuilder, err = NewDropDatabase(tt.ctx, orm, tt.dbName, tt.queryOptions)
			   			if tt.wantDropErr {
			   				tAssert.Error(err)
			   				return
			   			}

			   			tAssert.NoError(err)
			   			tAssert.NotNil(dbBuilder) */
		})
	}
}

func BenchmarkNewInsert(b *testing.B) {
	ctx := context.TODO()
	ormConfig := &Config{
		Host:     "localhost",
		Port:     9000,
		Username: "default",
		Password: "local12345",
		Database: "unpack",
		Insecure: true,
	}

	orm, err := NewORM(ctx, ormConfig)
	if err != nil {
		b.Fatalf("Failed to create ORM: %v", err)
	}

	if _, err := NewCreateDatabase(ctx, orm, "chorm", true, nil); err != nil {
		b.Fatalf("Failed to create database: %v", err)
	}

	if _, err := NewCreateTable(ctx, orm, &TestModel{}, nil); err != nil {
		b.Fatalf("Failed to create table: %v", err)
	}

	b.ResetTimer() // Reset the timer to exclude setup time

	for i := 0; i < b.N; i++ {
		model := &TestModel{
			Name:      fmt.Sprintf("test_%d", i),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		}

		_, _, err := NewInsert(ctx, orm, model, nil)
		if err != nil {
			b.Fatalf("Failed to insert: %v", err)
		}
	}
}
