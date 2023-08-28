package auth

type UserAuthReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type JwtVerifyRequest struct {
	Jwt string `json:"jwt"`
}
