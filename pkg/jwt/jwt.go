package jwt

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/golang-jwt/jwt"
)

var (
	errorJwtExpired           error = errors.New("the jwt is expired")
	errorJwtInvalid           error = errors.New("the jwt provided is invalid")
	errorJwtTokenTypeMismatch error = errors.New("the jwt provided of incorrect type")
)

type TokenType string

type TokenPair struct {
	RefreshToken string
	AuthToken    string
}

const (
	RefreshToken TokenType = "refresh"
	AuthToken    TokenType = "auth"
)

type JwtManager struct {
	secret      []byte
	method      jwt.SigningMethod
	tokenLife   time.Duration
	refreshLife time.Duration
}

type Claims struct {
	ExpiresAt  time.Time `json:"expiresAt"`
	ID         uuid.UUID `json:"id"`
	Email      string    `json:"email"`
	Username   string    `json:"username"`
	ProfilePic string    `json:"profilePicture"`
	TokenType  TokenType `json:"tokenType"`
}

func (c Claims) Valid() error {
	if c.ExpiresAt.After(time.Now()) {
		return nil
	}

	return errorJwtExpired
}

func New(secret string, tokenLifeMicroSeconds, refreshLifeMicroSeconds int64) *JwtManager {
	signingMethod := jwt.SigningMethodHS256
	tokenLife := time.Duration(tokenLifeMicroSeconds) * time.Microsecond
	refreshLife := time.Duration(refreshLifeMicroSeconds) * time.Microsecond

	jwtManager := JwtManager{
		secret:      []byte(secret),
		method:      signingMethod,
		tokenLife:   tokenLife,
		refreshLife: refreshLife,
	}

	return &jwtManager
}

func (j *JwtManager) GenerateTokenPair(uuid uuid.UUID, email, username, profilePic string) (TokenPair, error) {
	authToken, err := j.generateToken(uuid, email, username, profilePic, AuthToken)
	if err != nil {
		log.Error().Err(err).Msg("error signing token")
		return TokenPair{}, err
	}

	refreshToken, err := j.generateToken(uuid, email, username, profilePic, RefreshToken)
	if err != nil {
		log.Error().Err(err).Msg("error signing token")
		return TokenPair{}, err
	}

	pair := TokenPair{
		AuthToken:    authToken,
		RefreshToken: refreshToken,
	}

	return pair, nil
}

func (j *JwtManager) GenerateUsingRefreshToken(refreshToken string) (string, bool) {
	claims, err := j.Verify(refreshToken, RefreshToken)
	if err != nil {
		return "", false
	}

	authToken, err := j.generateToken(claims.ID, claims.Email, claims.Username, claims.ProfilePic, AuthToken)
	if err != nil {
		log.Error().Err(err).Msg("error signing token")
		return "", false
	}

	return authToken, true
}

func (j *JwtManager) generateToken(uuid uuid.UUID, email, username, profilePic string, tokenType TokenType) (string, error) {
	life := j.tokenLife
	if tokenType == RefreshToken {
		life = j.refreshLife
	}
	expirationTime := time.Now().Add(life)
	claims := &Claims{
		ExpiresAt:  expirationTime,
		ID:         uuid,
		Email:      email,
		Username:   username,
		ProfilePic: profilePic,
		TokenType:  tokenType,
	}

	token := jwt.NewWithClaims(j.method, claims)
	tokenString, err := token.SignedString(j.secret)
	if err != nil {
		log.Error().Err(err).Msg("error signing token")
		return "", err
	}

	return tokenString, nil
}

func (j *JwtManager) Verify(token string, tokenType TokenType) (*Claims, error) {
	var keyfunc jwt.Keyfunc = func(token *jwt.Token) (interface{}, error) {
		return j.secret, nil
	}

	claims := &Claims{}

	parsedToken, err := jwt.ParseWithClaims(token, claims, keyfunc)
	if err != nil {
		log.Error().Err(err).Msg("error parsing jwt token")
		return nil, err
	}

	if !parsedToken.Valid {
		log.Error().Msg("token is invalid")
		return nil, errorJwtInvalid
	}

	if claims.TokenType != tokenType {
		log.Error().Msg("token type mismatch")
		return nil, errorJwtTokenTypeMismatch
	}

	return claims, nil
}
