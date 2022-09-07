package redmine

import (
	"github.com/deathsgun/art/provider"
	"time"
)

type redmineImportProvider struct {
}

const redmineURL = "https://joope.de/"

func (r *redmineImportProvider) Name() string {
	return "redmine"
}

func (r *redmineImportProvider) ValidateLogin(username string, password string) (string, string, error) {
	redmineApi, err := NewRedmineAPI(redmineURL, AuthorizeHTTP(username, password))
	if err != nil {
		return username, password, nil
	}
	return redmineApi.RedmineUser.ApiKey, "", nil
}

func (r *redmineImportProvider) NeedsLogin() bool {
	return true
}

func (r *redmineImportProvider) Import(startDate time.Time) ([]*provider.Entry, error) {

	return nil, nil
}

func NewRedmineProvider() *redmineImportProvider {
	return &redmineImportProvider{}
}
