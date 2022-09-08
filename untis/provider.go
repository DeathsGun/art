package untis

import (
	"fmt"
	"github.com/deathsgun/art/login"
	"github.com/deathsgun/art/provider"
	"regexp"
	"sort"
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
	un, err := NewUntisAPI("bk-ahaus")
	if err != nil {
		return nil, err
	}
	username, password := login.GetLogin(u.Name())
	if username == "" || password == "" {
		return nil, provider.ErrNoLoginConfigured
	}

	err = un.Login(username, password)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = un.Logout()
	}()

	var result []*provider.Entry
	for i := 0; i < 5; i++ {
		date := startDate.Add(time.Duration(i*24) * time.Hour)
		entries, err := u.getEntriesForDate(un, date)
		if err != nil {
			return nil, err
		}
		if len(entries) > 0 {
			fmt.Printf("Got %d entries from \n")
		}
		result = append(result, entries...)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Date.Before(result[j].Date)
	})
	return result, nil
}

func (u *untisImportProvider) getEntriesForDate(un *Untis, date time.Time) ([]*provider.Entry, error) {
	entries, err := un.Details(date, date.Add(23*time.Hour))
	if err != nil {
		return nil, err
	}
	if len(entries) == 0 {
		return nil, nil
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
				Date:      startDateTime,
				Text:      message,
				Category:  provider.SUBJECTS,
				PrintDate: true,
			})
		}
	}
	return result, nil
}

func NewUntisProvider() provider.ImportProvider {
	return &untisImportProvider{}
}
