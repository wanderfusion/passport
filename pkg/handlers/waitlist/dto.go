package waitlist

type CreateWaitlistReq struct {
	Mail string `json:"mail"`
	Name string `json:"name"`
}
