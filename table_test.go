package chorm

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vahid-sohrabloo/chconn/v3"
)

func TestCreateTableBuilder(t *testing.T) {
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

			builder, err := NewCreateTable(tt.ctx, orm, tt.model, tt.queryOptions)
			if tt.wantInsertErr {
				tAssert.Error(err)
				return
			}

			tAssert.NoError(err)
			tAssert.NotNil(builder)

			fmt.Println("Create SQL: ", builder.SQL())

			builder, err = NewDropTable(tt.ctx, orm, tt.model, tt.queryOptions)
			if tt.wantInsertErr {
				tAssert.Error(err)
				return
			}

			tAssert.NoError(err)
			tAssert.NotNil(builder)

			fmt.Println("Drop SQL: ", builder.SQL())

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
