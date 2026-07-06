package v1

type CreateRequest struct {
	Nickname string `json:"nickname"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateResponse struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}
