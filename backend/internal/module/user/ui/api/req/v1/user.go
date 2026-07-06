package v1

type UserVO struct {
	ID       int64  `json:"id"`
	Nickname string `json:"nickname"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type LoginVO struct {
	Token string `json:"token"`
	User  UserVO `json:"user"`
}
