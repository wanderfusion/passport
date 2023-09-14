package app

import (
	"net/http"

	"github.com/akxcix/passport/pkg/handlers"
	authHandlers "github.com/akxcix/passport/pkg/handlers/auth"
	"github.com/akxcix/passport/pkg/services/auth"
	"github.com/rs/zerolog/log"
	limiter "github.com/ulule/limiter/v3"
	mhttp "github.com/ulule/limiter/v3/drivers/middleware/stdlib"
	"github.com/ulule/limiter/v3/drivers/store/memory"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

// -----------------------------------------------------------------------------------
func createRoutes(authService *auth.Service) *chi.Mux {
	authH := authHandlers.New(authService)

	r := chi.NewRouter()
	applyGlobalMiddlewares(r)
	r.Get("/health", handlers.HealthCheck)
	applyAuthRoutes(r, authH)

	return r
}

func applyGlobalMiddlewares(r *chi.Mux) {
	r.Use(corsHandler())
	r.Use(rateLimiteMiddleware())
	r.Use(handlers.LogRequest)
}

func corsHandler() func(http.Handler) http.Handler {
	return cors.Handler(
		cors.Options{
			AllowedOrigins:   []string{"https://*", "http://*"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: false,
			MaxAge:           300,
		},
	)
}

func rateLimiteMiddleware() func(http.Handler) http.Handler {
	rate, err := limiter.NewRateFromFormatted("2000-M")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize limiter")
	}

	store := memory.NewStore()

	middleware := mhttp.NewMiddleware(limiter.New(
		store,
		rate,
		limiter.WithTrustForwardHeader(true),
	))

	return middleware.Handler
}

func applyAuthRoutes(r *chi.Mux, authH *authHandlers.Handlers) {
	r.Route("/users", func(r chi.Router) {
		r.Get("/{userIDs}", authH.GetUsersUsingUUIDs)
		r.Post("/register", authH.RegisterUser)
		r.Post("/login", authH.LoginUser)
		r.With(authH.AuthMiddleware).Put("/me", authH.UpdateUser)
	})

	r.Route("/auth", func(r chi.Router) {
		r.Post("/tokens/renew-auth", authH.RenewAuthToken)
		r.Post("/tokens/renew-refresh", authH.RenewRefreshToken)
	})
}
