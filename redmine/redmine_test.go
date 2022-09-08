package redmine

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func PrepareRedmineAPI(t *testing.T) *Redmine {
	redmine, err := NewRedmineAPI("https://joope.de", AuthorizeHTTP(os.Getenv("REDMINE_USER"), os.Getenv("REDMINE_PASSWORD")))
	if err != nil {
		t.Fatal(err)
	}
	return redmine
}

func TestRedmineLogin(t *testing.T) {
	redmine := PrepareRedmineAPI(t)
	fmt.Printf("%+v", redmine.RedmineUser)
}

func TestRedmine_GetTimeEntries(t *testing.T) {
	redmine := PrepareRedmineAPI(t)
	timeEntries, err := redmine.GetTimeEntries(100, 0, time.Now().Add(-time.Hour*24*60), time.Now())
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%+v", timeEntries)
}

func TestRedmineGetIssues(t *testing.T) {
	redmine := PrepareRedmineAPI(t)
	issuesPage1, err := redmine.GetIssues(5, 1)
	if err != nil {
		t.Fatal(err)
	}
	issuesPage2, err := redmine.GetIssues(5, 2)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%+v", issuesPage1)
	fmt.Printf("%+v", issuesPage2)
}
