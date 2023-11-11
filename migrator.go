package chorm

import (
	"context"
	"fmt"
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
	orm  *ORM
	up   map[string]func(ctx context.Context, orm *ORM, migrator *Migrator) error
	down map[string]func(ctx context.Context, orm *ORM, migrator *Migrator) error
}

func NewMigrator(orm *ORM) (*Migrator, error) {
	// We need to first register the migration model for the ORM to function correctly.
	if err := orm.GetManager().RegisterModel(&Migration{}); err != nil {
		return nil, err
	}

	return &Migrator{
		orm:  orm,
		up:   make(map[string]func(ctx context.Context, orm *ORM, migrator *Migrator) error),
		down: make(map[string]func(ctx context.Context, orm *ORM, migrator *Migrator) error),
	}, nil
}

func (m *Migrator) GetUpMigrations() map[string]func(ctx context.Context, orm *ORM, migrator *Migrator) error {
	return m.up
}

func (m *Migrator) GetDownMigrations() map[string]func(ctx context.Context, orm *ORM, migrator *Migrator) error {
	return m.up
}

func (m *Migrator) GetUpMigrationNames() []string {
	toReturn := make([]string, 0, len(m.up))
	for name := range m.up {
		toReturn = append(toReturn, name)
	}
	return toReturn
}

func (m *Migrator) GetDownMigrationNames() []string {
	toReturn := make([]string, 0, len(m.down))
	for name := range m.down {
		toReturn = append(toReturn, name)
	}
	return toReturn
}

func (m *Migrator) RegisterMigration(
	name string,
	up func(ctx context.Context, orm *ORM, migrator *Migrator) error,
	down func(ctx context.Context, orm *ORM, migrator *Migrator) error,
) error {
	if _, ok := m.up[name]; ok {
		return fmt.Errorf("migration with name %s already exists", name)
	}

	m.up[name] = up
	m.down[name] = down
	return nil
}

func (m *Migrator) Migrate(ctx context.Context, queryOptions *chconn.QueryOptions) error {
	for name, up := range m.up {
		// Check if the migration has already been applied.

		// Apply the migration.
		if err := up(ctx, m.orm, m); err != nil {
			return err
		}

		// Create a new migration record.
		migration := &Migration{
			UUID:      uuid.New(),
			Name:      name,
			Migrated:  true,
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		}

		if _, _, err := m.orm.Insert(ctx, migration, queryOptions); err != nil {
			return err
		}

		zap.L().Info("Successfully migrated", zap.String("migration", name))
	}

	return nil
}

func (m *Migrator) Rollback(ctx context.Context, queryOptions *chconn.QueryOptions) error {
	return nil
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
