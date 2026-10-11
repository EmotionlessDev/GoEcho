package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/EmotionlessDev/GoEcho/internal/config"
	"github.com/EmotionlessDev/GoEcho/internal/httpx"
	"github.com/EmotionlessDev/GoEcho/internal/logger"
	"github.com/EmotionlessDev/GoEcho/internal/pg"
)

func main() {
	if err := run(); err != nil {
		// Логгер на этот момент может быть ещё не собран (упали на конфиге),
		// поэтому пишем в stderr напрямую.
		fmt.Fprintf(os.Stderr, "goecho-api: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.LoadGoEchoAPI()
	if err != nil {
		return err
	}

	log := logger.New(cfg.App)

	// SIGTERM приходит от `docker compose down` и от оркестратора;
	// SIGINT — от Ctrl+C при локальном запуске.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	// Пул поднимается до сервера: сервис без базы отдавать каталог не может,
	// и узнать об этом лучше на старте, а не на первом запросе.
	pool, err := pg.NewPool(ctx, cfg.Postgres)
	if err != nil {
		return err
	}

	defer pool.Close()

	health := httpx.NewHealth()
	health.Register("postgres", pg.HealthCheck(pool))

	router := httpx.NewRouter(log, health)

	log.InfoContext(ctx, "запускаю goecho-api",
		slog.String("env", cfg.Env),
		slog.String("addr", cfg.HTTP.Addr))

	if err := httpx.NewServer(cfg.HTTP, cfg.App, router, log).Run(ctx); err != nil {
		return err
	}

	log.InfoContext(ctx, "goecho-api остановлен")

	return nil
}
