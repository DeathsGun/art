package ihk

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

// https://www.bildung-ihk-nordwestfalen.de/tibrosBB/azubiHome.jsp
// https://www.bildung-ihk-nordwestfalen.de/tibrosBB/azubiHeftEditForm.jsp
// https://www.bildung-ihk-nordwestfalen.de/tibrosBB/logout.jsp

const (
	BaseRoute       = "/tibrosBB"
	StartRoute      = BaseRoute + "/BB_auszubildende.jsp"
	LoginRoute      = BaseRoute + "/azubiHome.jsp"
	ReportRoute     = BaseRoute + "/azubiHeft.jsp"
	EditReportRoute = BaseRoute + "/azubiHeftEditForm.jsp"
	AddReportRoute  = BaseRoute + "/azubiHeftAdd.jsp"
	UserAgent       = "Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/106.0.0.0 Mobile Safari/537.36"
)

var (
	ErrNoSessionID = errors.New("no session id received from server")
	ErrNoValidHTML = errors.New("got no valid html from server")
	ErrNoFormToken = errors.New("no token in form")
	ErrNoStartDate = errors.New("no start date")
	ErrNoEndDate   = errors.New("no end date")
)

type IIHKService interface {
	Login(ctx context.Context, server string, username string, password string) (string, error)
	Logout(ctx context.Context, server string, tokens string) error
	CreateNewReport(ctx context.Context, server string, instructorFallbackEmail string, fallbackDepartment string, tokens string) (*Report, error)
	CancelReport(ctx context.Context, server string, report *Report, tokens string) error
	SaveReport(ctx context.Context, server string, report *Report, tokens string) error
}

type service struct {
}

func (s *service) Login(ctx context.Context, server string, username string, password string) (string, error) {
	initialId, err := s.generateInitialJSessionID(ctx, server)
	if err != nil {
		return "", err
	}

	data := &url.Values{
		"login":    []string{username},
		"pass":     []string{password},
		"anmelden": []string{""},
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, server+LoginRoute, strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Fake some data
	req = s.authenticateRequest(req, server, initialId)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	if resp.StatusCode != http.StatusOK {
		_, _ = io.Copy(os.Stdout, resp.Body)
		return "", s.parseError(resp)
	}
	for _, cookie := range resp.Cookies() {
		if cookie.Name == "JSESSIONID" {
			return fmt.Sprintf("JSESSIONID=%s; %s", cookie.Value, initialId), nil
		}
	}
	return "", ErrNoSessionID
}

func (s *service) Logout(ctx context.Context, server string, token string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/tibrosBB/logout.jsp", server), nil)
	if err != nil {
		return err
	}

	req = s.authenticateRequest(req, server, token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return s.parseError(resp)
	}
	return resp.Body.Close()
}

func (s *service) generateInitialJSessionID(ctx context.Context, server string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, server+StartRoute, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", UserAgent)
	req.Header.Set("Referer", server+StartRoute)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	_ = resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", s.parseError(resp)
	}
	rawCookie := ""
	for _, cookie := range resp.Cookies() {
		if cookie.Name == "JSESSIONID" {
			rawCookie = "JSESSIONID=" + cookie.Value
		}
	}
	//
	//req, err = http.NewRequestWithContext(ctx, http.MethodGet, "https://www.bildung-ihk-nordwestfalen.de/favicon.ico", nil)
	//resp, err = http.DefaultClient.Do(req)
	//if err != nil {
	//	return "", err
	//}
	//_ = resp.Body.Close()
	//if resp.StatusCode != http.StatusOK {
	//	return "", s.parseError(resp)
	//}
	//
	//for _, cookie := range resp.Cookies() {
	//	if cookie.Name == "JSESSIONID" {
	//		rawCookie += "; JSESSIONID=" + cookie.Value
	//	}
	//}

	if rawCookie == "" {
		return "", ErrNoSessionID
	}

	return rawCookie, nil
}

func (s *service) SaveReport(ctx context.Context, server string, report *Report, token string) error {
	buffer, boundary, err := report.ToMultipartFormData("save")
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, server+AddReportRoute, buffer)
	if err != nil {
		return err
	}
	req = s.authenticateRequest(req, server, token)
	req.Header.Set("Content-Type", "multipart/form-data; boundary=--"+boundary)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	return s.parseError(resp)
}

func (s *service) CancelReport(ctx context.Context, server string, report *Report, token string) error {
	buffer, boundary, err := report.ToMultipartFormData("cancel")
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, server+AddReportRoute, buffer)
	if err != nil {
		return err
	}
	req = s.authenticateRequest(req, server, token)
	req.Header.Set("Content-Type", "multipart/form-data; boundary=--"+boundary)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return s.parseError(resp)
	}

	return nil
}

func (s *service) dummyRequest(ctx context.Context, server string, token string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, server+ReportRoute, nil)
	if err != nil {
		return err
	}
	req = s.authenticateRequest(req, server, token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	_ = resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return s.parseError(resp)
	}
	return nil
}

func (s *service) CreateNewReport(ctx context.Context, server string, instructorFallbackEmail string, fallbackDepartment string, token string) (*Report, error) {
	if err := s.dummyRequest(ctx, server, token); err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, server+EditReportRoute, strings.NewReader("neu="))
	if err != nil {
		return nil, err
	}
	req = s.authenticateRequest(req, server, token)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, s.parseError(resp)
	}

	buffer := &bytes.Buffer{}
	_, _ = io.Copy(buffer, resp.Body)

	doc, err := goquery.NewDocumentFromReader(buffer)
	if err != nil {
		return nil, err
	}

	form := doc.Find("form").First()
	if form == nil {
		_, _ = io.Copy(os.Stdout, buffer)
		return nil, ErrNoValidHTML
	}

	tokenInput := form.Find("input[name=token]").First()
	if tokenInput == nil {
		_, _ = io.Copy(os.Stdout, buffer)
		return nil, ErrNoValidHTML
	}

	token, ok := tokenInput.Attr("value")
	if !ok {
		_, _ = io.Copy(os.Stdout, buffer)
		return nil, ErrNoFormToken
	}

	emailInput := form.Find("input[name=ausbMail]").First()
	if emailInput == nil {
		return nil, ErrNoValidHTML
	}

	email, ok := emailInput.Attr("value")
	if !ok {
		email = instructorFallbackEmail
	}

	contentActivityInput := form.Find("input[name=ausbinhalt1]").First()
	if contentActivityInput == nil {
		return nil, ErrNoValidHTML
	}

	contentActivity, ok := contentActivityInput.Attr("value")
	if !ok {
		contentActivity = ""
	}

	contentSubjectsInput := form.Find("input[name=ausbinhalt2]").First()
	if contentSubjectsInput == nil {
		return nil, ErrNoValidHTML
	}

	contentSubjects, ok := contentSubjectsInput.Attr("value")
	if !ok {
		contentSubjects = ""
	}

	contentTrainingInput := form.Find("input[name=ausbinhalt3]").First()
	if contentTrainingInput == nil {
		return nil, ErrNoValidHTML
	}

	contentTraining, ok := contentTrainingInput.Attr("value")
	if !ok {
		contentTraining = ""
	}

	startDateInput := form.Find("input[name=edtvon]").First()
	if startDateInput == nil {
		return nil, ErrNoValidHTML
	}

	startDateText, ok := startDateInput.Attr("value")
	if !ok {
		return nil, ErrNoStartDate
	}

	startDate, err := time.Parse("02.01.2006", startDateText)
	if err != nil {
		return nil, err
	}

	endDateInput := form.Find("input[name=edtbis]").First()
	if endDateInput == nil {
		return nil, ErrNoValidHTML
	}

	endDateText, ok := endDateInput.Attr("value")
	if !ok {
		return nil, ErrNoEndDate
	}

	endDate, err := time.Parse("02.01.2006", endDateText)
	if err != nil {
		return nil, err
	}

	return &Report{
		Token:           token,
		InstructorEmail: email,
		Department:      fallbackDepartment,
		StartDate:       startDate,
		EndDate:         endDate,
		ContentActivity: contentActivity,
		ContentSubjects: contentSubjects,
		ContentTraining: contentTraining,
	}, nil
}

func (s *service) parseError(resp *http.Response) error {
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return fmt.Errorf("invalid status code: %d", resp.StatusCode)
	}
	errorDiv := doc.Find(".error").First()
	if errorDiv == nil {
		if resp.StatusCode == http.StatusOK {
			return nil
		}
		return fmt.Errorf("invalid status code: %d", resp.StatusCode)
	}
	return errors.New(errorDiv.Text())
}

func (s *service) authenticateRequest(req *http.Request, server string, cookie string) *http.Request {
	//req.Header.Set("Referer", server+LoginRoute)
	req.Header.Set("Origin", server)
	req.Header.Set("User-Agent", UserAgent)
	req.Header.Set("Cookie", cookie)
	req.Header.Set("Accept", "text/html")
	return req
}

func New() IIHKService {
	return &service{}
}
