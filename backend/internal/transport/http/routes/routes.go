package route

import (
	"fmt"
	"log/slog"
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
	r.Use(httplog.RequestLogger(deps.Logger, &httplog.Options{RecoverPanics: true, LogRequestHeaders: []string{"Origin"}}))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// ── profiler & swagger ───────────────────────────────────────────────
	registerSwaggerRoutes(r)
	r.Mount("/debug/pprof", profiler(deps.Config.Env))

	// ── public API ───────────────────────────────────────────────────────
	r.Route("/api", func(r chi.Router) {
		registerBaseRoutes(r)
		registerProductPublicRoutes(r, deps.Handlers.ProductHandler)
		registerCategoryPublicRoutes(r, deps.Handlers.CategoryHandler)
		registerPresetPublicRoutes(r, deps.Handlers.PresetHandler)

		// ── admin zone ──────────────────────────────────────────────────
		r.Route("/admin", func(r chi.Router) {
			// one‑time login link
			r.Get(fmt.Sprintf("/auth/%s", deps.Config.AdminCode), deps.Handlers.AuthHandler.Login)

			// JWT‑protected block
			cfg := deps.Config.JWTConfig
			jwtMW, _ := customMiddlewares.NewJWTMiddleware(customMiddlewares.JWTCfg{
				Secret: []byte(cfg.SecretKey), Algorithm: cfg.Algorithm, Issuer: cfg.Issuer, Audience: cfg.Audience,
			})
			r.Group(func(r chi.Router) {
				r.Use(jwtMW.CheckJWT)

				registerProductAdminRoutes(r, deps.Handlers.ProductHandler)
				registerCategoryAdminRoutes(r, deps.Handlers.CategoryHandler, deps.Handlers.CategoryAttributeHandler)
				registerPresetAdminRoutes(r, deps.Handlers.PresetHandler)
			})
		})
	})

	return r
}
