package models

import (
	"sort"

	"github.com/vahid-sohrabloo/chconn/v3/column"
)

type Declaration struct {
	DatabaseName string
	TableName    string
	Engine       string
	PartitionBy  string
	Settings     []string
	OrderBy      []string
	Fields       map[string]Field
}

func (d *Declaration) GetField(name string) (Field, bool) {
	field, ok := d.Fields[name]
	return field, ok
}

func (d *Declaration) GetFields() []Field {
	fields := make([]Field, 0, len(d.Fields))

	for _, field := range d.Fields {
		fields = append(fields, field)
	}

	// Sort fields by name
	sort.Slice(fields, func(i, j int) bool {
		return fields[i].Index < fields[j].Index
	})

	return fields
}

func (d *Declaration) GetFieldNames() []string {
	names := make([]string, 0, len(d.Fields))

	for _, field := range d.GetFields() {
		names = append(names, field.Name)
	}

	return names
}

func (d *Declaration) GetColumns() []column.ColumnBasic {
	names := make([]column.ColumnBasic, 0, len(d.Fields))

	for _, field := range d.GetFields() {
		names = append(names, field.Column)
	}

	return names
}

func (d *Declaration) GetPreparedColumns() []column.ColumnBasic {
	names := make([]column.ColumnBasic, 0, len(d.Fields))

	for _, field := range d.GetFields() {
		names = append(names, field.Column)
	}

	return names
}

func (d *Declaration) GetPKFields() []Field {
	fields := make([]Field, 0, len(d.Fields))

	for _, field := range d.GetFields() {
		if field.PrimaryKey {
			fields = append(fields, field)
		}
	}

	return fields
}
