package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	httpapi "github.com/EmotionlessDev/GoEcho/gen/http"
	"github.com/EmotionlessDev/GoEcho/internal/config"
	apiHttp "github.com/EmotionlessDev/GoEcho/internal/delivery/http"
	"github.com/EmotionlessDev/GoEcho/internal/httpx"
	"github.com/EmotionlessDev/GoEcho/internal/logger"
	"github.com/EmotionlessDev/GoEcho/internal/password"
	"github.com/EmotionlessDev/GoEcho/internal/pg"
	userHttp "github.com/EmotionlessDev/GoEcho/internal/user/delivery/http"
	userRepo "github.com/EmotionlessDev/GoEcho/internal/user/repository"
	userUsecase "github.com/EmotionlessDev/GoEcho/internal/user/usecase"
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

	// hasher
	hasher := password.NewBcryptPasswordHasher()

	userRepo := userRepo.NewPostgresUserRepository(pool)
	createUser := userUsecase.NewCreateUser(userRepo, hasher)
	userHandler := userHttp.NewHandler(createUser)

	apiHandler := apiHttp.NewAPIHandler(userHandler)
	// user routes
	strictHandler := httpapi.NewStrictHandler(apiHandler, nil)
	httpapi.HandlerFromMux(strictHandler, router)

	log.InfoContext(ctx, "запускаю goecho-api",
		slog.String("env", cfg.Env),
		slog.String("addr", cfg.HTTP.Addr))

	if err := httpx.NewServer(cfg.HTTP, cfg.App, router, log).Run(ctx); err != nil {
		return err
	}

	log.InfoContext(ctx, "goecho-api остановлен")

	return nil
}
