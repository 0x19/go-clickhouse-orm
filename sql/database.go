package sql

type DatabaseBuilder struct {
	DmlBuilder // Embedded struct
}

func NewCreateDatabaseBuilder() *DatabaseBuilder {
	return &DatabaseBuilder{
		DmlBuilder: DmlBuilder{
			queryType:  CreateDatabase,
			subQueries: []*DmlBuilder{},
		},
	}
}

func NewDropDatabaseBuilder() *DatabaseBuilder {
	return &DatabaseBuilder{
		DmlBuilder: DmlBuilder{
			queryType:  DropDatabase,
			subQueries: []*DmlBuilder{},
		},
	}
}
