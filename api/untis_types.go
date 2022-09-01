package api

type Person int
type PeriodType string
type PeriodState string

const (
	Teacher  Person = 2
	Students Person = 5
)

const (
	Lesson           PeriodType = ""
	OfficeHour       PeriodType = "oh"
	Standby          PeriodType = "sb"
	BreakSupervision PeriodType = "bs"
	Examination      PeriodType = "ex"
)

const (
	Default   PeriodState = ""
	Cancelled PeriodState = "cancelled"
	Irregular PeriodState = "irregular"
)

type TimetableRequest struct {
	Options TimetableOptions `json:"options"`
}

type TimetableOptions struct {
	Element           TimetableRequestElement `json:"element"`
	StartDate         int                     `json:"startDate,omitempty"`
	EndDate           int                     `json:"endDate,omitempty"`
	OnlyBaseTimetable bool                    `json:"only_base_timetable,omitempty"`
	ShowBookings      bool                    `json:"showBooking,omitempty"`
	ShowInfo          bool                    `json:"showInfo,omitempty"`
	ShowSubstText     bool                    `json:"showSubstText,omitempty"`
	ShowLessonText    bool                    `json:"showLsText,omitempty"`
	ShowLessonNumber  bool                    `json:"showLsNumber,omitempty"`
	ShowStudentGroup  bool                    `json:"showStudentgroup,omitempty"`
	ClassFields       []string                `json:"klasseFields,omitempty"`
	RoomFields        []string                `json:"roomFields,omitempty"`
	SubjectFields     []string                `json:"subjectFields,omitempty"`
	TeacherFields     []string                `json:"teacherFields"`
}

type TimetableRequestElement struct {
	Id      string `json:"id"`
	Type    int    `json:"type"`
	KeyType string `json:"keyType"`
}

type TimetablePeriod struct {
	Id               int                       `json:"id"`
	Date             int                       `json:"date,omitempty"`
	StartDate        int                       `json:"startDate,omitempty"`
	EndDate          int                       `json:"endDate,omitempty"`
	Type             PeriodType                `json:"lstype"`
	State            PeriodState               `json:"code"`
	Info             string                    `json:"info,omitempty"`
	SubstitutionText string                    `json:"substText,omitempty"`
	Text             string                    `json:"lstext,omitempty"`
	PeriodNumber     int                       `json:"lsnumber,omitempty"`
	StatsFlags       string                    `json:"statflags"`
	ActivityType     string                    `json:"activityType"`
	StudentGroup     string                    `json:"sg,omitempty"`
	BookingRemark    string                    `json:"bkRemark,omitempty"`
	BookingText      string                    `json:"bkText,omitempty"`
	Classes          []TimetablePeriodListItem `json:"kl"`
	Teachers         []TimetablePeriodListItem `json:"teachers"`
	Subjects         []TimetablePeriodListItem `json:"subjects"`
	Rooms            []TimetablePeriodListItem `json:"rooms"`
}

type TimetablePeriodListItem struct {
	Id          int    `json:"id"`
	Name        string `json:"name,omitempty"`
	LongName    string `json:"longname,omitempty"`
	ExternalKey string `json:"externalKey,omitempty"`
}

type ClassListRequest struct {
	SchoolYearId string `json:"schoolyearId,omitempty"`
}

type ClassListEntry struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	LongName  string `json:"longName"`
	ForeColor int    `json:"foreColor"`
	BackColor int    `json:"backColor"`
}

type LoginResponse struct {
	SessionId  string `json:"sessionId"`
	PersonType Person `json:"personType"`
	PersonId   int    `json:"personId"`
}

type LoginRequest struct {
	Username string `json:"user"`
	Password string `json:"password"`
	Client   string `json:"client"`
}
