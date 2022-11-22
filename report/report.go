package report

import (
	"bytes"
	"context"
	"fmt"
	"github.com/deathsgun/art/di"
	"github.com/deathsgun/art/i18n"
	"github.com/deathsgun/art/utils"
	"strings"
)

var weekdays = []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}

type Report struct {
	Entries []Entry
}

func (r *Report) Format(ctx context.Context, printDate bool) (map[Category][]byte, error) {
	buffers := map[Category]*bytes.Buffer{
		Activity: {},
		Subjects: {},
		Training: {},
	}
	translationService := di.Instance[i18n.ITranslationService]("i18n")
	var dates []string
	for _, v := range r.Entries {
		date := v.Date.Format("Monday 02.01.2006:")
		for _, weekday := range weekdays {
			translated := translationService.Translate(ctx, strings.ToUpper(weekday))
			date = strings.ReplaceAll(date, weekday, translated)
		}
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
