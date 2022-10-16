package report

import (
	"bytes"
	"fmt"
	"github.com/deathsgun/art/utils"
)

type Report struct {
	Entries []Entry
}

func (r *Report) Format(printDate bool) (map[Category][]byte, error) {
	buffers := map[Category]*bytes.Buffer{
		Activity: {},
		Subjects: {},
		Training: {},
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
		Activity: buffers[Activity].Bytes(),
		Subjects: buffers[Subjects].Bytes(),
		Training: buffers[Training].Bytes(),
	}, nil
}
