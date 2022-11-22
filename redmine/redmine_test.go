package redmine

import (
	"context"
	"os"
	"testing"
	"time"
)

func PrepareRedmineAPI(t *testing.T) *Redmine {
	if os.Getenv("REDMINE_USER") == "" {
		t.SkipNow()
	}
	redmine, err := NewRedmineAPI(context.Background(), "https://joope.de", &Authorization{
		User:     os.Getenv("REDMINE_USER"),
		Password: os.Getenv("REDMINE_PASSWORD"),
	})
	if err != nil {
		t.Fatal(err)
	}
	return redmine
}

func TestRedmineLogin(t *testing.T) {
	redmine := PrepareRedmineAPI(t)
	if redmine.Authorization.User == "" {
		t.Fatal("User is empty so login failed")
	}
}

func TestRedmine_GetTimeEntries(t *testing.T) {
	redmine := PrepareRedmineAPI(t)
	timeEntries, err := redmine.GetTimeEntries(context.Background(), 100, 0, time.Now().Add(-time.Hour*24*60), time.Now())
	if err != nil {
		t.Fatal(err)
	}
	if len(timeEntries.TimeEntries) > 0 {
		t.Fatal("time entries are empty")
	}
}
