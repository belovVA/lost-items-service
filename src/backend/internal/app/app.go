package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"log/slog"
	"lost-items-service/config"
	"lost-items-service/pkg"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type App struct {
	//pgCfg   config.PGConfig
	httpCfg config.HTTPConfig
	//
	//jwtSecret string
	//
	//pgDB *pgxpool.Pool
	//log *log.Slog

	router *chi.Mux
}

func NewApp(ctx context.Context) (*App, error) {
	pgCfg, err := config.PGConfigLoad()
	if err != nil {
		return nil, fmt.Errorf("error loading postgres config: %w", err)
	}

	_, err = config.HTTPConfigLoad()
	if err != nil {
		return nil, fmt.Errorf("error loading http config: %w", err)
	}

	_, err = config.JWTConfigLoad()
	if err != nil {
		return nil, fmt.Errorf("error loading jwt config: %w", err)
	}

	_, err = pkg.InitDBPool(ctx, pgCfg)
	if err != nil {
		return nil, fmt.Errorf("error initializing DB pool: %w", err)
	}

	//init logger
	//logger_ := slog

	// init repo
	//repo := repository.NewRepository(DBPool)

	// init service
	//service := service.NewService(repo, jwtSecret)

	//init router
	//r := handler.NewRouter(service_, logger_, jwtSecret)
	//return &App{
	//	router: r,
	//}

	return &App{
			router: nil,
		},
		nil
}

func (a *App) Run() error {
	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", a.httpCfg.GetPort()),
		Handler:      a.router,
		ReadTimeout:  a.httpCfg.GetTimeout(),
		WriteTimeout: a.httpCfg.GetTimeout(),
		IdleTimeout:  a.httpCfg.GetIdleTimeout(),
	}

	// Запуск сервера
	go func() {
		slog.Info("Starting HTTP server", "addr", server.Addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("HTTP server ListenAndServe failed", slog.Any("err", err))
		}
	}()

	// Слушаем сигналы остановки
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	slog.Info("Shutdown signal received")

	// Контекст с таймаутом на graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("Server shutdown failed", slog.Any("err", err))
		return err
	}

	select {
	case <-ctx.Done():
		slog.Warn("Shutdown timeout exceeded")
	default:
		slog.Info("Server exited gracefully")
	}

	return nil
}
