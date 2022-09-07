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
	testEntry := &provider.Entry{
		Date:     startDate,
		Category: provider.ACTIVITY,
		Text:     "Kyshaftiger Test um zu checken obs klappt // Spoiler: Ich musste das fett fixen @Linus oder @Gerhard // Excel-Provider",
	}
	var result []*provider.Entry
	result = append(result, testEntry)
	return result, nil
}

func NewRedmineProvider() *redmineImportProvider {
	return &redmineImportProvider{}
}
