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

func (h *Handlers) PostUser(w http.ResponseWriter, r *http.Request) {
	var req CreateUserReq
	unmarshaller := handlers.Unmarshalable[CreateUserReq]{}
	if err := unmarshaller.FromRequest(r, &req); err != nil {
		handlers.RespondWithError(w, r, err, http.StatusBadRequest)
		return
	}

	err := h.Service.RegisterUser(req.Username, req.Email)
	if err != nil {
		handlers.RespondWithError(w, r, err, http.StatusInternalServerError)
		return
	}

	handlers.RespondWithData(w, r, "successfully added to waitlist")
}

func (h *Handlers) PostGithubRegisterUser(w http.ResponseWriter, r *http.Request) {
	var req CreateGithubUserReq
	unmarshaller := handlers.Unmarshalable[CreateGithubUserReq]{}
	if err := unmarshaller.FromRequest(r, &req); err != nil {
		handlers.RespondWithError(w, r, err, http.StatusBadRequest)
		return
	}

	msg, err := h.Service.RegisterGithubUser(req.Code)
	if err != nil {
		handlers.RespondWithError(w, r, err, http.StatusInternalServerError)
		return
	}

	handlers.RespondWithData(w, r, msg)
}
