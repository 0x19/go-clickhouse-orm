package chorm

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vahid-sohrabloo/chconn/v3"
)

func TestCreateDatabaseBuilder(t *testing.T) {
	tests := []struct {
		name         string
		ctx          context.Context
		ormConfig    *Config
		queryOptions *chconn.QueryOptions
		dbName       string
		wantOrmErr   bool
		wantErr      bool
	}{
		{
			name: "Basic Create Database With No Name",
			ctx:  context.TODO(),
			ormConfig: &Config{
				Host:     "localhost",
				Port:     9000,
				Username: "default",
				Password: "local12345",
				Database: "unpack",
				Insecure: true,
			},
			wantErr: true,
			dbName:  "",
		},
		{
			name: "Basic Create Database With Name",
			ctx:  context.TODO(),
			ormConfig: &Config{
				Host:     "localhost",
				Port:     9000,
				Username: "default",
				Password: "local12345",
				Database: "unpack",
				Insecure: true,
			},
			wantErr: false,
			dbName:  "chorm",
		},
		{
			name: "Basic Create and Drop Database With Name",
			ctx:  context.TODO(),
			ormConfig: &Config{
				Host:     "localhost",
				Port:     9000,
				Username: "default",
				Password: "local12345",
				Database: "unpack",
				Insecure: true,
			},
			wantErr: false,
			dbName:  "chorm",
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

			builder, err := NewCreateDatabase(tt.ctx, orm, tt.dbName, true, tt.queryOptions)
			if tt.wantErr {
				tAssert.Error(err)
				return
			}

			tAssert.NoError(err)
			tAssert.NotNil(builder)

			builder, err = NewDropDatabase(tt.ctx, orm, tt.dbName, tt.queryOptions)
			if tt.wantErr {
				tAssert.Error(err)
				return
			}

			tAssert.NoError(err)
			tAssert.NotNil(builder)
		})
	}
}
