package report

import "time"

type Entry struct {
	Date     time.Time
	Text     string
	Category Category
}
