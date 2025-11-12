package database

import (
	"backend-master/configs"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

const (
	driverName = "pgx"
)

type DBManager interface {
	GetDB() *sqlx.DB
}

type dbManagerImpl struct {
	DBManager

	conn *sqlx.DB
}

func NewManager(
	cfg configs.DatabaseConfig,
	logger *zap.Logger,
) (DBManager, error) {
	// urlExample := "postgres://username:password@localhost:5432/database_name"
	url := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		cfg.PgUser,
		cfg.PgPass,
		cfg.PgHost,
		cfg.PgPort,
		cfg.PgDb,
	)

	pgCfg, err := pgx.ParseConfig(url)
	if err != nil {
		return nil, fmt.Errorf("failed to parse postgres config: %w", err)
	}

	pgLog := NewPgxLogger(logger)
	pgCfg.Tracer = pgLog

	nativeDB := stdlib.OpenDB(*pgCfg)

	nativeDB.SetMaxOpenConns(10)
	nativeDB.SetMaxIdleConns(5)

	return &dbManagerImpl{
		conn: sqlx.NewDb(nativeDB, driverName),
	}, nil
}

func (d *dbManagerImpl) GetDB() *sqlx.DB {
	return d.conn
}
