package db

import (
	"context"
	"fmt"
	"simple_gin_server/configs"

	"sync"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	MaxPgPoolConns = 10
	MinPgPoolConns = 2
)

type PgRepoInterface interface {
	Close()
	GetPool() *pgxpool.Pool
}

type PgRepo struct {
	mu   *sync.Mutex
	pool *pgxpool.Pool
}

// NewPgRepo creates a new PostgreSQL repository with properly configured connection pool
func NewPgRepo(ctx context.Context, conf *configs.Config) (*PgRepo, error) {
	// Parse the connection string into a pgxpool.Config
	poolConfig, err := pgxpool.ParseConfig(conf.Db.Dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse DB DSN: %w", err)
	}

	// Configure connection pool settings
	poolConfig.MaxConns = int32(MaxPgPoolConns)
	poolConfig.MinConns = int32(MinPgPoolConns)

	// Configure connection health checks
	poolConfig.HealthCheckPeriod = 1 * time.Minute
	poolConfig.MaxConnLifetime = 1 * time.Hour
	poolConfig.MaxConnIdleTime = 30 * time.Minute

	// Configure connection timeouts
	poolConfig.ConnConfig.ConnectTimeout = 5 * time.Second

	// Create the connection pool with context
	pool, err := pgxpool.ConnectConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Verify the connection
	connCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := pool.Ping(connCtx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &PgRepo{
		mu:   &sync.Mutex{},
		pool: pool,
	}, nil
}

// Close gracefully closes the connection pool
func (r *PgRepo) Close() {
	if r.pool != nil {
		r.pool.Close()
	}
}

// GetPool returns the connection pool (useful for transactions)
func (r *PgRepo) GetPool() *pgxpool.Pool {
	return r.pool
}
