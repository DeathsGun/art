package main

import (
	"github.com/deathsgun/art/redmine"
	"testing"
	"time"
)

func TestTimeFunction(t *testing.T) {
	monday := redmine.LerpToPreviousMonday(time.Now())
	println(monday.Format(time.RFC850))
}
