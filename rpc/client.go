package rpc

import (
	"bytes"
	"encoding/json"
	"github.com/google/uuid"
	"io"
	"net/http"
	"net/http/cookiejar"
	"time"
)

type Client struct {
	url    string
	client http.Client
}

func Dial(url string) (*Client, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}
	return &Client{
		url: url,
		client: http.Client{
			Jar:     jar,
			Timeout: 5 * time.Second,
		},
	}, nil
}

func (c *Client) Call(method string, params any, result any) error {
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
	req, err := http.NewRequest(http.MethodPost, c.url, buf)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

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
