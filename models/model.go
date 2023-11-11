package models

import (
	"reflect"
	"strings"

	"github.com/vahid-sohrabloo/chconn/v3"
	"github.com/vahid-sohrabloo/chconn/v3/column"
)

type Model interface {
	TableName() string
	Settings() []string
	GetDeclaration() *Declaration
	ScanRow(row chconn.SelectStmt) error
}

type Field struct {
	Index       int16 // Used for sorting and ordering algorithm
	Name        string
	PrimaryKey  bool
	Type        string
	Nullable    bool
	Default     string
	Comment     string
	Column      column.ColumnBasic
	ReflectType reflect.Type
	GoType      interface{}
}

func (f Field) GetDDL() string {
	var ddlBuilder strings.Builder

	ddlBuilder.WriteString(f.Name)
	ddlBuilder.WriteString(" ")
	ddlBuilder.WriteString(f.Type)

	if f.Nullable {
		ddlBuilder.WriteString(" NULL ")
	}

	if f.Default != "" {
		ddlBuilder.WriteString(" DEFAULT ")
		ddlBuilder.WriteString(f.Default)
	}

	if f.Comment != "" {
		ddlBuilder.WriteString(" COMMENT ")
		ddlBuilder.WriteString("'" + f.Comment + "'")
	}

	return ddlBuilder.String()
}
