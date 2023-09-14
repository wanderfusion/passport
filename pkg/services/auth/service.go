package auth

import (
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"

	"github.com/akxcix/passport/pkg/commons/stringutils"
	"github.com/akxcix/passport/pkg/config"
	"github.com/akxcix/passport/pkg/jwt"
	"github.com/akxcix/passport/pkg/repositories/auth"
)

type Service struct {
	JwtManager *jwt.JwtManager
	AuthRepo   *auth.Database
}

func New(dbConf *config.DatabaseConfig, jwtConf *config.Jwt) *Service {
	validateConfig(dbConf, jwtConf)

	authRepo := auth.New(dbConf)
	jwtManager := jwt.New(
		jwtConf.Secret,
		jwtConf.TokenValidMicroSeconds,
		jwtConf.RefreshValidMicroSeconds,
	)

	return &Service{
		JwtManager: jwtManager,
		AuthRepo:   authRepo,
	}
}

func validateConfig(dbConf *config.DatabaseConfig, jwtConf *config.Jwt) {
	if dbConf == nil || jwtConf == nil {
		log.Fatal().Msg("Invalid config")
	}
}

func (s *Service) GetUsersUsingUUIDs(userIDs []uuid.UUID) ([]auth.User, error) {
	users, err := s.AuthRepo.FetchUsersUsingUUIDs(userIDs)
	return users, err
}

func (s *Service) GetUsersUsingUsernames(usernames []string) ([]auth.User, error) {
	users, err := s.AuthRepo.FetchUsersUsingUsernames(usernames)
	return users, err
}

func (s *Service) RegisterUser(email, password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 8)
	if err != nil {
		return "", err
	}

	username := stringutils.GenerateUsername(email)
	defaultProfilePic := "https://avatar.vercel.sh/monsters.png"

	err = s.AuthRepo.RegisterUser(email, string(hashedPassword), username, defaultProfilePic)
	if err != nil {
		return "", err
	}

	msg := "Sign up successful"
	return msg, nil
}

func (s *Service) LoginUser(email, password string) (jwt.TokenPair, error) {
	user, err := s.AuthRepo.FetchUserDataByEmail(email)
	if err != nil {
		return jwt.TokenPair{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password))
	if err != nil {
		return jwt.TokenPair{}, err
	}

	authToken, err := s.JwtManager.GenerateToken(user.ID, user.Email, user.Username, user.ProfilePic, jwt.AuthToken)
	if err != nil {
		return jwt.TokenPair{}, err
	}
	refreshToken, err := s.JwtManager.GenerateToken(user.ID, user.Email, user.Username, user.ProfilePic, jwt.RefreshToken)
	if err != nil {
		return jwt.TokenPair{}, err
	}

	tokenPair := jwt.TokenPair{
		AuthToken:    authToken,
		RefreshToken: refreshToken,
	}

	return tokenPair, err
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

func (s *Service) RenewAuthToken(token string) (string, bool) {
	claims, err := s.JwtManager.Verify(token, jwt.RefreshToken)
	if err != nil {
		return "", false
	}

	authToken, err := s.JwtManager.GenerateToken(claims.ID, claims.Email, &claims.Username, &claims.ProfilePic, jwt.AuthToken)
	if err != nil {
		return "", false
	}

	return authToken, true
}

func (s *Service) RenewRefreshToken(token string) (string, bool) {
	claims, err := s.JwtManager.Verify(token, jwt.RefreshToken)
	if err != nil {
		return "", false
	}

	userId := claims.ID
	user, err := s.AuthRepo.FetchUserDataByID(userId)
	if err != nil {
		return "", false
	}

	refreshToken, err := s.JwtManager.GenerateToken(user.ID, user.Email, user.Username, user.ProfilePic, jwt.RefreshToken)
	if err != nil {
		return "", false
	}

	return refreshToken, true
}

func (s *Service) ValidateAuthToken(token string) (*jwt.Claims, bool) {
	claims, err := s.JwtManager.Verify(token, jwt.AuthToken)
	if err != nil {
		return nil, false
	}

	return claims, true
}
