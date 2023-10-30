package dml

import (
	"fmt"
	"strings"

	"github.com/0x19/go-clickhouse-model/models"
)

type QueryType int8

const (
	Select QueryType = iota
	Insert
	Update
	Delete
)

type DmlBuilder struct {
	table      string
	queryType  QueryType
	fields     []string
	subQueries []*DmlBuilder
}

func (d *DmlBuilder) Model(m models.Model) *DmlBuilder {
	d.table = m.TableName()
	return d
}

func (d *DmlBuilder) Select(fields ...interface{}) *DmlBuilder {
	for _, field := range fields {
		switch v := field.(type) {
		case string:
			d.fields = append(d.fields, v)
		case *DmlBuilder:
			d.subQueries = append(d.subQueries, v)
		}
	}
	return d
}

func (d *DmlBuilder) Fields(fields ...string) *DmlBuilder {
	d.fields = fields
	return d
}

func (d *DmlBuilder) Build() (string, error) {
	var queryBuilder strings.Builder

	switch d.queryType {
	case Select:
		queryBuilder.WriteString("SELECT ")

		if len(d.fields) == 0 && len(d.subQueries) == 0 {
			return "", fmt.Errorf("no fields selected or subqueries provided")
		}

		if len(d.fields) > 0 {
			queryBuilder.WriteString(strings.Join(d.fields, ", "))
		}

		for _, subQuery := range d.subQueries {
			subQueryString, err := subQuery.Build()
			if err != nil {
				return "", err
			}
			queryBuilder.WriteString(fmt.Sprintf(", (%s)", subQueryString))
		}

		queryBuilder.WriteString(" FROM ")
		queryBuilder.WriteString(d.table)

		queryBuilder.WriteString(";")

	case Insert:
		if len(d.fields) == 0 {
			return "", fmt.Errorf("no fields selected for insert")
		}

		queryBuilder.WriteString("INSERT INTO ")
		queryBuilder.WriteString(d.table)
		queryBuilder.WriteString(" (")
		queryBuilder.WriteString(strings.Join(d.fields, ", "))
		queryBuilder.WriteString(") VALUES")
	case Update:
		queryBuilder.WriteString("ALTER TABLE ")
		queryBuilder.WriteString(d.table)
	case Delete:
		queryBuilder.WriteString("DELETE FROM ")
		queryBuilder.WriteString(d.table)
	}

	return queryBuilder.String(), nil
}

func (d *DmlBuilder) String() string {
	response, _ := d.Build()
	return response
}
