package redmine

import (
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
	Authorization RedmineAuthorization
	RedmineUser   RedmineUser
}

func (r *Redmine) GetAccountInformation() (*AccountResponse, error) {
	result := &AccountResponse{}
	err := r.CreateHTTPRequest("/my/account.json", url.Values{}, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *Redmine) GetIssues(limit int, page int) (*IssueResponse, error) {
	result := &IssueResponse{}
	err := r.CreateHTTPRequest("/issues.json", url.Values{
		"limit":  []string{strconv.Itoa(limit)},
		"offset": []string{strconv.Itoa(limit * page)},
	}, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *Redmine) GetIssueDetails(limit int, page int, issueIds ...int) (*IssueResponse, error) {
	result := &IssueResponse{}
	s, err := json.Marshal(issueIds)
	if err != nil {
		return nil, err
	}
	err = r.CreateHTTPRequest("/issues.json", url.Values{
		"limit":    []string{strconv.Itoa(limit)},
		"offset":   []string{strconv.Itoa(limit * page)},
		"issue_id": []string{strings.Trim(string(s), "[]")},
	}, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *Redmine) GetTimeEntries(limit int, page int, from time.Time, to time.Time) (*TimeEntriesResponse, error) {
	result := &TimeEntriesResponse{}
	err := r.CreateHTTPRequest("/time_entries.json", url.Values{
		"limit":   []string{strconv.Itoa(limit)},
		"offset":  []string{strconv.Itoa(limit * page)},
		"user_id": []string{strconv.Itoa(r.RedmineUser.Id)},
		"from":    []string{from.Format("2006-01-02")},
		"to":      []string{to.Format("2006-01-02")},
	}, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *Redmine) CreateHTTPRequest(endpoint string, params url.Values, result any) error {
	endpointURL, err := url.Parse(r.Url + endpoint)
	if err != nil {
		return err
	}
	endpointURL.RawQuery = params.Encode()
	req, err := http.NewRequest(http.MethodGet, endpointURL.String(), nil)
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

func NewRedmineAPI(baseURL string, authorization *RedmineAuthorization) (*Redmine, error) {
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
		Authorization: *authorization,
	}
	accountResponse, err := redmine.GetAccountInformation()
	if err != nil {
		return nil, err
	}
	redmine.RedmineUser = accountResponse.User
	return redmine, nil
}
