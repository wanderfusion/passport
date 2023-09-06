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

func createRoutes(authService *auth.Service) *chi.Mux {
	rateLimiter := NewRateLimiter()
	limiterMiddleware := rateLimiter.rateLimitMiddleware()
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
	r.Use(limiterMiddleware)
	r.Use(handlers.LogRequest)

	// general routes
	r.Get("/health", handlers.HealthCheck)

	authHandlers := authHandlers.New(authService)
	r.Post("/users/register", authHandlers.RegisterUser)
	r.Post("/users/login", authHandlers.LoginUser)
	r.Post("/users/token/refresh", authHandlers.GenerateAuthToken)
	r.Post("/users/update", authHandlers.AuthMiddleware(authHandlers.UpdateUser))

	return r
}

// rate limiter -----------------------------------------------------------------------------------
type RateLimiter struct {
	store limiter.Store
	rate  limiter.Rate
}

func NewRateLimiter() *RateLimiter {
	rate, err := limiter.NewRateFromFormatted("2000-M")
	if err != nil {
		log.Fatal().Err(err).Msg("unable to initialise limiter")
	}

	store := memory.NewStore()

	limiter := RateLimiter{
		store: store,
		rate:  rate,
	}
	return &limiter
}

func (l *RateLimiter) rateLimitMiddleware() func(h http.Handler) http.Handler {
	middleware := mhttp.NewMiddleware(limiter.New(
		l.store,
		l.rate,
		limiter.WithTrustForwardHeader(true),
	))

	return middleware.Handler
}
