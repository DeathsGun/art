package redmine

import (
	"time"
)

func LerpToPreviousMonday(t time.Time) time.Time {
	return t.Add(-time.Second * time.Duration(t.Second())).Add(-time.Minute * time.Duration(t.Minute())).Add(-time.Hour * time.Duration(t.Hour())).Add(time.Hour * 24 * time.Duration(time.Monday-t.Weekday()))
}
