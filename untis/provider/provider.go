package provider

import (
	"context"
	"fmt"
	"github.com/deathsgun/art/auth"
	"github.com/deathsgun/art/config/model"
	"github.com/deathsgun/art/di"
	"github.com/deathsgun/art/i18n"
	"github.com/deathsgun/art/provider"
	"github.com/deathsgun/art/report"
	"github.com/deathsgun/art/untis"
	"regexp"
	"sort"
	"time"
)

type impl struct {
}

func (u *impl) ValidateConfig(ctx context.Context, _ *model.ProviderConfig) error {
	service := di.Instance[untis.IUntisService]("untis")
	return service.ValidateLogin(ctx, auth.Session(ctx)) // Should be valid tho
}

func (u *impl) Logo() string {
	return "WebUntis_Logo_orange.png"
}

func (u *impl) Id() string {
	return "PROVIDER_UNTIS"
}

func (u *impl) Capabilities() []provider.Capability {
	return []provider.Capability{
		provider.TypeImport,
	}
}

func NewProvider() provider.ImportProvider {
	return &impl{}
}

func (u *impl) Import(ctx context.Context, monday time.Time) ([]report.Entry, error) {
	var result []report.Entry
	for i := 0; i < 5; i++ {
		date := monday.Add(time.Duration(i*24) * time.Hour)
		entries, err := u.getEntriesForDate(ctx, date)
		if err != nil {
			return nil, err
		}
		if len(entries) > 0 {
			fmt.Printf("[untis] Got %d entries from %s\n", len(entries), date.Format("02.01.2006"))
		}
		result = append(result, entries...)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Date.Before(result[j].Date)
	})
	return result, nil
}

var numberRegex = regexp.MustCompile("\\d+")

func (u *impl) getEntriesForDate(ctx context.Context, date time.Time) ([]report.Entry, error) {
	untisService := di.Instance[untis.IUntisService]("untis")
	entries, err := untisService.GetCalendarForDay(ctx, date, date.Add(23*time.Hour))
	if err != nil {
		return nil, err
	}
	if len(entries) == 0 {
		return nil, nil
	}

	sortedEntries := map[string][]untis.CalendarEntry{}
	for _, entry := range entries {
		value, ok := sortedEntries[entry.StartDateTime]
		if !ok {
			value = []untis.CalendarEntry{}
		}
		value = append(value, entry)
		sortedEntries[entry.StartDateTime] = value
	}
	var result []report.Entry
	for _, sameTimes := range sortedEntries {
		for _, entry := range sameTimes {
			startDateTime, err := time.Parse("2006-01-02T15:04:05", entry.StartDateTime)
			if err != nil {
				return nil, err
			}
			translationService := di.Instance[i18n.ITranslationService]("translation")
			message := translationService.Translate(ctx, "UNTIS_ENTRY_NOT_FILLED")
			if entry.Status == "CANCELLED" && len(sameTimes) == 2 {
				continue
			} else if entry.Status == "CANCELLED" && len(sameTimes) == 1 {
				message = translationService.Translate(ctx, "UNTIS_LESSON_CANCELLED")
			}
			if entry.TeachingContent != "" {
				message = entry.TeachingContent
			}
			if entry.Subject != nil {
				message = numberRegex.ReplaceAllString(entry.Subject.ShortName, "") + ": " + message
			} else {
				message = entry.LessonInfo + ": " + message
			}

			result = append(result, report.Entry{
				Date:     startDateTime,
				Text:     message,
				Category: report.Subjects,
			})
		}
	}
	return result, nil
}
