package auth

import (
	"github.com/rs/zerolog/log"

	"github.com/akxcix/passport/pkg/clients/oauth/github"
	"github.com/akxcix/passport/pkg/config"
	"github.com/akxcix/passport/pkg/repositories/auth"
)

type Service struct {
	AuthRepo          *auth.Database
	GithubOauthClient *github.GithubOauthClient
}

func New(conf *config.DatabaseConfig) *Service {
	if conf == nil {
		log.Fatal().Msg("Conf is nil")
	}

	authRepo := auth.New(conf)
	githubOauthClient := github.New("https://github.com", "6f9e7c286c45c0aa4169", "774e7cc8ecff7d2a4d2e0c341b3e379aa1bcc468")

	svc := &Service{
		AuthRepo:          authRepo,
		GithubOauthClient: githubOauthClient,
	}

	return svc
}

func (s *Service) RegisterUser(username, mail string) error {
	return s.AuthRepo.RegisterUser(username, mail)
}

func (s *Service) RegisterGithubUser(code string) (string, error) {
	return s.GithubOauthClient.GetToken(code)
}
