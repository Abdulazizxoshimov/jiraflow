package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/jira-backend/jiraflow-backend/internal/pkg/config"
	"github.com/jira-backend/jiraflow-backend/internal/pkg/logger"
)

type Postgres struct {
	Builder squirrel.StatementBuilderType
	DB      *pgxpool.Pool
}

// New opens a pgxpool connection using config, pings the database, and starts
// a background monitor goroutine that exits when ctx is cancelled.
func New(ctx context.Context, cfg *config.Config, log logger.Logger) (*Postgres, error) {
	poolCfg, err := pgxpool.ParseConfig(cfg.Postgres.DSN())
	if err != nil {
		return nil, fmt.Errorf("postgres: parse config: %w", err)
	}

	poolCfg.MinConns = cfg.Postgres.MinConns
	poolCfg.MaxConns = cfg.Postgres.MaxConns
	poolCfg.MaxConnIdleTime = cfg.Postgres.MaxConnIdleTime
	poolCfg.MaxConnLifetime = cfg.Postgres.MaxConnLifetime
	poolCfg.HealthCheckPeriod = cfg.Postgres.HealthCheckPeriod
	poolCfg.MaxConnLifetimeJitter = 15 * time.Minute
	poolCfg.ConnConfig.ConnectTimeout = 10 * time.Second

	pingCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(pingCtx, poolCfg)
	if err != nil {
		return nil, fmt.Errorf("postgres: create pool: %w", err)
	}

	if err := pool.Ping(pingCtx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("postgres: ping: %w", err)
	}

	go monitor(ctx, pool, uint32(poolCfg.MaxConns), 30*time.Second, log)

	return &Postgres{
		DB:      pool,
		Builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}, nil
}

// Close gracefully closes the connection pool.
func (p *Postgres) Close() {
	if p.DB != nil {
		p.DB.Close()
	}
}

func monitor(ctx context.Context, pool *pgxpool.Pool, maxConns uint32, interval time.Duration, log logger.Logger) {
	t := time.NewTicker(interval)
	defer t.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			s := pool.Stat()
			log.Debug(ctx, "pgxpool stats",
				logger.Int("acquired", int(s.AcquiredConns())),
				logger.Int("idle", int(s.IdleConns())),
				logger.Int("total", int(s.TotalConns())),
				logger.Int("max", int(maxConns)),
				logger.Int64("acquire_count", s.AcquireCount()),
				logger.Int64("empty_acquire_count", s.EmptyAcquireCount()),
				logger.Int64("canceled_acquire_count", s.CanceledAcquireCount()),
				logger.Int64("acquire_duration_ms", s.AcquireDuration().Milliseconds()),
			)
		}
	}
}
