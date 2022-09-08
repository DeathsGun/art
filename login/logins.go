package login

// Login stores the credentials for the provider
type Login struct {
	Name     string `json:"provider"`
	Username string `json:"username"`
	Password string `json:"password"`
}
