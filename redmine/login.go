package redmine

import (
	"encoding/base64"
	"errors"
)

type RedmineAuthType int

const (
	HTTP   RedmineAuthType = 0
	APIKEY RedmineAuthType = 1
)

type RedmineAuthorization struct {
	Type            RedmineAuthType
	RedmineUser     string
	RedminePassword string
}

func (r *RedmineAuthorization) BuildAuthorizationHeader() string {
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(r.RedmineUser+":"+r.RedminePassword))
}

func (r *RedmineAuthorization) getUsername() (string, error) {
	if r.Type == HTTP {
		return r.RedmineUser, nil
	}
	return "", errors.New("only the HTTP-Authorization has the Username")
}

func (r *RedmineAuthorization) getPassword() (string, error) {
	if r.Type == HTTP {
		return r.RedminePassword, nil
	}
	return "", errors.New("only the HTTP-Authorization has the Password")
}

func (r *RedmineAuthorization) getAPIKey() (string, error) {
	if r.Type == APIKEY {
		return r.RedmineUser, nil
	}
	return "", errors.New("there is no API Key for the HTTP-Authorization")
}

func AuthorizeHTTP(username string, password string) *RedmineAuthorization {
	return &RedmineAuthorization{
		Type:            HTTP,
		RedmineUser:     username,
		RedminePassword: password,
	}
}

func AuthorizeAPIKey(apiKey string) *RedmineAuthorization {
	return &RedmineAuthorization{
		Type:            APIKEY,
		RedmineUser:     apiKey,
		RedminePassword: "",
	}
}
