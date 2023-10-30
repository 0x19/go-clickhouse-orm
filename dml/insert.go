package dml

type InsertBuilder struct {
	DmlBuilder // Embedded struct
}

func NewInsertBuilder() *InsertBuilder {
	return &InsertBuilder{
		DmlBuilder: DmlBuilder{
			queryType:  Insert,
			subQueries: []*DmlBuilder{},
		},
	}
}
