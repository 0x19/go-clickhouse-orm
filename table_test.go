package chorm

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/0x19/go-clickhouse-model/models"
	"github.com/stretchr/testify/assert"
	"github.com/vahid-sohrabloo/chconn/v2"
	"github.com/vahid-sohrabloo/chconn/v2/column"
)

type TestModel struct {
	models.Model

	Name      string
	CreatedAt int64
	UpdatedAt int64
}

func (d *TestModel) TableName() string {
	return d.GetDeclaration().TableName
}

func (d *TestModel) GetDeclaration() *models.Declaration {
	return &models.Declaration{
		DatabaseName: "chorm",
		TableName:    "dummy_model",
		Engine:       "MergeTree()",
		PartitionBy:  "toYYYYMM(created_at)",
		Settings: []string{
			"index_granularity = 8192",
		},
		OrderBy: []string{
			"created_at",
		},
		Fields: map[string]models.Field{
			"name": {
				Index:  0,
				Name:   "name",
				Type:   "String",
				GoType: "string",
				Column: func() column.ColumnBasic {
					c := column.NewString()
					c.SetName([]byte("name"))
					return c
				}(),
			},
			"created_at": {
				Index:      1,
				Name:       "created_at",
				PrimaryKey: true,
				Type:       "DateTime",
				GoType:     "time.Time",
				Column: func() column.ColumnBasic {
					c := column.New[time.Time]()
					c.SetName([]byte("created_at"))
					return c
				}(),
			},
			"updated_at": {
				Index:  2,
				Name:   "updated_at",
				Type:   "DateTime",
				GoType: "time.Time",
				Column: func() column.ColumnBasic {
					c := column.New[time.Time]()
					c.SetName([]byte("updated_at"))
					return c
				}(),
			},
		},
	}
}

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
				Name: "test",
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

			record, builder, err := NewCreateTable(tt.ctx, orm, tt.model, tt.queryOptions)
			if tt.wantInsertErr {
				tAssert.Error(err)
				return
			}

			tAssert.NoError(err)
			tAssert.NotNil(record)
			tAssert.NotNil(builder)

			fmt.Println("Create SQL: ", builder.SQL())

			builder, err = NewDropTable(tt.ctx, orm, tt.model, tt.queryOptions)
			if tt.wantInsertErr {
				tAssert.Error(err)
				return
			}

			tAssert.NoError(err)
			tAssert.NotNil(record)
			tAssert.NotNil(builder)

			fmt.Println("Create SQL: ", builder.SQL())

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
