package app

import (
	"github.com/akxcix/passport/pkg/handlers"
	authHandlers "github.com/akxcix/passport/pkg/handlers/auth"
	"github.com/akxcix/passport/pkg/services/auth"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func createRoutes(authService *auth.Service) *chi.Mux {
	r := chi.NewRouter()

	// global middlewares
	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	r.Use(handlers.LogRequest)

	// general routes
	r.Get("/health", handlers.HealthCheck)

	authHandlers := authHandlers.New(authService)
	r.Post("/users", authHandlers.PostUser)
	r.Post("/users/register/github", authHandlers.PostGithubRegisterUser)

	return r
}
