package postgres

import (
	"database/sql"
	"fmt"

	"github.com/lyracampos/go-clean-architecture/config"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bunotel"
	"go.uber.org/zap"
)

const (
	DuplicateKeyPrefix = "duplicate key value violates unique constraint"
	NoRowsInResultSet  = "no rows in result set"
)

type Client struct {
	log    *zap.SugaredLogger
	DB     *bun.DB
	Config *config.Config
}

func NewClient(log *zap.SugaredLogger, config *config.Config) (*Client, error) {
	sqlDB := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(config.Database.ConnectionString)))
	sqlDB.SetMaxOpenConns(config.Database.MaxOpenConnections)
	sqlDB.SetMaxIdleConns(config.Database.MaxIdleConnections)

	newDB := bun.NewDB(sqlDB, pgdialect.New())
	newDB.AddQueryHook(bunotel.NewQueryHook())

	// 	if conf.Log.Environment == config.LogEnvironmentDevelopment {
	// 		newDB.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))
	// 	}

	if err := newDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to connecto to postgres database: %w", err)
	}

	log.Info("postgres client started successfully")

	return &Client{
		DB:     newDB,
		Config: config,
		log:    log,
	}, nil
}
