package chorm

type Migrator struct {
	orm *ORM
}

func NewMigrator(orm *ORM) *Migrator {
	return &Migrator{orm: orm}
}
