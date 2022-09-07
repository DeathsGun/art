package redmine

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"time"
)

type Redmine struct {
	Url           string
	Client        http.Client
	Authorization RedmineAuthorization
	RedmineUser   RedmineUser
}

func (r *Redmine) getAPIKey() (*AccountResponse, error) {
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

func (r *Redmine) GetRelevantIssues() (*IssueResponse, error) {
	return nil, nil
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
	accountResponse, err := redmine.getAPIKey()
	if err != nil {
		return nil, err
	}
	redmine.RedmineUser = accountResponse.User
	return redmine, nil
}
