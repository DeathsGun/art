package untis

type CalendarResponse struct {
	CalendarEntries []CalendarEntry `json:"calendarEntries"`
}

type CalendarEntry struct {
	EndDateTime     string        `json:"endDateTime"`
	Id              int           `json:"id"`
	Classes         []Class       `json:"klasses"`
	Lesson          Lesson        `json:"lesson"`
	LessonInfo      string        `json:"lessonInfo"`
	Permissions     []string      `json:"permissions"`
	Rooms           []Room        `json:"rooms"`
	SingleEntries   []SingleEntry `json:"singleEntries"`
	StartDateTime   string        `json:"startDateTime"`
	Status          string        `json:"status"`
	SubType         SubType       `json:"subType"`
	Subject         Subject       `json:"subject"`
	SubstText       string        `json:"substText"`
	Teachers        []Teacher     `json:"teachers"`
	TeachingContent string        `json:"teachingContent"`
	Type            string        `json:"type"`
}

type SingleEntry struct {
	CreatedAt            interface{}   `json:"createdAt"`
	EndDateTime          string        `json:"endDateTime"`
	Id                   int           `json:"id"`
	LastUpdate           interface{}   `json:"lastUpdate"`
	Permissions          []string      `json:"permissions"`
	StartDateTime        string        `json:"startDateTime"`
	TeachingContent      interface{}   `json:"teachingContent"`
	TeachingContentFiles []interface{} `json:"teachingContentFiles"`
}

type Subject *struct {
	DisplayName  string `json:"displayName"`
	HasTimetable bool   `json:"hasTimetable"`
	Id           int    `json:"id"`
	LongName     string `json:"longName"`
	ShortName    string `json:"shortName"`
}

type Teacher struct {
	DisplayName  string      `json:"displayName"`
	HasTimetable bool        `json:"hasTimetable"`
	Id           int         `json:"id"`
	LongName     string      `json:"longName"`
	ShortName    string      `json:"shortName"`
	Status       string      `json:"status"`
	ImageUrl     interface{} `json:"imageUrl"`
}

type SubType struct {
	DisplayInPeriodDetails bool   `json:"displayInPeriodDetails"`
	DisplayName            string `json:"displayName"`
	Id                     int    `json:"id"`
}

type Room struct {
	DisplayName  string `json:"displayName"`
	HasTimetable bool   `json:"hasTimetable"`
	Id           int    `json:"id"`
	LongName     string `json:"longName"`
	ShortName    string `json:"shortName"`
	Status       string `json:"status"`
}

type Lesson struct {
	LessonId     int `json:"lessonId"`
	LessonNumber int `json:"lessonNumber"`
}

type Class struct {
	DisplayName  string `json:"displayName"`
	HasTimetable bool   `json:"hasTimetable"`
	Id           int    `json:"id"`
	LongName     string `json:"longName"`
	ShortName    string `json:"shortName"`
}
