package auth

import (
	"net/http"

	"github.com/akxcix/passport/pkg/handlers"
	"github.com/akxcix/passport/pkg/services/auth"
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

	msg, err := h.Service.RegisterUser(req.Username, req.Password)
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

	msg, err := h.Service.LoginUser(req.Username, req.Password)
	if err != nil {
		handlers.RespondWithError(w, r, err, http.StatusInternalServerError)
		return
	}

	handlers.RespondWithData(w, r, msg)
}

func (h *Handlers) ValidateJwt(w http.ResponseWriter, r *http.Request) {
	var req JwtVerifyRequest
	if err := handlers.FromRequest(r, &req); err != nil {
		handlers.RespondWithError(w, r, err, http.StatusBadRequest)
		return
	}

	isValid := h.Service.ValidateJwt(req.Jwt)
	if !isValid {
		handlers.RespondWithError(w, r, ErrInvalidJwt, http.StatusUnauthorized)
		return
	}

	handlers.RespondWithData(w, r, true)
}
