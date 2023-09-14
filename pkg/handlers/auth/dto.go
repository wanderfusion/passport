package auth

import "github.com/google/uuid"

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

type GetUsersUsingUUIDsResponse struct {
	Users []UserDTO `json:"users"`
}

// dto --------------------------------------------------------------------------------------------
type TokenPairDTO struct {
	RefreshToken string `json:"refreshToken"`
	AuthToken    string `json:"authToken"`
}

type UserDTO struct {
	ID         uuid.UUID `json:"id"`
	Username   *string   `json:"username"`
	ProfilePic *string   `json:"profilePic"`
}
