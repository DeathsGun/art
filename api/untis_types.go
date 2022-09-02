package api

type ElementType int

const (
	LessonElement  ElementType = 1
	TeacherElement ElementType = 2
	SubjectElement ElementType = 3
	RoomElement    ElementType = 4
	StudentElement ElementType = 5
)

//region Person search

type PersonSearchRequest struct {
	Type       ElementType `json:"type"`
	ForeName   string      `json:"fn"`
	SurName    string      `json:"sn"`
	DayOfBirth int         `json:"dob"`
}

type PersonSearchResponse int

//endregion

//region SchoolYears

type SchoolYear struct {
	Id        int    `json:"id"`
	Name      string `json:"name,omitempty"`
	StartDate int    `json:"startDate,omitempty"`
	EndDate   int    `json:"endDate,omitempty"`
}

//endregion

//region TimeGrid

type Day int

const (
	Sunday = iota + 1
	Monday
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
)

type TimeGrid struct {
	Day       Day         `json:"day"`
	TimeUnits []TimeUnits `json:"timeUnits,omitempty"`
}

type TimeUnits struct {
	StartTime int `json:"startTime"`
	EndTime   int `json:"endTime"`
}

//endregion

//region Holidays

type Holiday struct {
	Id        int    `json:"id"`
	Name      string `json:"name,omitempty"`
	LongName  string `json:"longName,omitempty"`
	StartData int    `json:"startDate"`
	EndDate   int    `json:"endDate"`
}

//endregion

//region Department

type Department struct {
	Id       int    `json:"id"`
	Name     string `json:"name,omitempty"`
	LongName string `json:"longName,omitempty"`
}

//endregion

//region Rooms

type Room struct {
	Id        int    `json:"id"`
	Name      string `json:"name,omitempty"`
	LongName  string `json:"longName,omitempty"`
	ForeColor string `json:"foreColor,omitempty"`
	BackColor string `json:"backColor,omitempty"`
}

//endregion

//region Subjects

type Subject struct {
	Id        int    `json:"id"`
	LongName  string `json:"longName,omitempty"`
	ForeColor string `json:"foreColor"`
	BackColor string `json:"backColor"`
}

//endregion

//region Students

type Gender string

const (
	Male   = "male"
	Female = "female"
)

type Student struct {
	Id       int    `json:"id"`
	Name     string `json:"name,omitempty"`
	ForeName string `json:"foreName,omitempty"`
	LongName string `json:"longName,omitempty"`
	Gender   Gender `json:"gender,omitempty"`
}

//endregion

//region Timetable

type PeriodType string

const (
	Lesson           PeriodType = ""
	OfficeHour       PeriodType = "oh"
	Standby          PeriodType = "sb"
	BreakSupervision PeriodType = "bs"
	Examination      PeriodType = "ex"
)

type PeriodState string

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
	Teachers         []TimetablePeriodListItem `json:"te"`
	Subjects         []TimetablePeriodListItem `json:"su"`
	Rooms            []TimetablePeriodListItem `json:"ro"`
}

type TimetablePeriodListItem struct {
	Id          int    `json:"id"`
	Name        string `json:"name,omitempty"`
	LongName    string `json:"longname,omitempty"`
	ExternalKey string `json:"externalkey,omitempty"`
}

//endregion

//region ClassList

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

//endregion

//region Login

type LoginResponse struct {
	SessionId  string      `json:"sessionId"`
	PersonType ElementType `json:"personType"`
	PersonId   int         `json:"personId"`
}

type LoginRequest struct {
	Username string `json:"user"`
	Password string `json:"password"`
	Client   string `json:"client"`
}

//endregion
