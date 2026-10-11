// Package pg — доступ к Postgres: пул соединений pgx и запуск миграций goose.
package pg

import (
	"context"
	"fmt"
	"time"

	"github.com/EmotionlessDev/GoEcho/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

const pingTimeout = 2 * time.Second

// NewPool поднимает пул соединений и сразу проверяет, что база отвечает.
// Проверка на старте важна для стенда: сервис с недоступной базой должен
// сказать об этом в логе один раз, а не отдавать 500 на каждый запрос.
func NewPool(ctx context.Context, cfg config.Postgres) (*pgxpool.Pool, error) {
	poolCfg, err := pgxpool.ParseConfig(cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("разбор POSTGRES_DSN: %w", err)
	}

	poolCfg.MaxConns = cfg.MaxConns
	poolCfg.MinConns = cfg.MinConns
	poolCfg.MaxConnLifetime = cfg.MaxConnLifetime
	poolCfg.MaxConnIdleTime = cfg.MaxConnIdleTime
	poolCfg.ConnConfig.ConnectTimeout = cfg.ConnectTimeout

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, fmt.Errorf("пул соединений с Postgres: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()

		return nil, fmt.Errorf("база не отвечает: %w", err)
	}

	return pool, nil
}

// HealthCheck — проверка готовности для /readyz: база отвечает на ping.
func HealthCheck(pool *pgxpool.Pool) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		ctx, cancel := context.WithTimeout(ctx, pingTimeout)
		defer cancel()

		if err := pool.Ping(ctx); err != nil {
			return fmt.Errorf("ping: %w", err)
		}

		return nil
	}
}
