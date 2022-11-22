package redmine

import (
	"encoding/base64"
)

type AuthType int

type Authorization struct {
	User     string
	Password string
}

func (r *Authorization) BuildAuthorizationHeader() string {
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(r.User+":"+r.Password))
}
