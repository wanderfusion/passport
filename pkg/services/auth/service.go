package auth

import (
	"github.com/google/uuid"
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

func (s *Service) LoginUser(email, password string) (string, error) {
	user, err := s.AuthRepo.FetchUserDataByEmail(email)
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password))
	if err != nil {
		return "", err
	}

	username := ""
	profilePic := ""

	if user.Username != nil {
		username = *user.Username
	}

	if user.ProfilePic != nil {
		profilePic = *user.ProfilePic
	}

	jwtString, err := s.JwtManager.GenerateJWT(user.ID, email, username, profilePic)
	return jwtString, err
}

func (s *Service) UpdateUser(id uuid.UUID, username, profilePic string) error {
	user := auth.User{}
	user.ID = id

	user.Username = &username
	if username == "" {
		user.Username = nil
	}

	user.ProfilePic = &profilePic
	if profilePic == "" {
		user.ProfilePic = nil
	}

	err := s.AuthRepo.UpdateUserProfile(user)
	return err
}

func (s *Service) ValidateJwt(token string) (*jwt.Claims, bool) {
	claims, err := s.JwtManager.Verify(token)
	if err != nil {
		return nil, false
	}

	return claims, true
}
