package provider

import (
	"bytes"
	"fmt"
	"github.com/deathsgun/art/utils"
	"time"
)

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

func (r *Report) Format(printDate bool) (map[Category][]byte, error) {
	buffers := map[Category]*bytes.Buffer{
		ACTIVITY: {},
		SUBJECTS: {},
		TRAINING: {},
	}
	var dates []string
	for _, v := range r.Entries {
		date := v.Date.Format("Monday 02.01.2006:")
		if !utils.Contains(dates, date) {
			if printDate {
				_, err := buffers[v.Category].WriteString("\t" + date + "\n")
				if err != nil {
					return nil, nil
				}
			}
			dates = append(dates, date)
		}
		_, err := buffers[v.Category].WriteString(fmt.Sprintf("\t\t- %s\n", v.Text))
		if err != nil {
			return nil, err
		}
	}
	return map[Category][]byte{
		ACTIVITY: buffers[ACTIVITY].Bytes(),
		SUBJECTS: buffers[SUBJECTS].Bytes(),
		TRAINING: buffers[TRAINING].Bytes(),
	}, nil
}

type Entry struct {
	Date     time.Time
	Text     string
	Category Category
}
