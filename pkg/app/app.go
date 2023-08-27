package app

import (
	"fmt"
	"net/http"

	"github.com/akxcix/passport/pkg/config"
	"github.com/akxcix/passport/pkg/services/auth"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type application struct {
	Config      *config.Config
	AuthService *auth.Service
	Routes      *chi.Mux
}

func readConfigs() *config.Config {
	config, err := config.Read("./config.yml")
	if err != nil {
		log.Fatal().Err(err)
	}

	return config
}

func createServices(conf *config.Config) *auth.Service {
	if conf == nil {
		log.Fatal().Msg("Conf is nil")
	}

	authService := auth.New(conf.Database)

	return authService
}

func new() *application {
	config := readConfigs()

	authService := createServices(config)
	routes := createRoutes(authService)

	app := application{
		Config:      config,
		AuthService: authService,
		Routes:      routes,
	}

	return &app
}

func Run() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	app := new()

	addr := fmt.Sprintf("%s:%s", app.Config.Server.Host, app.Config.Server.Port)
	log.Info().Msg(fmt.Sprintf("Running application at %s", addr))
	err := http.ListenAndServe(addr, app.Routes)
	log.Fatal().Err(err).Msg("Crashed")
}
