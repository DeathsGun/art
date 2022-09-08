package redmine

import (
	"encoding/base64"
)

type RedmineAuthType int

type RedmineAuthorization struct {
	RedmineUser     string
	RedminePassword string
}

func (r *RedmineAuthorization) BuildAuthorizationHeader() string {
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(r.RedmineUser+":"+r.RedminePassword))
}

func AuthorizeHTTP(username string, password string) *RedmineAuthorization {
	return &RedmineAuthorization{
		RedmineUser:     username,
		RedminePassword: password,
	}
}
