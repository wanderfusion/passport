package auth

type UserAuthReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type JwtVerifyRequest struct {
	Jwt string `json:"jwt"`
}
