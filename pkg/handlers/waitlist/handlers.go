package waitlist

import (
	"net/http"

	"github.com/akxcix/passport/pkg/handlers"
	"github.com/akxcix/passport/pkg/services/waitlist"
)

type Handlers struct {
	Service *waitlist.Service
}

func New(s *waitlist.Service) *Handlers {
	h := Handlers{
		Service: s,
	}

	return &h
}

func (h *Handlers) PostWaitlist(w http.ResponseWriter, r *http.Request) {
	var req CreateWaitlistReq
	unmarshaller := handlers.Unmarshalable[CreateWaitlistReq]{}
	if err := unmarshaller.FromRequest(r, &req); err != nil {
		handlers.RespondWithError(w, r, err, http.StatusBadRequest)
		return
	}

	err := h.Service.AddToWaitlist(req.Mail, req.Name)
	if err != nil {
		handlers.RespondWithError(w, r, err, http.StatusInternalServerError)
		return
	}

	handlers.RespondWithData(w, r, "successfully added to waitlist")
}
