package app

import (
	"context"
	"errors"
	"fmt"
	log "log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"lost-items-service/internal/app/pkg/postgres"
	"lost-items-service/internal/app/pkg/redis"
	"lost-items-service/internal/client/cache"
	"lost-items-service/internal/client/cache/redis"
	"lost-items-service/internal/config"
	cfg "lost-items-service/internal/config/env"
	"lost-items-service/internal/db/pgxdb"
	"lost-items-service/internal/handler"
	"lost-items-service/internal/repository"
	"lost-items-service/internal/service"
	"lost-items-service/pkg/logger"
)

type App struct {
	httpCfg     config.HTTPConfig
	db          pgxdb.DB
	redisClient cache.RedisClient
	router      *chi.Mux
}

func NewApp(ctx context.Context) (*App, error) {
	var (
		envPath    = ".env"
		configPath = "./configs/config.yaml"
	)

	logger := logger.InitLogger()

	if err := config.LoadEnv(envPath); err != nil {
		return nil, fmt.Errorf("error loading env file, %s", envPath)
	}

	pgCfg, err := cfg.PGConfigLoad(configPath)
	if err != nil {
		return nil, fmt.Errorf("error loading postgres config: %w", err)
	}

	log.Info(pgCfg.DSN())
	htppCfg, err := cfg.HTTPConfigLoad(configPath)
	if err != nil {
		return nil, fmt.Errorf("error loading http config: %w", err)
	}

	jwtCfg, err := cfg.JWTConfigLoad()
	if err != nil {
		return nil, fmt.Errorf("error loading jwt config: %w", err)
	}

	redisCfg, err := cfg.RedisConfigLoad(configPath)
	if err != nil {
		return nil, fmt.Errorf("error loading redis config: %w", err)
	}

	pgPool, err := postgres.InitDBPool(ctx, pgCfg)
	if err != nil {
		return nil, fmt.Errorf("error initializing DB pool: %w", err)
	}

	redisPool, err := redis.InitRedisPool(ctx, redisCfg)
	if err != nil {
		return nil, fmt.Errorf("error initializing Redis pool: %w", err)
	}

	redisClient := redisclient.NewClient(redisPool, redisCfg)

	pgDB := pgxdb.NewPgxDB(pgPool)
	//init repo

	repo := repository.NewRepository(pgDB, redisClient)

	// init service
	serv := service.NewService(repo, jwtCfg.Jwt)

	//init router
	r := handler.NewRouter(serv, jwtCfg.Jwt, logger)

	return &App{
			router:      r,
			httpCfg:     htppCfg,
			db:          pgDB,
			redisClient: redisClient,
		},
		nil
}

func (a *App) Run() error {
	defer func() {
		a.db.Close()
		a.redisClient.Close()
	}()

	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", a.httpCfg.GetPort()),
		Handler:      a.router,
		ReadTimeout:  a.httpCfg.GetTimeout(),
		WriteTimeout: a.httpCfg.GetTimeout(),
		IdleTimeout:  a.httpCfg.GetIdleTimeout(),
	}

	// Запуск сервера
	go func() {
		log.Info("Starting HTTP server", "addr", server.Addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("HTTP server ListenAndServe failed", log.Any("err", err))
		}
	}()

	// Слушаем сигналы остановки
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Info("Shutdown signal received")

	// Контекст с таймаутом на graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Error("Server shutdown failed", log.Any("err", err))
		return err
	}

	select {
	case <-ctx.Done():
		log.Warn("Shutdown timeout exceeded")
	default:
		log.Info("Server exited gracefully")
	}

	return nil
}
