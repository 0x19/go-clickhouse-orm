package chorm

import (
	"time"

	"github.com/0x19/go-clickhouse-orm/models"
	"github.com/vahid-sohrabloo/chconn/v2/column"
	"github.com/vahid-sohrabloo/chconn/v2/types"
)

type TestModel struct {
	models.Model

	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (d TestModel) TableName() string {
	return d.GetDeclaration().TableName
}

func (d TestModel) GetNameField() *column.String {
	c := column.NewString()
	c.SetName([]byte("name"))
	c.Append(d.Name)
	return c
}

func (d TestModel) GetCreatedAtField() *column.Date[types.DateTime] {
	c := column.NewDate[types.DateTime]()
	c.SetName([]byte("created_at"))
	c.Append(d.CreatedAt)
	return c
}

func (d TestModel) GetUpdatedAtField() *column.Date[types.DateTime] {
	c := column.NewDate[types.DateTime]()
	c.SetName([]byte("updated_at"))
	c.Append(d.UpdatedAt)
	return c
}

func (d TestModel) GetDeclaration() *models.Declaration {
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
				Column: d.GetNameField(),
			},
			"created_at": {
				Index:      1,
				Name:       "created_at",
				PrimaryKey: true,
				Type:       "DateTime",
				GoType:     "time.Time",
				Column:     d.GetCreatedAtField(),
			},
			"updated_at": {
				Index:  2,
				Name:   "updated_at",
				Type:   "DateTime",
				GoType: "time.Time",
				Column: d.GetUpdatedAtField(),
			},
		},
	}
}
