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
	Router   chi.Router
	Logger   *slog.Logger
	Handlers *restHTTP.Handlers
	Config   *config.Config
}

func NewRouter(deps Deps) chi.Router {
	r := deps.Router

	// ── global middleware ────────────────────────────────────────────────
	r.Use(middleware.RequestID, middleware.RealIP, middleware.Recoverer)
	r.Use(middleware.Timeout(30*time.Second), middleware.Compress(5))
	r.Use(httplog.RequestLogger(slog.Default(), &httplog.Options{
		RecoverPanics:      true,
		LogRequestHeaders:  []string{"Origin", "User-Agent", "Accept", "Content-Type", "X-Request-ID"},
		LogResponseHeaders: []string{"Content-Type", "Content-Length", "X-Request-ID"},
		LogRequestBody: func(req *http.Request) bool {
			return req.Method == http.MethodPost && strings.HasPrefix(req.URL.Path, "/debug/")
		},
	}))
	if deps.Config.HTTPServer.EnableCORS {
		r.Use(cors.Handler(cors.Options{
			AllowedOrigins:   []string{"*"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			AllowCredentials: true,
			MaxAge:           300,
		}))
	}

	// ── profiler & swagger ───────────────────────────────────────────────
	isDev := deps.Config.Env == config.EnvLocal || deps.Config.Env == config.EnvDev
	if isDev && deps.Config.HTTPServer.EnablePProf {
		r.Mount("/debug/pprof", profiler(deps.Config.Env))
	}

	// ── public API ───────────────────────────────────────────────────────
	r.Route("/api", func(r chi.Router) {
		if isDev && deps.Config.Swagger.Enabled {
			registerSwaggerRoutes(r)
		}
		registerBaseRoutes(r)
		registerProductPublicRoutes(r, deps.Handlers.ProductHandler)
		registerCategoryWithAttrsPublicRoutes(r, deps.Handlers.CategoryHandler, deps.Handlers.AttributeHandler)
		registerPresetPublicRoutes(r, deps.Handlers.PresetHandler)
		// ── admin zone ──────────────────────────────────────────────────
		r.Route("/admin", func(r chi.Router) {
			// one‑time login link
			r.Get(fmt.Sprintf("/auth/%s", deps.Config.AdminCode), deps.Handlers.AuthHandler.Login)

			// JWT‑protected block
			cfg := deps.Config.JWTConfig
			jwtMW, _ := customMiddlewares.NewJWTMiddleware(customMiddlewares.JWTCfg{
				Secret: []byte(cfg.JWTSecret), Algorithm: cfg.Algorithm, Issuer: cfg.Issuer, Audience: cfg.Audience,
			})
			r.Group(func(r chi.Router) {
				r.Use(jwtMW.CheckJWT)

				registerProductAdminRoutes(r, deps.Handlers.ProductHandler)
				registerCategoryWithAttrsAdminRoutes(r, deps.Handlers.CategoryHandler, deps.Handlers.AttributeHandler)
				registerPresetAdminRoutes(r, deps.Handlers.PresetHandler)
			})
		})
	})

	return r
}
