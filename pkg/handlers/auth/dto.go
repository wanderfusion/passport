package auth

type UserAuthReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdateUserReq struct {
	Username   string `json:"username"`
	ProfilePic string `json:"profilePic"`
}

type JwtVerifyRequest struct {
	Jwt string `json:"jwt"`
}
