package sql

type UpdateBuilder struct {
	DmlBuilder
}

func NewUpdateBuilder() *UpdateBuilder {
	return &UpdateBuilder{
		DmlBuilder: DmlBuilder{
			queryType:  Update,
			subQueries: []*DmlBuilder{},
		},
	}
}
