package rpc

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"io"
	"net/http"
	"time"
)

type Client struct {
	url     string
	session string
	client  http.Client
}

func New(url string, session string) *Client {
	return &Client{
		url:     url,
		session: session,
		client: http.Client{
			Timeout: 20 * time.Second,
		},
	}
}

func (c *Client) Call(ctx context.Context, method string, params any, result any) error {
	buf := &bytes.Buffer{}

	err := json.NewEncoder(buf).Encode(&Request{
		Id:      uuid.NewString(),
		Method:  method,
		Params:  params,
		JsonRpc: "2.0",
	})
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.url, buf)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	if c.session != "" {
		req.AddCookie(&http.Cookie{
			Name:  "JSESSIONID",
			Value: c.session,
		})
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	response := &Response{Result: result}
	err = json.NewDecoder(resp.Body).Decode(response)
	if err != nil {
		return err
	}
	if response.Error != nil {
		return response.Error
	}
	return nil
}

func (c *Client) Client() *http.Client {
	return &c.client
}
