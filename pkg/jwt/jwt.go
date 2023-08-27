package jwt

import (
	"errors"
	"time"

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
	Username  string    `json:"username"`
	ExpiresAt time.Time `json:"expiresAt"`
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

func (j *JwtManager) GenerateJWT(username string) (string, error) {
	expirationTime := time.Now().Add(j.validity)
	claims := &Claims{
		Username:  username,
		ExpiresAt: expirationTime,
	}

	token := jwt.NewWithClaims(j.method, claims)
	tokenString, err := token.SignedString(j.secret)
	if err != nil {
		log.Error().Err(err).Msg("error signing token")
		return "", err
	}

	return tokenString, nil
}

func (j *JwtManager) Verify(token string) error {
	var keyfunc jwt.Keyfunc = func(token *jwt.Token) (interface{}, error) {
		return j.secret, nil
	}

	claims := &Claims{}

	parsedToken, err := jwt.ParseWithClaims(token, claims, keyfunc)
	if err != nil {
		log.Error().Err(err).Msg("error parsing jwt token")
		return err
	}

	if !parsedToken.Valid {
		log.Error().Msg("token is invalid")
		return errorJwtInvalid
	}

	return nil
}
