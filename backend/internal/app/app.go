package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/Neimess/zorkin-store-project/internal/config"
	repository "github.com/Neimess/zorkin-store-project/internal/infrastructure"
	"github.com/Neimess/zorkin-store-project/internal/server/rest"
	"github.com/Neimess/zorkin-store-project/internal/service"
	"github.com/Neimess/zorkin-store-project/internal/transport/http/restHTTP"
	"github.com/Neimess/zorkin-store-project/pkg/database/psql"
	"github.com/Neimess/zorkin-store-project/pkg/secret/jwt"

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
	fmt.Println(dep.Config)
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
	repos, err := repository.New(repository.Deps{
		DB:     db,
		Logger: dep.Logger,
	})
	if err != nil {
		logNew.Error("repositories initialization failed", slog.Any("error", err))
	}
	logNew.Debug("repositories initialized")

	// 3. Сервисы — бизнес-логика, используют репозиторий

	// 4. Генераторы
	jwtGenerator := jwt.NewGenerator(jwt.JWTConfig{
		Secret:    dep.Config.JWTConfig.JWTSecret,
		Issuer:    dep.Config.JWTConfig.Issuer,
		Audience:  dep.Config.JWTConfig.Audience,
		Algorithm: dep.Config.JWTConfig.Algorithm,
	})

	services, err := service.New(
		service.NewDeps(
			jwtGenerator,
			dep.Logger,
			repos.ProductRepository,
			repos.CategoryRepository,
			repos.PresetRepository,
			repos.AttributeRepository,
			repos.CoefficientRepository,
			repos.ServiceRepository,
		),
	)
	if err == nil {
		logNew.Debug("services wired")
	} else {
		log.Error("services initialization failed", slog.Any("error", err))
		return nil, fmt.Errorf("application.services: %w", err)
	}

	handlersDeps, err := restHTTP.NewDeps(
		dep.Logger,
		services.ProductService,
		services.CategoryService,
		services.AuthService,
		services.PresetService,
		services.AttributeService,
		services.CoefficientService,
		services.ServiceService,
	)
	if err != nil {
		logNew.Error("handlers dependencies initialization failed", slog.Any("error", err))
		return nil, fmt.Errorf("application.handlersdeps: %w", err)
	}
	// 4. Хендлеры — принимают интерфейсы сервисов (Interface Segregation)
	restHandlers, err := restHTTP.New(handlersDeps)
	if err == nil {
		logNew.Debug("handlers initialized")
	} else {
		log.Error("handlers initialization failed", slog.Any("error", err))
		return nil, fmt.Errorf("application.handlers: %w", err)
	}
	// 5. Сервер — HTTP-сервер
	depsServer, err := rest.NewDeps(
		dep.Config,
		restHandlers,
		dep.Logger,
	)
	if err != nil {
		logNew.Error("server dependencies initialization failed", slog.Any("error", err))
		return nil, fmt.Errorf("application.serverdeps: %w", err)
	}
	srv, err := rest.NewServer(depsServer)
	if err != nil {
		logNew.Error("server initialization failed", slog.Any("error", err))
		return nil, fmt.Errorf("application.server: %w", err)
	}

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

func (a *Application) ServerHandler() http.Handler {
	return a.server.Handler()
}

func (a *Application) DB() *sqlx.DB {
	return a.db
}
