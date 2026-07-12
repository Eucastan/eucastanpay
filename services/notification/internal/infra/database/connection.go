package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	commonconfig "github.com/Eucastan/eucastanpay/common/pkg/config"
	"github.com/Eucastan/eucastanpay/services/notification/config"
	"github.com/Eucastan/eucastanpay/services/notification/migrations"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/sirupsen/logrus"
)

type DBConnect struct {
	DB     *pgxpool.Pool
	sqlDB  *sql.DB
	logger *logrus.Logger
}

func NewPostgresDB(cfg *config.Config, logger *logrus.Logger) *DBConnect {
	pgxConfig, err := pgxpool.ParseConfig(cfg.SharedCfg.Dsn)
	if err != nil {
		logger.WithError(err).Fatal("failed to parse database config")
	}

	pgxConfig.MaxConns = 25
	pgxConfig.MinConns = 5
	pgxConfig.MaxConnIdleTime = 30 * time.Minute
	pgxConfig.MaxConnLifetime = 2 * time.Hour
	pgxConfig.HealthCheckPeriod = 5 * time.Second

	pgxConfig.ConnConfig.RuntimeParams["search_path"] = cfg.SharedCfg.Schema

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, pgxConfig)
	if err != nil {
		logger.WithError(err).Fatal("failed to connect to postgres")
	}

	if err := pool.Ping(ctx); err != nil {
		logger.WithError(err).Fatal(err)
	}

	sqlDB := stdlib.OpenDBFromPool(pool)

	// Run migrations
	if err := runMigrations(sqlDB, cfg); err != nil {
		logger.WithError(err).Fatal("Database migration failed")
	}

	return &DBConnect{
		DB:     pool,
		sqlDB:  sqlDB,
		logger: logger,
	}
}

func runMigrations(db *sql.DB, cfg *config.Config) error {
	goose.SetBaseFS(migrations.FS)

	if err := commonconfig.ValidateSchema(cfg.SharedCfg.Schema); err != nil {
		return err
	}

	_, err := db.Exec(fmt.Sprintf(`
	CREATE SCHEMA IF NOT EXISTS %s;
	SET search_path TO %s;
	`,
		cfg.SharedCfg.Schema,
		cfg.SharedCfg.Schema,
	))
	if err != nil {
		return err
	}

	entries, err := migrations.FS.ReadDir(".")
	if err != nil {
		return err
	}

	for _, e := range entries {
		println("embedded:", e.Name())
	}

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}
	return goose.Up(db, ".")
}

func (c *DBConnect) CloseDB() {
	if c.sqlDB != nil {
		c.sqlDB.Close()
	}
	if c.DB != nil {
		c.DB.Close()
		c.logger.Info("Database connection pool closed successfully")
	}
}
