package models

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
)

type TestModel struct {
	Model       `clickhouse:"table:test_model, engine:MergeTree(), partition:toYYYYMM(created_at), order:created_at"`
	ID          int64  `clickhouse:"name:id, primary:true, type:Int64, nullable:false, default:0, comment:'Unique ID'"`
	Name        string `clickhouse:"name:name"`
	Description string `clickhouse:"name:description, type:String, default:'No description'"`
}

func (t *TestModel) Settings() []string {
	return []string{
		"index_granularity=8192",
	}
}

func (t *TestModel) TableName() string {
	return "test_model"
}

func TestModelManager(t *testing.T) {
	tAssert := assert.New(t)

	cfg := &ManagerConfig{
		DatabaseName: "test",
	}

	manager := NewManager(cfg)
	tAssert.NotNil(manager)

	err := manager.RegisterModel(&TestModel{})
	tAssert.NoError(err)

	declaration, err := manager.GetDeclaration(&TestModel{})
	tAssert.NoError(err)

	spew.Dump(declaration)
}
