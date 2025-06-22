package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Neimess/zorkin-store-project/docs"
	"github.com/Neimess/zorkin-store-project/internal/app"
	"github.com/Neimess/zorkin-store-project/internal/config"
	"github.com/Neimess/zorkin-store-project/pkg/args"
	logger "github.com/Neimess/zorkin-store-project/pkg/log"
)

const op = "cmd.main"

var (
	version   = "dev"
	commit    = "none"
	buildDate = "unknown"
)

// @title Zorkin Store API
// @securityDefinitions.apikey BearerAuth
// @in                         header
// @name                       Authorization
// @description                Type **"Bearer <JWT>"** here
func main() {
	// 1. аргументы + конфиг
	cfg := config.MustLoad(args.Parse())
	printVersion()
	initSwagger(cfg.Swagger)
	// 2. логгер (корневой + контекст op)
	root := logger.MustInitLogger(cfg.Env)
	logMain := root.With(slog.String("component", "cmd"),
		slog.String("op", op),
		slog.String("version", cfg.Version),
		slog.String("env", cfg.Env),
	)
	slog.SetDefault(root)
	logMain.Info("main started")

	// 4. DI, создание приложения
	application := mustCreateApp(cfg, root, logMain)

	// 5. graceful-shutdown контекст
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// 6. пуск сервера
	go runApp(ctx, stop, application, logMain)

	// 7. ожидание сигнала и финальный shutdown
	<-ctx.Done()
	gracefulShutdown(application, root, logMain)
}

/* -------- helpers -------- */

func mustCreateApp(cfg *config.Config, root *slog.Logger, log *slog.Logger) *app.Application {
	deps := &app.Deps{Config: cfg, Logger: root}

	application, err := app.NewApplication(deps)
	if err != nil {
		log.Error("application init failed", slog.Any("error", err))
		os.Exit(1)
	}
	log.Info("application initialized")
	return application
}

func runApp(ctx context.Context, cancel context.CancelFunc, a *app.Application, log *slog.Logger) {
	if err := a.Run(ctx); err != nil {
		log.Error("run returned error", slog.Any("error", err))
		cancel()
	}
}

func initSwagger(cfg config.SwaggerInfo) {
	docs.SwaggerInfo.Host = cfg.Host
	docs.SwaggerInfo.Schemes = cfg.Schemes
	docs.SwaggerInfo.Version = cfg.Version
}

func gracefulShutdown(a *app.Application, root, log *slog.Logger) {
	log.With(slog.String("op", "cmd.main.graceful_shutdown")).Info("starting graceful shutdown")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := a.Shutdown(ctx); err != nil {
		root.Warn("graceful shutdown failed", slog.Any("error", err))
	} else {
		root.Info("graceful shutdown complete")
	}
}

func printVersion() {
	fmt.Printf("Version: %s, Commit: %s, Built: %s\n", version, commit, buildDate)
}
