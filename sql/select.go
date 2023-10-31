package sql

type SelectBuilder struct {
	DmlBuilder // Embedded struct
}

func NewSelectBuilder() *SelectBuilder {
	return &SelectBuilder{
		DmlBuilder: DmlBuilder{
			queryType:  Select,
			subQueries: []*DmlBuilder{},
		},
	}
}
