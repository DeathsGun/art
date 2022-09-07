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

func (u *untisImportProvider) ValidateLogin(username string, password string) error {
	un, err := NewUntisAPI("bk-ahaus")
	if err != nil {
		return err
	}
	err = un.Login(username, password)
	if err != nil {
		return err
	}
	return un.Logout()
}

func (u *untisImportProvider) NeedsLogin() bool {
	return true
}

var numberRegex = regexp.MustCompile("\\d+")

func (u *untisImportProvider) Import(startDate string) ([]*provider.Entry, error) {
	t := time.Now()
	t = t.Add(time.Hour * time.Duration(-t.Hour())).
		Add(time.Minute * time.Duration(-t.Minute())).
		Add(time.Second * time.Duration(-t.Second()))
	if startDate != "" {
		tim, err := time.Parse("02.01.06", startDate)
		if err != nil {
			return nil, fmt.Errorf("wrong format expected format dd.mm.yy e.g. 02.01.06: %s", err)
		}
		t = tim
	}
	fmt.Printf("Importing for %s as start date\n", t.Format(time.RFC3339))
	un, err := NewUntisAPI("bk-ahaus")
	if err != nil {
		return nil, err
	}
	username, password := login.GetLogin("untis")
	if username == "" || password == "" {
		return nil, errors.New("untis login not configured")
	}

	err = un.Login(username, password)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = un.Logout()
	}()

	entries, err := un.Details(t, t.Add(23*time.Hour))
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
