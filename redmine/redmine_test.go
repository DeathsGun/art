package redmine

import (
	"fmt"
	"os"
	"testing"
)

func TestRedmine(t *testing.T) {
	redmine, err := NewRedmineAPI("https://joope.de", AuthorizeHTTP(os.Getenv("REDMINE_USER"), os.Getenv("REDMINE_PASSWORD")))
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%+v", redmine.RedmineUser)
}
