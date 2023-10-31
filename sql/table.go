package sql

type TableBuilder struct {
	DmlBuilder // Embedded struct
}

func NewCreateTableBuilder() *TableBuilder {
	return &TableBuilder{
		DmlBuilder: DmlBuilder{
			queryType:  CreateTable,
			subQueries: []*DmlBuilder{},
		},
	}
}
