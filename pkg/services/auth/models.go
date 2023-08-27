package auth

import "time"

type JwtData struct {
	Username  string    `json:"username"`
	ExpiresAt time.Time `json:"expiresAt"`
}
