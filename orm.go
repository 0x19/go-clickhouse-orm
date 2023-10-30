package gchm

import (
	"context"
	"time"

	"github.com/vahid-sohrabloo/chconn/v2"
	"github.com/vahid-sohrabloo/chconn/v2/chpool"
)

type ORM struct {
	ctx context.Context
	cfg *Config
	db  chpool.Pool
}

func NewORM(ctx context.Context, cfg *Config) (*ORM, error) {
	if cfg == nil {
		return nil, ErrNoConfigProvided
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	toReturn := &ORM{ctx: ctx, cfg: cfg}

	if err := toReturn.Connect(); err != nil {
		return nil, err
	}

	return toReturn, nil
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

func (o *ORM) GetConn() chpool.Pool {
	return o.db
}

func (o *ORM) Close() {
	o.db.Close()
}
