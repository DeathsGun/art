package untis

import (
	"errors"
	"fmt"
	"github.com/deathsgun/art/login"
	"github.com/deathsgun/art/provider"
	"regexp"
	"time"
)

type untisImportProvider struct {
}

func (u *untisImportProvider) Name() string {
	return "untis"
}

func (u *untisImportProvider) ValidateLogin(username string, password string) (string, string, error) {
	un, err := NewUntisAPI("bk-ahaus")
	if err != nil {
		return "", "", err
	}
	err = un.Login(username, password)
	if err != nil {
		return "", "", err
	}
	return username, password, un.Logout()
}

func (u *untisImportProvider) NeedsLogin() bool {
	return true
}

var numberRegex = regexp.MustCompile("\\d+")

func (u *untisImportProvider) Import(startDate time.Time) ([]*provider.Entry, error) {
	startDate = startDate.Add(time.Hour * time.Duration(-startDate.Hour())).
		Add(time.Minute * time.Duration(-startDate.Minute())).
		Add(time.Second * time.Duration(-startDate.Second()))
	fmt.Printf("[Untis] Importing for %s as start date\n", startDate.Format(time.RFC3339))
	un, err := NewUntisAPI("bk-ahaus")
	if err != nil {
		return nil, err
	}
	username, password := login.GetLogin(u.Name())
	if username == "" || password == "" {
		return nil, errors.New(fmt.Sprintf("%s login not configured", u.Name()))
	}

	err = un.Login(username, password)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = un.Logout()
	}()

	entries, err := un.Details(startDate, startDate.Add(23*time.Hour))
	if err != nil {
		return nil, err
	}
	if len(entries) == 0 {
		return nil, fmt.Errorf("no entries for date %s", startDate)
	}

	sortedEntries := map[string][]CalendarEntry{}
	for _, entry := range entries {
		value, ok := sortedEntries[entry.StartDateTime]
		if !ok {
			value = []CalendarEntry{}
		}
		value = append(value, entry)
		sortedEntries[entry.StartDateTime] = value
	}
	var result []*provider.Entry
	for _, sameTimes := range sortedEntries {
		for _, entry := range sameTimes {
			startDateTime, err := time.Parse("2006-01-02T15:04:05", entry.StartDateTime)
			if err != nil {
				return nil, err
			}
			message := "Lehrer hat noch nicht eingetragen !!BITTE FÜLLEN!!"
			if entry.Status == "CANCELLED" && len(sameTimes) == 2 {
				continue
			} else if entry.Status == "CANCELLED" && len(sameTimes) == 1 {
				message = "Entfall"
			}
			if entry.TeachingContent != "" {
				message = entry.TeachingContent
			}
			if entry.Subject != nil {
				message = numberRegex.ReplaceAllString(entry.Subject.ShortName, "") + ": " + message
			} else {
				message = entry.LessonInfo + ": " + message
			}

			result = append(result, &provider.Entry{
				Date:     startDateTime,
				Text:     message,
				Category: provider.SUBJECTS,
			})
		}
	}
	return result, nil
}

func NewUntisProvider() provider.ImportProvider {
	return &untisImportProvider{}
}
