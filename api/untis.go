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

func (u *Untis) FindPersonId(element ElementType, foreName string, surName string, dayOfBirth ...int) (PersonSearchResponse, error) {
	dob := 0
	if len(dayOfBirth) > 0 {
		dob = dayOfBirth[0]
	}
	var result = PersonSearchResponse(0)
	err := u.client.Call("getPersonId", &PersonSearchRequest{
		Type:       element,
		ForeName:   foreName,
		SurName:    surName,
		DayOfBirth: dob,
	}, &result)
	if err != nil {
		return 0, err
	}
	return result, nil
}

func (u *Untis) GetCurrentSchoolYear() (*SchoolYear, error) {
	result := &SchoolYear{}
	err := u.client.Call("getCurrentSchoolyear", nil, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (u *Untis) GetSchoolYears() (*[]SchoolYear, error) {
	result := &[]SchoolYear{}
	err := u.client.Call("getSchoolyears", nil, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (u *Untis) GetTimeGrids() (*[]TimeGrid, error) {
	result := &[]TimeGrid{}
	err := u.client.Call("getTimegridUnits", nil, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (u *Untis) GetHolidays() (*[]Holiday, error) {
	result := &[]Holiday{}
	err := u.client.Call("getHolidays", nil, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (u *Untis) GetSubjects() (*[]Subject, error) {
	result := &[]Subject{}
	err := u.client.Call("getSubjects", nil, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (u *Untis) GetDepartments() (*[]Department, error) {
	result := &[]Department{}
	err := u.client.Call("getDepartments", nil, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (u *Untis) GetRooms() (*[]Room, error) {
	result := &[]Room{}
	err := u.client.Call("getRooms", nil, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (u *Untis) GetStudents() (*[]Student, error) {
	result := &[]Student{}
	err := u.client.Call("getStudents", nil, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (u *Untis) GetClassList(schoolYearId ...string) (*[]ClassListEntry, error) {
	syId := ""
	if len(schoolYearId) > 0 {
		syId = schoolYearId[0]
	}
	result := &[]ClassListEntry{}
	err := u.client.Call("getKlassen", ClassListRequest{SchoolYearId: syId}, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (u *Untis) GetTimetable(options *TimetableOptions) (*[]TimetablePeriod, error) {
	result := &[]TimetablePeriod{}
	err := u.client.Call("getTimetable", TimetableRequest{Options: *options}, result)
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
