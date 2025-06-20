package route

import (
	"fmt"
	"net/http"
	"net/http/pprof"
	"time"

	_ "github.com/Neimess/zorkin-store-project/docs"
	customMiddlewares "github.com/Neimess/zorkin-store-project/pkg/http/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httplog/v3"
	httpSwagger "github.com/swaggo/http-swagger"
)

func NewRouter(deps Deps) chi.Router {
	r := deps.Router
	jwtMiddleware, err := customMiddlewares.NewJWTMiddleware(customMiddlewares.JWTCfg{
		Secret:    []byte(deps.Config.JWTConfig.SecretKey),
		Algorithm: deps.Config.JWTConfig.Algorithm,
		Issuer:    deps.Config.JWTConfig.Issuer,
		Audience:  deps.Config.JWTConfig.Audience,
	})
	if err != nil {
		panic("failed to create JWT middleware: " + err.Error())
	}

	// middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(30 * time.Second))
	r.Use(middleware.Compress(5))
	r.Use(httplog.RequestLogger(deps.Logger, &httplog.Options{
		RecoverPanics:      true,
		LogRequestHeaders:  []string{"Origin"},
		LogResponseHeaders: []string{},
	}))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// API group
	r.Route("/api", func(r chi.Router) {

		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("OK"))
		})

		ph := deps.Handlers.ProductHandler
		r.Route("/product", func(r chi.Router) {

			r.Get("/{id}", ph.GetDetailed)
		})
		ch := deps.Handlers.CategoryHandler
		r.Route("/category", func(r chi.Router) {

			r.Get("/{id}", ch.Get)
			r.Get("/", ch.List)
			r.Post("/", ch.Create)
			r.Put("/{id}", ch.Update)
			r.Delete("/{id}", ch.Delete)
		})
		authH := deps.Handlers.AuthHandler
		r.Route("/admin", func(r chi.Router) {
			r.Get(fmt.Sprintf("/auth/%s", deps.Config.AdminCode), authH.Login)


			r.Route("/product", func(r chi.Router) {
				r.Use(jwtMiddleware.CheckJWT)

				r.Post("/", ph.Create)
				r.Post("/detailed", ph.CreateWithAttributes)
			})
		})
	})

	r.Get("/swagger", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/swagger/index.html", http.StatusMovedPermanently)
	})
	r.Get("/swagger/*", httpSwagger.WrapHandler)

	r.Mount("/debug/pprof", profiler())

	return r
}

func profiler() http.Handler {
	r := chi.NewRouter()
	r.Handle("/", http.HandlerFunc(pprof.Index))
	r.Handle("/cmdline", http.HandlerFunc(pprof.Cmdline))
	r.Handle("/profile", http.HandlerFunc(pprof.Profile))
	r.Handle("/symbol", http.HandlerFunc(pprof.Symbol))
	r.Handle("/trace", http.HandlerFunc(pprof.Trace))
	r.Handle("/allocs", pprof.Handler("allocs"))
	r.Handle("/block", pprof.Handler("block"))
	r.Handle("/goroutine", pprof.Handler("goroutine"))
	r.Handle("/heap", pprof.Handler("heap"))
	r.Handle("/mutex", pprof.Handler("mutex"))
	r.Handle("/threadcreate", pprof.Handler("threadcreate"))
	return r
}
