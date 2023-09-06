package auth

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/akxcix/passport/pkg/handlers"
	"github.com/akxcix/passport/pkg/services/auth"
	"github.com/google/uuid"
)

type Handlers struct {
	Service *auth.Service
}

func New(s *auth.Service) *Handlers {
	h := Handlers{
		Service: s,
	}

	return &h
}

func (h *Handlers) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var req UserAuthReq
	if err := handlers.FromRequest(r, &req); err != nil {
		handlers.RespondWithError(w, r, err, http.StatusBadRequest)
		return
	}

	if req.Email == "" {
		handlers.RespondWithError(w, r, errors.New("email is empty"), http.StatusBadRequest)
		return
	}

	if req.Password == "" {
		handlers.RespondWithError(w, r, errors.New("password is empty"), http.StatusBadRequest)
		return
	}

	msg, err := h.Service.RegisterUser(req.Email, req.Password)
	if err != nil {
		handlers.RespondWithError(w, r, err, http.StatusInternalServerError)
		return
	}

	handlers.RespondWithData(w, r, msg)
}

func (h *Handlers) LoginUser(w http.ResponseWriter, r *http.Request) {
	var req UserAuthReq
	if err := handlers.FromRequest(r, &req); err != nil {
		handlers.RespondWithError(w, r, err, http.StatusBadRequest)
		return
	}

	if req.Email == "" {
		handlers.RespondWithError(w, r, errors.New("email is empty"), http.StatusBadRequest)
		return
	}

	if req.Password == "" {
		handlers.RespondWithError(w, r, errors.New("password is empty"), http.StatusBadRequest)
		return
	}

	tokenPair, err := h.Service.LoginUser(req.Email, req.Password)
	if err != nil {
		handlers.RespondWithError(w, r, err, http.StatusInternalServerError)
		return
	}

	tokenPairDTO := TokenPairDTO{
		RefreshToken: tokenPair.RefreshToken,
		AuthToken:    tokenPair.AuthToken,
	}

	handlers.RespondWithData(w, r, tokenPairDTO)
}

func (h *Handlers) GenerateAuthToken(w http.ResponseWriter, r *http.Request) {
	var req JwtVerifyRequest
	if err := handlers.FromRequest(r, &req); err != nil {
		handlers.RespondWithError(w, r, err, http.StatusBadRequest)
		return
	}

	authToken, isValid := h.Service.ValidateRefreshToken(req.Jwt)
	if !isValid {
		handlers.RespondWithError(w, r, ErrInvalidJwt, http.StatusUnauthorized)
		return
	}

	handlers.RespondWithData(w, r, authToken)
}

func (h *Handlers) UpdateUser(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(handlers.UserIdContextKey).(uuid.UUID)
	if !ok {
		handlers.RespondWithError(w, r, ErrInvalidJwt, http.StatusBadRequest)
		return
	}
	var req UpdateUserReq
	if err := handlers.FromRequest(r, &req); err != nil {
		handlers.RespondWithError(w, r, err, http.StatusBadRequest)
		return
	}

	err := h.Service.UpdateUser(userID, req.Username, req.ProfilePic)
	if err != nil {
		handlers.RespondWithError(w, r, err, http.StatusInternalServerError)
		return
	}

	handlers.RespondWithData(w, r, "update successful")
}

func (h *Handlers) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bearerToken := r.Header.Get("Authorization")

		// Check if token exists
		if bearerToken == "" {
			handlers.RespondWithError(w, r, ErrInvalidJwt, http.StatusUnauthorized)
			return
		}

		// Extract token from header
		tokenParts := strings.Split(bearerToken, " ")
		if len(tokenParts) != 2 {
			handlers.RespondWithError(w, r, ErrInvalidJwt, http.StatusUnauthorized)
			return
		}
		tokenString := tokenParts[1]

		// Validate token
		claims, isValid := h.Service.ValidateAuthToken(tokenString)
		if claims == nil || !isValid {
			handlers.RespondWithError(w, r, ErrInvalidJwt, http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), handlers.UserIdContextKey, claims.ID)

		// If token is valid, forward to the actual handler
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
