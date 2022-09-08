package provider

import "time"

type Category int

func (c Category) Text() string {
	switch c {
	case ACTIVITY:
		return "Betriebliche Tätigkeiten"
	case SUBJECTS:
		return "Berufsschule (Unterrichtsthemen)"
	case TRAINING:
		return "Unterweisungen, betrieblicher Unterricht, sonstige Schulung"
	}
	return ""
}

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
	Date      time.Time
	Text      string
	Category  Category
	PrintDate bool
}
