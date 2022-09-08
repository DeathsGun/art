package redmine

import (
	"fmt"
	"github.com/deathsgun/art/login"
	"github.com/deathsgun/art/provider"
	"github.com/deathsgun/art/utils"
	"time"
)

type redmineImportProvider struct {
}

const redmineURL = "https://joope.de/"
const printDate = false

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
	username, password := login.GetLogin(r.Name())
	if username == "" && password == "" {
		return nil, provider.ErrNoLoginConfigured
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
		return nil, nil
	}
	fmt.Printf("[redmine] Got %d entries from %s - %s\n",
		len(entries.TimeEntries),
		startDate.Format("02.01.2006"),
		endDate.Format("02.01.2006"),
	)

	sortedEntries := map[string][]TimeEntry{}
	for _, entry := range entries.TimeEntries {
		value, ok := sortedEntries[entry.Project.Name]
		if !ok {
			value = []TimeEntry{}
		}
		value = append(value, entry)
		sortedEntries[entry.Project.Name] = value
	}
	var addedKeys []string
	var result []*provider.Entry
	for key, sameTimes := range sortedEntries {
		for _, entry := range sameTimes {
			startDateTime, err := time.Parse("2006-01-02", entry.SpentOn)
			if err != nil {
				return nil, err
			}
			if !utils.Contains(addedKeys, key) {
				result = append(result, &provider.Entry{
					Date:     startDateTime,
					Text:     key + ":",
					Category: provider.ACTIVITY,
				})
				addedKeys = append(addedKeys, key)
			}
			result = append(result, &provider.Entry{
				Date:     startDateTime.Add(time.Hour * 10),
				Text:     "\t" + entry.Comments,
				Category: provider.ACTIVITY,
			})
		}
	}
	return result, nil
}

func NewRedmineProvider() *redmineImportProvider {
	return &redmineImportProvider{}
}
