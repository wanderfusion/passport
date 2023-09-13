package jwt

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/golang-jwt/jwt"
)

// JwtManager manages the creation and verification of JWTs.
type JwtManager struct {
	secret      []byte
	method      jwt.SigningMethod
	tokenLife   time.Duration
	refreshLife time.Duration
}

// Claims represents the payload of JWT.
type Claims struct {
	ExpiresAt  time.Time `json:"expiresAt"`
	ID         uuid.UUID `json:"id"`
	Email      string    `json:"email"`
	Username   string    `json:"username"`
	ProfilePic string    `json:"profilePicture"`
	TokenType  TokenType `json:"tokenType"`
}

// TokenType is used to distinguish between different types of tokens.
type TokenType string

// TokenPair holds both the refresh and auth tokens.
type TokenPair struct {
	RefreshToken string
	AuthToken    string
}

// Constants for types of tokens.
const (
	RefreshToken TokenType = "refresh"
	AuthToken    TokenType = "auth"
)

// Declare custom errors for JWT verification.
var (
	errorJwtExpired           error = errors.New("the jwt is expired")
	errorJwtInvalid           error = errors.New("the jwt provided is invalid")
	errorJwtTokenTypeMismatch error = errors.New("the jwt provided of incorrect type")
)

// Valid checks if the token has expired.
func (c Claims) Valid() error {
	if c.ExpiresAt.After(time.Now()) {
		return nil
	}

	return errorJwtExpired
}

// New creates a new JwtManager.
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

// GenerateToken creates a new JWT with the given claims.
func (j *JwtManager) GenerateToken(uuid uuid.UUID, email string, username, profilePic *string, tokenType TokenType) (string, error) {
	if username == nil {
		*username = ""
	}

	if profilePic == nil {
		*profilePic = ""
	}

	// Choose token life based on the token type.
	life := j.tokenLife
	if tokenType == RefreshToken {
		life = j.refreshLife
	}
	expirationTime := time.Now().Add(life)
	claims := &Claims{
		ExpiresAt:  expirationTime,
		ID:         uuid,
		Email:      email,
		Username:   *username,
		ProfilePic: *profilePic,
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

// Verify validates a given JWT.
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
