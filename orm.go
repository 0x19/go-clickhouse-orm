package chorm

import (
	"context"
	"time"

	"github.com/0x19/go-clickhouse-orm/models"
	"github.com/vahid-sohrabloo/chconn/v3"
	"github.com/vahid-sohrabloo/chconn/v3/chpool"
)

type ORM struct {
	ctx      context.Context
	cfg      *Config
	db       chpool.Pool
	manager  *models.Manager
	migrator *Migrator
}

func NewORM(ctx context.Context, cfg *Config) (*ORM, error) {
	if cfg == nil {
		return nil, ErrNoConfigProvided
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	manager := models.NewManager(&models.ManagerConfig{
		DatabaseName: cfg.Database,
	})

	toReturn := &ORM{ctx: ctx, cfg: cfg, manager: manager}

	migrator, err := NewMigrator(toReturn)
	if err != nil {
		return nil, err
	}
	toReturn.migrator = migrator

	if err := toReturn.Connect(); err != nil {
		return nil, err
	}

	return toReturn, nil
}

func (o *ORM) GetMigrator() *Migrator {
	return o.migrator
}

func (o *ORM) GetManager() *models.Manager {
	return o.manager
}

func (o *ORM) GetDatabaseName() string {
	return o.cfg.Database
}

func (o *ORM) Connect() error {
	config, err := chpool.ParseConfig(o.cfg.GetDSN())
	if err != nil {
		return err
	}

	config.ConnConfig.ClientName = "unpack-services"
	config.ConnConfig.Compress = chconn.CompressLZ4
	config.MaxConns = 10
	config.MaxConnLifetime = time.Duration(10) * time.Minute
	config.MaxConnIdleTime = time.Duration(10) * time.Minute
	config.MinConns = 5

	db, err := chpool.NewWithConfig(config)
	if err != nil {
		return err
	}
	o.db = db

	return o.db.Ping(o.ctx)
}

func (o *ORM) GetContext() context.Context {
	return o.ctx
}

func (o *ORM) GetConfig() *Config {
	return o.cfg
}

func (o *ORM) CreateDatabase(ctx context.Context, dbName string, useDb bool, queryOptions *chconn.QueryOptions) (*DatabaseBuilder, error) {
	return NewCreateDatabase(ctx, o, dbName, useDb, queryOptions)
}

func (o *ORM) DropDatabase(ctx context.Context, dbName string, queryOptions *chconn.QueryOptions) (*DatabaseBuilder, error) {
	return NewDropDatabase(ctx, o, dbName, queryOptions)
}

func (o *ORM) Insert(ctx context.Context, model models.Model, queryOptions *chconn.QueryOptions) (models.Model, *InsertBuilder[models.Model], error) {
	return NewInsert(ctx, o, model, queryOptions)
}

func (o *ORM) GetConn() chpool.Pool {
	return o.db
}

func (o *ORM) Close() {
	o.db.Close()
}
