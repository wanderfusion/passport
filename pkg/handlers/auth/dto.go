package auth

type CreateUserReq struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

type CreateGithubUserReq struct {
	Code string `json:"code"`
}
