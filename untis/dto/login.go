package dto

type LoginResponse struct {
	SessionId  string `json:"sessionId"`
	PersonType int    `json:"personType"`
	PersonId   int    `json:"personId"`
}

type LoginRequest struct {
	Username string `json:"user"`
	Password string `json:"password"`
	Client   string `json:"client"`
}
