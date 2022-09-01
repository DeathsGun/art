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

func (u *Untis) GetClassList(request ClassListRequest) (*[]ClassListEntry, error) {
	result := &[]ClassListEntry{}
	err := u.client.Call("getKlassen", request, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (u *Untis) GetTimetable(options TimetableOptions) (*[]TimetablePeriod, error) {
	result := &[]TimetablePeriod{}
	err := u.client.Call("getTimetable", TimetableRequest{Options: options}, result)
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
