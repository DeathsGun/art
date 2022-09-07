package redmine

import (
	"fmt"
	"github.com/deathsgun/art/login"
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

	startDateZero := startDate.Add(time.Hour * time.Duration(-startDate.Hour())).
		Add(time.Minute * time.Duration(-startDate.Minute())).
		Add(time.Second * time.Duration(-startDate.Second()))

	mondayDiff := time.Monday - startDateZero.Weekday()
	monday := startDateZero.Add(time.Hour * time.Duration(24*mondayDiff))

	fmt.Printf("%s", monday.Format(time.UnixDate))

	fmt.Printf("[Redmine] Importing for %s as start date\n", startDate.Format(time.RFC3339))
	username, password := login.GetLogin(r.Name())
	ra, err := NewRedmineAPI(redmineURL, AuthorizeHTTP(username, password))
	if err != nil {
		return nil, err
	}

	apiKey, err := ra.GetAccountInformation()
	var result []*provider.Entry
	result = append(result, &provider.Entry{Category: provider.ACTIVITY, Date: startDate, Text: apiKey.User.ApiKey})
	return result, nil
}

func NewRedmineProvider() *redmineImportProvider {
	return &redmineImportProvider{}
}
