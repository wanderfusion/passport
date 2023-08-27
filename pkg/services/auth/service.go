package auth

import (
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"

	"github.com/akxcix/passport/pkg/config"
	"github.com/akxcix/passport/pkg/jwt"
	"github.com/akxcix/passport/pkg/repositories/auth"
)

type Service struct {
	JwtManager *jwt.JwtManager
	AuthRepo   *auth.Database
}

func New(dbConf *config.DatabaseConfig, jwtConf *config.Jwt) *Service {
	if dbConf == nil {
		log.Fatal().Msg("dbConf is nil")
	}

	if jwtConf == nil {
		log.Fatal().Msg("jwtConf is nil")
	}

	authRepo := auth.New(dbConf)
	jwtManager := jwt.New(jwtConf.Secret, jwtConf.ValidMins)

	svc := &Service{
		JwtManager: jwtManager,
		AuthRepo:   authRepo,
	}

	return svc
}

func (s *Service) RegisterUser(username, password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 8)
	if err != nil {
		return "", err
	}

	err = s.AuthRepo.RegisterUser(username, string(hashedPassword))
	if err != nil {
		return "", err
	}

	msg := "Sign up successful"
	return msg, nil
}

func (s *Service) LoginUser(username, password string) (string, error) {
	hashedPassword, err := s.AuthRepo.FetchHashByUsername(username)
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return "", err
	}

	jwtString, err := s.JwtManager.GenerateJWT(username)
	return jwtString, err
}

func (s *Service) ValidateJwt(token string) bool {
	err := s.JwtManager.Verify(token)
	if err != nil {
		return false
	}

	return true
}
