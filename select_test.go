package chorm

import (
	"context"
	"testing"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
	"github.com/vahid-sohrabloo/chconn/v2"
)

func TestSelectBuilder(t *testing.T) {
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
		wantSelectErr bool
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
			wantSelectErr: false,
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

			record, insertBuilder, err := NewInsert(tt.ctx, orm, tt.model, tt.queryOptions)
			if tt.wantInsertErr {
				tAssert.Error(err)
				return
			}

			tAssert.NoError(err)
			tAssert.NotNil(record)
			tAssert.NotNil(insertBuilder)

			t.Logf("Insert SQL: %s", insertBuilder.SQL())

			instance, err := NewSelect[TestModel](tt.ctx, orm, tt.queryOptions)
			if tt.wantSelectErr {
				tAssert.Error(err)
				return
			}

			tAssert.NoError(err)
			tAssert.NotNil(instance)

			var records []TestModel

			instance.Database(tt.dbName)
			instance.Select("*")

			t.Logf("Select SQL: %s", instance.SQL())

			err = instance.Scan(tt.ctx, tt.queryOptions, records)
			tAssert.NoError(err)

			spew.Dump(records)
		})
	}
}
