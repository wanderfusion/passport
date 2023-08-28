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
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))
	r.Use(handlers.LogRequest)

	// general routes
	r.Get("/health", handlers.HealthCheck)

	authHandlers := authHandlers.New(authService)
	r.Post("/users/register", authHandlers.RegisterUser)
	r.Post("/users/login", authHandlers.LoginUser)
	r.Post("/users/verify", authHandlers.ValidateJwt)
	r.Post("/users/update", authHandlers.AuthMiddleware(authHandlers.UpdateUser))

	return r
}
