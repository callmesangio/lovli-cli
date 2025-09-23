package lovli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Redirection struct {
	ShortURL string `json:"short_url"`
}

type Error struct {
	Text string
	Err  error
}

func (e *Error) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s (%s)", e.Text, e.Err.Error())
	}
	return fmt.Sprintf("%s", e.Text)
}

func (e *Error) Unwrap() error { return e.Err }

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
		return nil, &Error{"transport error", err}
	}
	defer res.Body.Close()
	if res.StatusCode == http.StatusOK {
		return unmarshal(res.Body)
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

func unmarshal(body io.ReadCloser) (*Redirection, error) {
	red := &Redirection{}
	err := json.NewDecoder(body).Decode(red)
	if err != nil {
		return nil, &Error{"unmarshaling error", err}
	}
	return red, nil
}

func errorBy(statusCode int) error {
	switch statusCode {
	case http.StatusBadRequest:
		return &Error{"invalid URL", nil}
	case http.StatusTooManyRequests:
		return &Error{"try again later", nil}
	default:
		return &Error{fmt.Sprintf("unexpected error (%d)", statusCode), nil}
	}
}
