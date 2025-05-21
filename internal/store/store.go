package store

import (
	"context"
	"fmt"
	"log"
	"log/slog"

	"ppTodolistService/internal/config"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	pool *pgxpool.Pool
	lg   *slog.Logger
}

func MustNew(lg *slog.Logger, cfg *config.Store) *Store {

	connString := fmt.Sprintf(
		`user=%s password=%s host=%s port=%d dbname=%s sslmode=%s pool_max_conns=%d pool_max_conn_lifetime=%s pool_max_conn_idle_time=%s`,
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
		cfg.SSLMode,
		cfg.PoolMaxConns,
		cfg.PoolMaxConnLifetime.String(),
		cfg.PoolMaxConnIdleTime.String(),
	)

	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		log.Fatalf("failed to initialize store: %v\n", err)
	}

	config.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeCacheDescribe

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatalf("failed to initialize store: %v\n", err)
	}
	if err := pool.Ping(context.Background()); err != nil {
		log.Fatalf("failed to initialize store: %v\n", err)
	}
	return &Store{
		pool,
		lg,
	}
}
