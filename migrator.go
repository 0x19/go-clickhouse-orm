package chorm

import (
	"context"
	"time"

	"github.com/0x19/go-clickhouse-orm/models"
	"github.com/google/uuid"
	"github.com/vahid-sohrabloo/chconn/v3"
	"go.uber.org/zap"
)

type Migration struct {
	models.Model `clickhouse:"table:migrations, engine:MergeTree(), order: uuid"`

	UUID      uuid.UUID `clickhouse:"name:uuid, type:UUID, primary:true"`
	Name      string    `clickhouse:"name:name, type:String"`
	Migrated  bool      `clickhouse:"name:migrated, type:Boolean"`
	CreatedAt time.Time `clickhouse:"name:created_at, type:DateTime"`
	UpdatedAt time.Time `clickhouse:"name:updated_at, type:DateTime"`
}

func (m *Migration) TableName() string {
	return "migrations"
}

func (m *Migration) Settings() []string {
	return []string{}
}

type Migrator struct {
	orm *ORM
}

func NewMigrator(orm *ORM) (*Migrator, error) {

	// We need to first register the migration model for the ORM to function correctly.
	if err := orm.GetManager().RegisterModel(&Migration{}); err != nil {
		return nil, err
	}

	return &Migrator{orm: orm}, nil
}

func (m *Migrator) Setup(ctx context.Context, queryOptions *chconn.QueryOptions) error {
	if _, err := m.orm.CreateDatabase(ctx, m.orm.GetDatabaseName(), true, queryOptions); err != nil {
		return err
	}

	zap.L().Info("Successfully created database", zap.String("database", m.orm.GetDatabaseName()))

	if _, err := NewCreateTable(ctx, m.orm, &Migration{}, queryOptions); err != nil {
		return err
	}

	zap.L().Info("Successfully created migrations table")

	return nil
}

func (m *Migrator) Destroy(ctx context.Context, queryOptions *chconn.QueryOptions) error {
	if _, err := NewDropTable(ctx, m.orm, &Migration{}, queryOptions); err != nil {
		return err
	}

	zap.L().Info("Successfully created migrations table")

	return nil
}
