package jwt

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/golang-jwt/jwt"
)

var (
	errorJwtExpired error = errors.New("the jwt is expired")
	errorJwtInvalid error = errors.New("the jwt provided is invalid")
)

type JwtManager struct {
	secret   []byte
	method   jwt.SigningMethod
	validity time.Duration
}

type Claims struct {
	ExpiresAt  time.Time `json:"expiresAt"`
	ID         uuid.UUID `json:"id"`
	Email      string    `json:"email"`
	Username   string    `json:"username"`
	ProfilePic string    `json:"profilePicture"`
}

func (c Claims) Valid() error {
	if c.ExpiresAt.After(time.Now()) {
		return nil
	}

	return errorJwtExpired
}

func New(secret string, validMins int) *JwtManager {
	signingMethod := jwt.SigningMethodHS256
	validity := time.Duration(validMins) * time.Minute

	jwtManager := JwtManager{
		secret:   []byte(secret),
		method:   signingMethod,
		validity: validity,
	}

	return &jwtManager
}

func (j *JwtManager) GenerateJWT(uuid uuid.UUID, email, username, profilePic string) (string, error) {
	expirationTime := time.Now().Add(j.validity)
	claims := &Claims{
		ExpiresAt:  expirationTime,
		ID:         uuid,
		Email:      email,
		Username:   username,
		ProfilePic: profilePic,
	}

	token := jwt.NewWithClaims(j.method, claims)
	tokenString, err := token.SignedString(j.secret)
	if err != nil {
		log.Error().Err(err).Msg("error signing token")
		return "", err
	}

	return tokenString, nil
}

func (j *JwtManager) Verify(token string) (*Claims, error) {
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

	return claims, nil
}
