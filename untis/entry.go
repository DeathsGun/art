package untis

type CalendarResponse struct {
	CalendarEntries []CalendarEntry `json:"calendarEntries"`
}

type CalendarEntry struct {
	StartDateTime   string   `json:"startDateTime"`
	Status          string   `json:"status"`
	TeachingContent string   `json:"teachingContent"`
	Subject         *Subject `json:"subject"`
	LessonInfo      string   `json:"lessonInfo"`
}

type Subject struct {
	DisplayName string `json:"displayName"`
	Id          int    `json:"id"`
	LongName    string `json:"longName"`
	ShortName   string `json:"shortName"`
}
