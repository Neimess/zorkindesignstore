package rest

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/Neimess/zorkin-store-project/internal/config"
	"github.com/Neimess/zorkin-store-project/internal/transport/http/restHTTP"
	route "github.com/Neimess/zorkin-store-project/internal/transport/http/routes"
	"github.com/go-chi/chi/v5"
)

type Deps struct {
	Server   config.HTTPServer
	Config   *config.Config
	Logger   *slog.Logger
	Handlers *restHTTP.Handlers
}

type Server struct {
	httpServer *http.Server
	logger     *slog.Logger
}

func NewServer(dep Deps) *Server {
	r := chi.NewRouter()
	route.NewRouter(route.Deps{
		Router:   r,
		Handlers: dep.Handlers,
		Logger:   dep.Logger,
		Config:   dep.Config,
	})

	srv := &http.Server{
		Handler:           r,
		Addr:              dep.Server.Address,
		ReadTimeout:       dep.Server.ReadTimeout,
		ReadHeaderTimeout: dep.Server.ReadTimeout,
		WriteTimeout:      dep.Server.WriteTimeout,
		IdleTimeout:       dep.Server.IdleTimeout,
		MaxHeaderBytes:    dep.Server.MaxHeaderBytes,
	}

	return &Server{
		httpServer: srv,
		logger:     dep.Logger.With(slog.String("component", "rest_server")),
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
