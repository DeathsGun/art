package redmine

import "time"

type AccountResponse struct {
	User User `json:"user"`
}

type User struct {
	Id        int       `json:"id"`
	UserName  string    `json:"login"`
	IsAdmin   bool      `json:"admin"`
	FirstName string    `json:"firstname"`
	LastName  string    `json:"lastname"`
	Mail      string    `json:"mail"`
	CreatedOn time.Time `json:"created_on"`
	LastLogin time.Time `json:"last_login_on"`
	ApiKey    string    `json:"api_key"`
}
