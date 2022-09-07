package redmine

import (
	"fmt"
	"testing"
)

func TestRedmine(t *testing.T) {
	redmine, err := NewRedmineAPI("https://joope.de", AuthorizeHTTP("dh_schmidt_jona", "josm2019!"))
	if err != nil {
		t.Fatal(err)
	}
	account, err := redmine.GetAPIKey()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%+v", account)
}
