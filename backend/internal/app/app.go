package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/Neimess/zorkin-store-project/internal/config"
	repository "github.com/Neimess/zorkin-store-project/internal/repository/psql"
	"github.com/Neimess/zorkin-store-project/internal/service"
	"github.com/Neimess/zorkin-store-project/internal/transport/http/restHTTP"
	"github.com/Neimess/zorkin-store-project/pkg/database/psql"
	"github.com/Neimess/zorkin-store-project/server/rest"

	// httpDelivery "github.com/Neimess/zorkin-store-project/internal/delivery/http"

	"github.com/jmoiron/sqlx"
)

type Application struct {
	cfg    *config.Config
	db     *sqlx.DB
	server *rest.Server
	logger *slog.Logger
}

func NewApplication(dep *Deps) (*Application, error) {
	log := dep.Logger.With(slog.String("component", "app"))
	logNew := log.With(slog.String("op", "app.new"))

	start := time.Now()
	db, err := psql.New(
		context.Background(),
		psql.WithHost(dep.Config.Storage.Host),
		psql.WithUser(dep.Config.Storage.User),
		psql.WithPassword(dep.Config.Storage.Password),
		psql.WithPort(dep.Config.Storage.Port),
		psql.WithConnLifetime(dep.Config.Storage.ConnMaxLifetime),
		psql.WithMaxConns(dep.Config.Storage.MaxOpenConns),
		psql.WithSSL(dep.Config.Storage.SSLMode),
		psql.WithDB(dep.Config.Storage.DBName),
	)
	if err != nil {
		log.Error("db connect failed", slog.Any("error", err))
		return nil, fmt.Errorf("application.dbconnect: %w", err)
	}

	logNew.Info("db connected",
		slog.String("host", dep.Config.Storage.Host),
		slog.Int("max_open_conns", dep.Config.Storage.MaxOpenConns),
		slog.Duration("took", time.Since(start)),
	)

	// 2. Репозитории — отвечают за доступ к данным
	repos := repository.New(repository.Deps{
		DB:     db,
		Logger: dep.Logger,
	})
	logNew.Debug("repositories initialized")

	// 3. Сервисы — бизнес-логика, используют репозиторий

	services := service.New(
		service.Deps{
			Logger:             dep.Logger,
			Config:             dep.Config,
			ProductRepository:  repos.ProductRepository,
			CategoryRepository: repos.CategoryRepository,
		},
	)
	logNew.Debug("services wired")

	// 4. Хендлеры — принимают интерфейсы сервисов (Interface Segregation)
	restHandlers := restHTTP.New(&restHTTP.Deps{
		ProductService:  services.ProductService,
		CategoryService: services.CategoryService,
		Logger:          dep.Logger,
	})
	// 5. Сервер — HTTP-сервер
	srv := rest.NewServer(dep.Config.HTTPServer, restHandlers, dep.Logger)
	logNew.Info("http server constructed",
		slog.String("addr", dep.Config.HTTPServer.Address),
	)
	return &Application{
		cfg:    dep.Config,
		db:     db,
		server: srv,
		logger: log,
	}, nil
}

func (a *Application) Run(ctx context.Context) error {
	const op = "app.app.run"
	log := a.logger.With("op", op)

	log.Info("server started", slog.String("address", a.cfg.HTTPServer.Address))

	if err := a.server.Run(); err != nil && err != http.ErrServerClosed {
		log.Error("server initialization error", slog.Any("error", err))
		return err
	}
	return nil
}

func (a *Application) Shutdown(ctx context.Context) error {
	const op = "app.shutdown"
	log := a.logger.With(slog.String("op", op))

	log.Info("attempting graceful shutdown")

	if err := a.server.Shutdown(ctx); err != nil {
		log.Error("HTTP server shutdown failed", slog.Any("error", err))
	} else {
		log.Info("HTTP server shutdown completed")
	}

	if err := a.db.Close(); err != nil {
		log.Error("DB close failed", slog.Any("error", err))
		return err
	}
	log.Info("shutdown completed successfully")
	return nil

}
