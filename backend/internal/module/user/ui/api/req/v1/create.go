package v1

type CreateRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateResponse struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}
