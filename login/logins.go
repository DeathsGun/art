package login

type Login struct {
	Name     string `json:"provider"`
	Username string `json:"username"`
	Password string `json:"password"`
}
