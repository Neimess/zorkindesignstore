// internal/adapter/rest/server.go
package rest

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/Neimess/zorkin-store-project/internal/config"
	"github.com/Neimess/zorkin-store-project/internal/transport/http/restHTTP"
	"github.com/Neimess/zorkin-store-project/internal/transport/http/routes"
	"github.com/go-chi/chi/v5"
)

type Server struct {
	httpServer *http.Server
	logger     *slog.Logger
}

func NewServer(cfg config.HTTPServer, handlers *restHTTP.Handlers, logger *slog.Logger) *Server {
	r := chi.NewRouter()
	route.NewRouter(route.Deps{
		Router:   r,
		Handlers: handlers,
		Logger:   logger,
	})

	srv := &http.Server{
		Handler:           r,
		Addr:              cfg.Address,
		ReadTimeout:       cfg.ReadTimeout,
		ReadHeaderTimeout: cfg.ReadTimeout,
		WriteTimeout:      cfg.WriteTimeout,
		IdleTimeout:       cfg.IdleTimeout,
		MaxHeaderBytes:    cfg.MaxHeaderBytes,
	}

	return &Server{
		httpServer: srv,
		logger:     logger.With(slog.String("component", "rest_server")),
	}
}

func (s *Server) Run() error {
	s.logger.Info("http server starting",
		slog.String("address", s.httpServer.Addr),
		slog.String("read_timeout", s.httpServer.ReadTimeout.String()),
		slog.String("write_timeout", s.httpServer.WriteTimeout.String()),
		slog.String("idle_timeout", s.httpServer.IdleTimeout.String()),
	)

	if err := http.ListenAndServe(s.httpServer.Addr, s.httpServer.Handler); err != nil {
		s.logger.Error("http server error", slog.Any("error", err))
		return err
	}

	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("http server shutdown initiated")
	return s.httpServer.Shutdown(ctx)
}
