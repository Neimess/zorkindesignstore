package route

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	_ "github.com/Neimess/zorkin-store-project/docs"
	"github.com/Neimess/zorkin-store-project/internal/config"
	"github.com/Neimess/zorkin-store-project/internal/transport/http/restHTTP"
	customMiddlewares "github.com/Neimess/zorkin-store-project/pkg/httputils/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httplog/v3"
)

type Deps struct {
	config   *config.Config
	logger   *slog.Logger
	router   chi.Router
	handlers *restHTTP.Handlers
}

func NewDeps(cfg *config.Config, logger *slog.Logger, router chi.Router, handlers *restHTTP.Handlers) (Deps, error) {
	if cfg == nil || logger == nil || handlers == nil {
		return Deps{}, fmt.Errorf("invalid dependencies")
	}
	return Deps{
		config:   cfg,
		logger:   logger,
		router:   router,
		handlers: handlers,
	}, nil
}

func NewRouter(deps Deps) chi.Router {
	r := deps.router

	// ── global middleware ────────────────────────────────────────────────
	r.Use(middleware.RequestID, middleware.RealIP, middleware.Recoverer)
	r.Use(middleware.Timeout(30*time.Second), middleware.Compress(5))
	r.Use(httplog.RequestLogger(deps.logger.With("component", "http"), &httplog.Options{
		RecoverPanics:      true,
		LogRequestHeaders:  []string{"Origin", "User-Agent", "Accept", "Content-Type", "X-Request-ID"},
		LogResponseHeaders: []string{"Content-Type", "Content-Length", "X-Request-ID"},
		LogRequestBody: func(req *http.Request) bool {
			return req.Method == http.MethodPost && strings.HasPrefix(req.URL.Path, "/debug/")
		},
	}))
	isDev := deps.config.Env == config.EnvLocal || deps.config.Env == config.EnvDev
	if isDev && deps.config.HTTPServer.EnableCORS {
		r.Use(cors.Handler(cors.Options{
			AllowedOrigins:   []string{"*"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			AllowCredentials: true,
			MaxAge:           300,
		}))
	}

	// ── profiler & swagger ───────────────────────────────────────────────
	if isDev && deps.config.HTTPServer.EnablePProf {
		r.Mount("/debug/pprof", profiler(deps.config.Env))
	}

	// ── public API ───────────────────────────────────────────────────────
	r.Route("/api", func(r chi.Router) {
		if isDev && deps.config.Swagger.Enabled {
			registerSwaggerRoutes(r)
		}
		registerBaseRoutes(r)
		registerProductPublicRoutes(r, deps.handlers.ProductHandler)
		registerCategoryWithAttrsPublicRoutes(r, deps.handlers.CategoryHandler, deps.handlers.AttributeHandler)
		registerPresetPublicRoutes(r, deps.handlers.PresetHandler)
		// ── admin zone ──────────────────────────────────────────────────
		r.Route("/admin", func(r chi.Router) {
			// one‑time login link
			r.Get(fmt.Sprintf("/auth/%s", deps.config.AdminCode), deps.handlers.AuthHandler.Login)

			// JWT‑protected block
			cfg := deps.config.JWTConfig
			jwtMW, _ := customMiddlewares.NewJWTMiddleware(customMiddlewares.JWTCfg{
				Secret: []byte(cfg.JWTSecret), Algorithm: cfg.Algorithm, Issuer: cfg.Issuer, Audience: cfg.Audience,
			})
			r.Group(func(r chi.Router) {
				r.Use(jwtMW.CheckJWT)

				registerProductAdminRoutes(r, deps.handlers.ProductHandler)
				registerCategoryWithAttrsAdminRoutes(r, deps.handlers.CategoryHandler, deps.handlers.AttributeHandler)
				registerPresetAdminRoutes(r, deps.handlers.PresetHandler)
				registerCoefficientsAdminRoutes(r, deps.handlers.CoefficientsHandler)
				registerServiceAdminRoutes(r, deps.handlers.ServiceHandler)
			})
		})
	})

	return r
}
