package lovli

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Redirection struct {
	ShortURL string `json:"short_url"`
}

type Client struct {
	client   *http.Client
	endpoint string
}

func NewClient() *Client {
	return &Client{
		client:   &http.Client{Timeout: 10 * time.Second},
		endpoint: "https://lovli.fyi/redirections",
	}
}

func (c *Client) Shorten(longURL *string) (*Redirection, error) {
	req := c.request(longURL)
	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode == http.StatusOK {
		return redirection(res.Body)
	}
	return nil, errorBy(res.StatusCode)
}

func (c *Client) request(longURL *string) *http.Request {
	body := fmt.Appendf([]byte{}, `{"location": "%s"}`, *longURL)
	req, _ := http.NewRequest("POST", c.endpoint, bytes.NewBuffer(body))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	return req
}

func redirection(body io.ReadCloser) (*Redirection, error) {
	red := &Redirection{}
	err := json.NewDecoder(body).Decode(red)
	if err != nil {
		return nil, err
	}
	return red, nil
}

func errorBy(statusCode int) error {
	switch statusCode {
	case http.StatusBadRequest:
		return errors.New("Invalid URL")
	case http.StatusTooManyRequests:
		return errors.New("Try again later")
	default:
		return fmt.Errorf("Unexpected error (%d)", statusCode)
	}
}
