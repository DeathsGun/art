package redmine

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"
)

type Redmine struct {
	Url           string
	Client        http.Client
	Authorization RedmineAuthorization
}

func (r *Redmine) GetAPIKey() (*AccountResponse, error) {
	result := &AccountResponse{}
	err := r.CreateHTTPRequest("/my/account.json", url.Values{}, result)
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
	return &Redmine{
		Url: baseURL,
		Client: http.Client{
			Jar:     jar,
			Timeout: 5 * time.Second,
		},
		Authorization: *authorization,
	}, nil
}
