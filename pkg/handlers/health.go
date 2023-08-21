package handlers

import (
	"net/http"
)

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	data := "I fly like paper, get high like planes. If you catch me at the border I got visas in my name"
	RespondWithData(w, r, data)
}
