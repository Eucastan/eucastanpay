package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/Eucastan/eucastanpay/services/admin/config"
	"github.com/Eucastan/eucastanpay/services/admin/migrations"
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
	pgxConfig, err := pgxpool.ParseConfig(cfg.Dsn)
	if err != nil {
		logger.WithError(err).Warn("failed to parse database config")
	}

	pgxConfig.MaxConns = 25
	pgxConfig.MinConns = 5
	pgxConfig.MaxConnIdleTime = 30 * time.Minute
	pgxConfig.MaxConnLifetime = 30 * time.Hour
	pgxConfig.HealthCheckPeriod = 5 * time.Second

	pool, err := pgxpool.NewWithConfig(context.Background(), pgxConfig)
	if err != nil {
		logger.WithError(err).Error("failed to connect to postgres")
	}

	sqlDB := stdlib.OpenDBFromPool(pool)

	// Run migrations
	if err := runMigrations(sqlDB); err != nil {
		logger.WithError(err).Fatal("Database migration failed")
	}

	return &DBConnect{
		DB:     pool,
		sqlDB:  sqlDB,
		logger: logger,
	}
}

func runMigrations(db *sql.DB) error {
	goose.SetBaseFS(migrations.FS)

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
