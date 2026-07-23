package checks

import (
	"context"
	"time"

	"github.com/Eucastan/eucastanpay/common/pkg/healthcheck"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Database struct {
	db *pgxpool.Pool
}

func NewDatabase(db *pgxpool.Pool) *Database {
	return &Database{
		db: db,
	}
}

func (d *Database) Name() string {
	return "database"
}

func (d *Database) Check(ctx context.Context) healthcheck.Component {

	started := time.Now()

	ctx, cancel := context.WithTimeout(
		ctx,
		3*time.Second,
	)
	defer cancel()

	if err := d.db.Ping(ctx); err != nil {
		return healthcheck.Component{
			Name:     d.Name(),
			Status:   healthcheck.Unhealthy,
			Error:    err.Error(),
			Duration: time.Since(started).String(),
		}
	}

	stats := d.db.Stat()
	return healthcheck.Component{
		Name:     d.Name(),
		Status:   healthcheck.Healthy,
		Duration: time.Since(started).String(),
		Details: map[string]interface{}{
			"total_connections":        stats.TotalConns(),
			"idle_connections":         stats.IdleConns(),
			"acquired_connections":     stats.AcquiredConns(),
			"constructing_connections": stats.ConstructingConns(),
			"max_connections":          stats.MaxConns(),
		},
	}
}
