package dto

type LoginRequest struct {
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
	School   string `json:"school" form:"school"`
}
