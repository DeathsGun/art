package redmine

import (
	"fmt"
	"testing"
)

func TestRedmine(t *testing.T) {
	redmine, err := NewRedmineAPI("https://joope.de", AuthorizeHTTP("", ""))
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%+v", redmine.RedmineUser)
}
