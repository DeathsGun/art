package redmine

import (
	"errors"
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

	startDate = LerpToPreviousMonday(startDateZero)

	fmt.Printf("[Redmine] Importing for %s as start date\n", startDate.Format(time.RFC3339))
	username, password := login.GetLogin(r.Name())
	if username == "" && password == "" {
		return nil, errors.New(fmt.Sprintf("%s login not configured", r.Name()))
	}
	ra, err := NewRedmineAPI(redmineURL, AuthorizeHTTP(username, password))
	if err != nil {
		return nil, err
	}
	endDate := startDate.Add(time.Hour * 24 * 5)
	entries, err := ra.GetTimeEntries(100, 0, startDate, endDate)
	if err != nil {
		return nil, err
	}
	if len(entries.TimeEntries) == 0 {
		return nil, fmt.Errorf("no entries for date %s", startDate)
	}

	sortedEntries := map[string][]TimeEntry{}
	for _, entry := range entries.TimeEntries {
		value, ok := sortedEntries[entry.SpentOn]
		if !ok {
			value = []TimeEntry{}
		}
		value = append(value, entry)
		sortedEntries[entry.SpentOn] = value
	}
	var result []*provider.Entry
	for _, sameTimes := range sortedEntries {
		for _, entry := range sameTimes {
			startDateTime, err := time.Parse("2006-01-02", entry.SpentOn)
			if err != nil {
				return nil, err
			}

			result = append(result, &provider.Entry{
				Date:      startDateTime,
				Text:      entry.Comments,
				Category:  provider.ACTIVITY,
				PrintDate: false,
			})
		}
	}
	return result, nil
}

func NewRedmineProvider() *redmineImportProvider {
	return &redmineImportProvider{}
}
