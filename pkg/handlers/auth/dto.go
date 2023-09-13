package auth

// req --------------------------------------------------------------------------------------------
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

// res --------------------------------------------------------------------------------------------
type UserAuthRes struct {
	TokenPair TokenPairDTO `json:"tokenPair"`
}

type JwtVerifyResponse struct {
	Jwt string `json:"jwt"`
}

// dto --------------------------------------------------------------------------------------------
type TokenPairDTO struct {
	RefreshToken string `json:"refreshToken"`
	AuthToken    string `json:"authToken"`
}
