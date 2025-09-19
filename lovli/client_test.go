package lovli

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
)

type roundTripTransport func(*http.Request) *http.Response

func (f roundTripTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	defer req.Body.Close()
	if res := f(req); res != nil {
		return res, nil
	}
	return nil, errors.New("Transport error")
}

func newTestClient(tr roundTripTransport) *Client {
	return &Client{
		client:   &http.Client{Transport: tr},
		endpoint: "https://service.endpoint.example.com",
	}
}

func longURL() *string {
	u := "https://long.url.example.com"
	return &u
}

func TestRequestPayload(t *testing.T) {
	client := newTestClient(func(req *http.Request) *http.Response {
		data, _ := io.ReadAll(req.Body)
		if string(data) != `{"location": "https://long.url.example.com"}` {
			t.Errorf("Unexpected payload: %s", data)
		}
		return &http.Response{}
	})

	client.Shorten(longURL())
}

func TestRequestDetails(t *testing.T) {
	client := newTestClient(func(req *http.Request) *http.Response {
		if req.URL.String() != "https://service.endpoint.example.com" {
			t.Errorf("Unexpected URL: %s", req.URL)
		}

		if req.Method != "POST" {
			t.Errorf("Unexpected method: %s", req.Method)
		}

		return &http.Response{}
	})

	client.Shorten(longURL())
}

func TestRequestHeaders(t *testing.T) {
	client := newTestClient(func(req *http.Request) *http.Response {
		if len(req.Header.Values("Content-Type")) != 1 {
			t.Error(`Missing or multiple "Content-Type" header(s)`)
		}
		if len(req.Header.Values("Accept")) != 1 {
			t.Error(`Missing or multiple "Accept" header(s)`)
		}

		if req.Header.Get("Content-Type") != "application/json" {
			t.Error(`"Content-Type" header is not "application/json"`)
		}
		if req.Header.Get("Accept") != "application/json" {
			t.Error(`"Accept" header is not "application/json"`)
		}

		return &http.Response{}
	})

	client.Shorten(longURL())
}

func TestTransportError(t *testing.T) {
	var redirection *Redirection
	var err error
	client := newTestClient(func(req *http.Request) *http.Response {
		return nil
	})

	redirection, err = client.Shorten(longURL())

	if redirection != nil {
		t.Error("`redirection` expected to be `nil`")
	}
	if err == nil {
		t.Fatal("`err` expected not to be `nil`")
	}
	if err.Error() != `Post "https://service.endpoint.example.com": Transport error` {
		t.Errorf("Unexpected error message: %s", err)
	}
}

func TestRedirectionUnmarshalingFailure(t *testing.T) {
	var redirection *Redirection
	var err error
	client := newTestClient(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader("{")),
		}
	})

	redirection, err = client.Shorten(longURL())

	if redirection != nil {
		t.Error("`redirection` expected to be `nil`")
	}
	if err == nil {
		t.Fatal("`err` expected not to be `nil`")
	}
	if err.Error() != "unexpected EOF" {
		t.Errorf("Unexpected error message: %s", err)
	}
}

func TestResponseStatusCodeFailure(t *testing.T) {
	var redirection *Redirection
	var err error

	for statusCode, message := range map[int]string{
		http.StatusBadRequest:         "Invalid URL",
		http.StatusTooManyRequests:    "Try again later",
		http.StatusServiceUnavailable: "Unexpected error (503)",
	} {
		client := newTestClient(func(req *http.Request) *http.Response {
			return &http.Response{StatusCode: statusCode}
		})

		redirection, err = client.Shorten(longURL())

		if redirection != nil {
			t.Error("`redirection` expected to be `nil`")
		}
		if err == nil {
			t.Error("`err` expected not to be `nil`")
			continue
		}
		if err.Error() != message {
			t.Errorf("Unexpected error message: %s", err)
		}
	}
}

func TestSuccess(t *testing.T) {
	var redirection *Redirection
	var err error
	client := newTestClient(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`{"short_url": "https://example.com/abcd"}`)),
		}
	})

	redirection, err = client.Shorten(longURL())

	if err != nil {
		t.Error("`err` expected to be `nil`")
	}
	if redirection == nil {
		t.Fatal("`redirection` expected not to be `nil`")
	}
	if redirection.ShortURL != "https://example.com/abcd" {
		t.Errorf("Unexpected short URL: %s", redirection.ShortURL)
	}
}
