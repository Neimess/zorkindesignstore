package rest

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/Neimess/zorkin-store-project/internal/config"
	"github.com/Neimess/zorkin-store-project/internal/transport/http/restHTTP"
	route "github.com/Neimess/zorkin-store-project/internal/transport/http/routes"
	"github.com/go-chi/chi/v5"
)

type Deps struct {
	cfg      *config.Config
	handlers *restHTTP.Handlers
	log      *slog.Logger
}

func NewDeps(cfg *config.Config, handlers *restHTTP.Handlers, logger *slog.Logger) (Deps, error) {
	if cfg == nil || handlers == nil || logger == nil {
		return Deps{}, errors.New("invalid dependencies")
	}
	return Deps{
		cfg:      cfg,
		handlers: handlers,
		log:      logger,
	}, nil
}

type Server struct {
	httpServer *http.Server
	log        *slog.Logger
}

func NewServer(dep Deps) (*Server, error) {

	r := chi.NewRouter()
	deps, err := route.NewDeps(
		dep.cfg,
		dep.log.With("component", "restHTTP.routes"),
		r,
		dep.handlers,
	)
	if err != nil {
		dep.log.Error("failed to create routes dependencies", slog.Any("error", err))
		return nil, fmt.Errorf("failed to create routes dependencies: %w", err)
	}
	route.NewRouter(deps)
	srv := &http.Server{
		Handler:           r,
		Addr:              dep.cfg.HTTPServer.Address,
		ReadTimeout:       dep.cfg.HTTPServer.ReadTimeout,
		ReadHeaderTimeout: dep.cfg.HTTPServer.ReadTimeout,
		WriteTimeout:      dep.cfg.HTTPServer.WriteTimeout,
		IdleTimeout:       dep.cfg.HTTPServer.IdleTimeout,
		MaxHeaderBytes:    dep.cfg.HTTPServer.MaxHeaderBytes,
	}

	return &Server{
		httpServer: srv,
		log:        dep.log,
	}, nil
}

func (s *Server) Run() error {
	s.log.Info("http server starting",
		slog.String("address", s.httpServer.Addr),
		slog.String("read_timeout", s.httpServer.ReadTimeout.String()),
		slog.String("write_timeout", s.httpServer.WriteTimeout.String()),
		slog.String("idle_timeout", s.httpServer.IdleTimeout.String()),
	)

	if err := http.ListenAndServe(s.httpServer.Addr, s.httpServer.Handler); err != nil {
		s.log.Error("http server error", slog.Any("error", err))
		return err
	}

	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.log.Info("http server shutdown initiated")
	return s.httpServer.Shutdown(ctx)
}
