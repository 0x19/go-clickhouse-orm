package chorm

import (
	"time"

	"github.com/0x19/go-clickhouse-orm/models"
	"github.com/vahid-sohrabloo/chconn/v3"
	"github.com/vahid-sohrabloo/chconn/v3/column"
	"github.com/vahid-sohrabloo/chconn/v3/types"
)

type TestModel struct {
	models.Model

	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (d *TestModel) TableName() string {
	return d.GetDeclaration().TableName
}

func (d *TestModel) GetNameField() column.ColumnBasic {
	c := column.NewString()
	c.SetType([]byte("String"))
	c.SetName([]byte("name"))
	c.Append(d.Name)
	return c
}

func (d *TestModel) GetCreatedAtField() column.ColumnBasic {
	c := column.NewDate[types.DateTime]()
	c.SetType([]byte("DateTime"))
	c.SetName([]byte("created_at"))
	c.Append(d.CreatedAt)
	return c
}

func (d *TestModel) GetUpdatedAtField() *column.Date[types.DateTime] {
	c := column.NewDate[types.DateTime]()
	c.SetType([]byte("DateTime"))
	c.SetName([]byte("updated_at"))
	c.Append(d.UpdatedAt)
	return c
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

func (d *TestModel) ScanRow(stmt chconn.SelectStmt) error {
	if err := stmt.Rows().Scan(&d.Name, &d.CreatedAt, &d.UpdatedAt); err != nil {
		return err
	}

	return nil
}
