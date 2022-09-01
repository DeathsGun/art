package api

import (
	"github.com/deathsgun/art/rpc"
	"github.com/google/uuid"
)

type Untis struct {
	client        *rpc.Client
	loginResponse LoginResponse
}

func (u *Untis) Login(username string, password string) error {
	loginResponse := &LoginResponse{}
	err := u.client.Call("authenticate", LoginRequest{Username: username, Password: password, Client: uuid.NewString()}, loginResponse)
	if err == nil {
		u.loginResponse = *loginResponse
	}
	return err
}

func (u *Untis) Logout() error {
	return u.client.Call("logout", nil, nil)
}

func (u *Untis) GetTimetable(request *TimetableRequest) (*[]TimetablePeriod, error) {
	result := &[]TimetablePeriod{}
	err := u.client.Call("getTimetable", request, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func NewUntisAPI(schoolName string) (*Untis, error) {
	client, err := rpc.Dial("https://asopo.webuntis.com/WebUntis/jsonrpc.do?school=" + schoolName)
	if err != nil {
		return nil, err
	}
	return &Untis{client: client}, nil
}

type Person int
type LessonType string
type PeriodState string

const (
	Teacher  Person = 2
	Students Person = 5
)

const (
	Lesson           LessonType = "ls"
	OfficeHour       LessonType = "oh"
	Standby          LessonType = "sb"
	BreakSupervision LessonType = "bs"
	Examination      LessonType = "ex"
)

const (
	Default   PeriodState = ""
	Cancelled PeriodState = "cancelled"
	Irregular PeriodState = "irregular"
)

type TimetablePeriod struct {
	Id           int              `json:"id"`
	Date         int              `json:"date"`
	StartDate    int              `json:"startDate"`
	EndDate      int              `json:"endDate"`
	Class        []PeriodListType `json:"kl"`
	Teachers     []PeriodListType `json:"te"`
	Subjects     []PeriodListType `json:"su"`
	Rooms        []PeriodListType `json:"ro"`
	Type         LessonType       `json:"lstype"`
	State        PeriodState      `json:"code"`
	Text         string           `json:"lstext"`
	Flags        string           `json:"statflags"`
	ActivityType string           `json:"activityType"`
}

type PeriodListType struct {
	Id int `json:"id"`
}

type LoginResponse struct {
	SessionId  string `json:"sessionId"`
	PersonType Person `json:"personType"`
	PersonId   int    `json:"personId"`
}

type TimetableRequest struct {
	Id        int `json:"id"`
	Type      int `json:"type"`
	StartDate int `json:"startDate"`
	EndDate   int `json:"endDate"`
}

type LoginRequest struct {
	Username string `json:"user"`
	Password string `json:"password"`
	Client   string `json:"client"`
}
