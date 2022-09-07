package untis

import (
	"encoding/json"
	"fmt"
	"github.com/deathsgun/art/rpc"
	"github.com/google/uuid"
	"io"
	"net/http"
	"net/url"
	"time"
)

const BaseUrl = "https://asopo.webuntis.com/WebUntis"

type Untis struct {
	client        *rpc.Client
	loginResponse LoginResponse
	token         string
}

func (u *Untis) Login(username string, password string) error {
	loginResponse := &LoginResponse{}
	err := u.client.Call("authenticate", LoginRequest{Username: username, Password: password, Client: uuid.NewString()}, loginResponse)
	resp, err := u.client.Client().Get(BaseUrl + "/api/token/new")
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received invalid status code: %d", resp.StatusCode)
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	u.token = string(data)
	u.loginResponse = *loginResponse
	return nil
}

func (u *Untis) Logout() error {
	return u.client.Call("logout", nil, nil)
}

func (u *Untis) Details(startTime time.Time, endTime time.Time) ([]CalendarEntry, error) {
	detailsUrl, err := url.Parse(BaseUrl + "/api/rest/view/v1/calendar-entry/detail")
	if err != nil {
		return nil, err
	}
	query := url.Values{
		"elementId":     []string{fmt.Sprintf("%d", u.loginResponse.PersonId)},
		"elementType":   []string{fmt.Sprintf("%d", u.loginResponse.PersonType)},
		"startDateTime": []string{startTime.Format("2006-01-02T15:04:05")},
		"endDateTime":   []string{endTime.Format("2006-01-02T15:04:05")},
	}
	detailsUrl.RawQuery = query.Encode()
	req, err := http.NewRequest(http.MethodGet, detailsUrl.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", u.token))

	resp, err := u.client.Client().Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	response := &CalendarResponse{}
	err = json.NewDecoder(resp.Body).Decode(response)
	if err != nil {
		return nil, err
	}
	return response.CalendarEntries, nil
}

func NewUntisAPI(schoolName string) (*Untis, error) {
	client, err := rpc.Dial(BaseUrl + "/jsonrpc.do?school=" + schoolName)
	if err != nil {
		return nil, err
	}
	return &Untis{client: client}, nil
}
