package provider

import "time"

type Category int

const (
	ACTIVITY Category = iota
	TRAINING
	SUBJECTS
)

type Report struct {
	Id      int
	Entries []*Entry
}

type Entry struct {
	Date     time.Time
	Text     string
	Category Category
}
