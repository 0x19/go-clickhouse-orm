package sql

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
	CreateDatabase
	DropDatabase
	CreateTable
	DropTable
	AlterTable
)

type DmlBuilder struct {
	table       string
	database    string
	queryType   QueryType
	engine      string
	partitionBy string
	primaryKeys []string
	orderBy     []string
	settings    []string
	fields      []string
	subQueries  []*DmlBuilder
}

func (d *DmlBuilder) Model(m models.Model) *DmlBuilder {
	d.table = m.TableName()
	return d
}

func (d *DmlBuilder) Database(database string) *DmlBuilder {
	d.database = database
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

func (d *DmlBuilder) Engine(engine string) *DmlBuilder {
	d.engine = engine
	return d
}

func (d *DmlBuilder) PartitionBy(partitionBy string) *DmlBuilder {
	d.partitionBy = partitionBy
	return d
}

func (d *DmlBuilder) PrimaryKeys(fields ...string) *DmlBuilder {
	d.primaryKeys = fields
	return d
}

func (d *DmlBuilder) OrderBy(fields ...string) *DmlBuilder {
	d.orderBy = fields
	return d
}

func (d *DmlBuilder) Settings(settings ...string) *DmlBuilder {
	d.settings = settings
	return d
}

func (d *DmlBuilder) Build() (string, error) {
	var queryBuilder strings.Builder

	switch d.queryType {
	case CreateDatabase:
		queryBuilder.WriteString("CREATE DATABASE IF NOT EXISTS ")
		queryBuilder.WriteString(d.database)
		queryBuilder.WriteString(";")
	case DropDatabase:
		queryBuilder.WriteString("DROP DATABASE IF EXISTS ")
		queryBuilder.WriteString(d.database)
		queryBuilder.WriteString(";")
	case CreateTable:
		queryBuilder.WriteString("CREATE TABLE IF NOT EXISTS ")
		queryBuilder.WriteString(d.database + "." + d.table)
		queryBuilder.WriteString(" (")
		queryBuilder.WriteString(strings.Join(d.fields, ", "))
		queryBuilder.WriteString(") ")

		if d.engine != "" {
			queryBuilder.WriteString("ENGINE = ")
			queryBuilder.WriteString(d.engine)
			queryBuilder.WriteString(" ")
		}

		if d.partitionBy != "" {
			queryBuilder.WriteString("PARTITION BY ")
			queryBuilder.WriteString(d.partitionBy)
			queryBuilder.WriteString(" ")
		}

		if len(d.primaryKeys) > 0 {
			queryBuilder.WriteString("PRIMARY KEY ")

			if len(d.primaryKeys) == 1 {
				queryBuilder.WriteString(d.primaryKeys[0])
				queryBuilder.WriteString(" ")
			} else {
				queryBuilder.WriteString("(")
				for i, field := range d.primaryKeys {
					queryBuilder.WriteString(field)
					if i < len(d.primaryKeys)-1 {
						queryBuilder.WriteString(", ")
					}
				}
				queryBuilder.WriteString(") ")
			}
		}

		if len(d.orderBy) > 0 {
			queryBuilder.WriteString("ORDER BY ")

			if len(d.orderBy) == 1 {
				queryBuilder.WriteString(d.orderBy[0])
				queryBuilder.WriteString(" ")
			} else {
				queryBuilder.WriteString("(")
				for i, field := range d.orderBy {
					queryBuilder.WriteString(field)
					if i < len(d.orderBy)-1 {
						queryBuilder.WriteString(", ")
					}
				}
				queryBuilder.WriteString(") ")
			}
		}

		if len(d.settings) > 0 {
			queryBuilder.WriteString("SETTINGS ")
			queryBuilder.WriteString(strings.Join(d.settings, ", "))
		}

		queryBuilder.WriteString(";")
	case DropTable:
		queryBuilder.WriteString("DROP TABLE IF EXISTS ")
		queryBuilder.WriteString(d.database + "." + d.table)
		queryBuilder.WriteString(";")
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
