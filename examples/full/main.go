package main

import (
	"context"
	"time"

	chorm "github.com/0x19/go-clickhouse-orm"
	"github.com/0x19/go-clickhouse-orm/models"
	"github.com/vahid-sohrabloo/chconn/v3"
	"github.com/vahid-sohrabloo/chconn/v3/column"
	"github.com/vahid-sohrabloo/chconn/v3/types"
	"go.uber.org/zap"
)

var (
	dbName = "chorm"
)

type TestModel struct {
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (d *TestModel) TableName() string {
	return d.GetDeclaration().TableName
}

func (d *TestModel) GetNameField() *column.String {
	c := column.NewString()
	c.SetName([]byte("name"))
	c.Append(d.Name)
	return c
}

func (d *TestModel) GetCreatedAtField() *column.Date[types.DateTime] {
	c := column.NewDate[types.DateTime]()
	c.SetName([]byte("created_at"))
	c.Append(d.CreatedAt)
	return c
}

func (d *TestModel) GetUpdatedAtField() *column.Date[types.DateTime] {
	c := column.NewDate[types.DateTime]()
	c.SetName([]byte("updated_at"))
	c.Append(d.UpdatedAt)
	return c
}

func (d *TestModel) Settings() []string {
	return []string{}
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

func (d *TestModel) ScanRow(row chconn.SelectStmt) error {
	row.Rows().Scan(&d.Name, &d.CreatedAt, &d.UpdatedAt)
	return nil
}

func main() {
	logger, err := zap.NewDevelopmentConfig().Build()
	if err != nil {
		panic(err)
	}
	zap.ReplaceGlobals(logger)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	config := chorm.Config{
		Host:     "localhost",
		Port:     9000,
		Username: "default",
		Password: "local12345",
		Database: "unpack",
		Insecure: true,
	}

	orm, err := chorm.NewORM(ctx, &config)
	if err != nil {
		panic(err)
	}

	dbBuilder, err := chorm.NewCreateDatabase(ctx, orm, dbName, true, nil)
	if err != nil {
		panic(err)
	}

	zap.L().Info("Create Database SQL", zap.String("sql", dbBuilder.SQL()))

	tblBuilder, err := chorm.NewCreateTable(ctx, orm, &TestModel{}, nil)
	if err != nil {
		panic(err)
	}

	zap.L().Info("Create Table SQL", zap.String("sql", tblBuilder.SQL()))

	model := &TestModel{
		Name:      "test",
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	record, builder, err := chorm.NewInsert(ctx, orm, model, nil)
	if err != nil {
		panic(err)
	}

	zap.L().Info("Insert SQL", zap.String("sql", builder.SQL()))
	zap.L().Info("Insert Record", zap.Any("record", record))

	tblDropBuilder, err := chorm.NewDropTable(ctx, orm, &TestModel{}, nil)
	if err != nil {
		panic(err)
	}

	zap.L().Info("Drop Table SQL", zap.String("sql", tblDropBuilder.SQL()))

	dbDropBuilder, err := chorm.NewDropDatabase(ctx, orm, dbName, nil)
	if err != nil {
		panic(err)
	}

	zap.L().Info("Drop Database SQL", zap.String("sql", dbDropBuilder.SQL()))
}
