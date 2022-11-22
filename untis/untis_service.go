package untis

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/deathsgun/art/rpc"
	"github.com/deathsgun/art/untis/dto"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type IUntisService interface {
	Login(ctx context.Context, school *School, username string, password string) (*Session, error)
	Logout(ctx context.Context, session *Session) error
	ValidateLogin(ctx context.Context, session *Session) error
	SearchSchools(ctx context.Context, query string) ([]School, error)
	GetCalendarForDay(ctx context.Context, startTime time.Time, endTime time.Time) ([]CalendarEntry, error)
}

type service struct {
}

func (s *service) SearchSchools(ctx context.Context, query string) ([]School, error) {
	body := fmt.Sprintf("{\"id\": \"wu_schulsuche-%d\", \"method\": \"searchSchool\", \"params\": [{\"search\": \"%s\"}], \"jsonrpc\": \"2.0\"}", time.Now().UTC().UnixMilli(), query)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://schoolsearch.webuntis.com/schoolquery2", strings.NewReader(body))
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	var result []School
	response := &dto.SchoolResponse{}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid status code: %d", resp.StatusCode)
	}
	err = json.NewDecoder(resp.Body).Decode(response)
	if err != nil {
		return nil, err
	}
	for _, school := range response.Result.Schools {
		result = append(result, School{
			Server:      school.Server,
			Address:     school.Address,
			DisplayName: school.DisplayName,
			LoginName:   school.LoginName,
		})
	}
	return result, nil
}

func (s *service) Login(ctx context.Context, school *School, username, password string) (*Session, error) {
	endpoint := fmt.Sprintf("https://%s/WebUntis", school.Server)
	client := rpc.New(fmt.Sprintf("%s/jsonrpc.do?school=%s", endpoint, school.LoginName), "")
	response := &dto.LoginResponse{}
	err := client.Call(ctx, "authenticate", dto.LoginRequest{
		Username: username,
		Password: password,
	}, response)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint+"/api/token/new", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.AddCookie(&http.Cookie{
		Name:  "JSESSIONID",
		Value: response.SessionId,
	})
	resp, err := http.DefaultClient.Do(req)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid status %d", resp.StatusCode)
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return &Session{
		Endpoint:   endpoint,
		SessionId:  response.SessionId,
		School:     school.LoginName,
		Token:      string(data),
		PersonType: response.PersonType,
		PersonId:   response.PersonId,
	}, nil
}

func (s *service) Logout(ctx context.Context, session *Session) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, session.Endpoint+"/saml/logout", nil)
	if err != nil {
		return err
	}
	req.AddCookie(&http.Cookie{
		Name:  "JSESSIONID",
		Value: session.SessionId,
	})
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusFound {
		return fmt.Errorf("invalid status code: %d", resp.StatusCode)
	}
	return nil
}

func (s *service) ValidateLogin(ctx context.Context, session *Session) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/api/rest/view/v1/timegrid", session.Endpoint), nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+session.Token)
	req.Header.Set("Accept", "application/json")
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("invalid status code: %d", resp.StatusCode)
	}
	return nil
}

func (s *service) GetCalendarForDay(ctx context.Context, startTime time.Time, endTime time.Time) ([]CalendarEntry, error) {
	session := ctx.Value("session").(*Session)
	detailsUrl, err := url.Parse(session.Endpoint + "/api/rest/view/v1/calendar-entry/detail")
	if err != nil {
		return nil, err
	}
	query := url.Values{
		"elementId":     []string{fmt.Sprintf("%d", session.PersonId)},
		"elementType":   []string{fmt.Sprintf("%d", session.PersonType)},
		"startDateTime": []string{startTime.Format("2006-01-02T15:04:05")},
		"endDateTime":   []string{endTime.Format("2006-01-02T15:04:05")},
	}
	detailsUrl.RawQuery = query.Encode()
	req, err := http.NewRequest(http.MethodGet, detailsUrl.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", session.Token))

	resp, err := http.DefaultClient.Do(req)
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

func NewService() IUntisService {
	return &service{}
}
