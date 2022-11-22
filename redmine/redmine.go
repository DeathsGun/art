package redmine

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Redmine struct {
	Url           string
	Client        http.Client
	Authorization *Authorization
	User          *User
}

func (r *Redmine) GetAccountInformation(ctx context.Context) (*AccountResponse, error) {
	result := &AccountResponse{}
	err := r.CreateHTTPRequest(ctx, "/my/account.json", url.Values{}, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *Redmine) GetIssues(ctx context.Context, limit int, page int) (*IssueResponse, error) {
	result := &IssueResponse{}
	err := r.CreateHTTPRequest(ctx, "/issues.json", url.Values{
		"limit":  []string{strconv.Itoa(limit)},
		"offset": []string{strconv.Itoa(limit * page)},
	}, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *Redmine) GetIssueDetails(ctx context.Context, limit int, page int, issueIds ...int) (*IssueResponse, error) {
	result := &IssueResponse{}
	s, err := json.Marshal(issueIds)
	if err != nil {
		return nil, err
	}
	err = r.CreateHTTPRequest(ctx, "/issues.json", url.Values{
		"limit":    []string{strconv.Itoa(limit)},
		"offset":   []string{strconv.Itoa(limit * page)},
		"issue_id": []string{strings.Trim(string(s), "[]")},
	}, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *Redmine) GetTimeEntries(ctx context.Context, limit int, page int, from time.Time, to time.Time) (*TimeEntriesResponse, error) {
	result := &TimeEntriesResponse{}
	err := r.CreateHTTPRequest(ctx, "/time_entries.json", url.Values{
		"limit":   []string{strconv.Itoa(limit)},
		"offset":  []string{strconv.Itoa(limit * page)},
		"user_id": []string{strconv.Itoa(r.User.Id)},
		"from":    []string{from.Format("2006-01-02")},
		"to":      []string{to.Format("2006-01-02")},
	}, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *Redmine) CreateHTTPRequest(ctx context.Context, endpoint string, params url.Values, result any) error {
	endpointURL, err := url.Parse(r.Url + endpoint)
	if err != nil {
		return err
	}
	endpointURL.RawQuery = params.Encode()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpointURL.String(), nil)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", r.Authorization.BuildAuthorizationHeader())

	resp, err := r.Client.Do(req)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return err
	}
	return nil
}

func NewRedmineAPI(ctx context.Context, baseURL string, authorization *Authorization) (*Redmine, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}
	redmine := &Redmine{
		Url: baseURL,
		Client: http.Client{
			Jar:     jar,
			Timeout: 5 * time.Second,
		},
		Authorization: authorization,
	}
	accountResponse, err := redmine.GetAccountInformation(ctx)
	if err != nil {
		return nil, err
	}
	redmine.User = &accountResponse.User
	return redmine, nil
}
