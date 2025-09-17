package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

const endpoint string = "https://lovli.fyi/redirections"

type Redirection struct {
	ShortUrl string `json:"short_url"`
}

func main() {
	parseCli()

	url, err := url()
	if err != nil {
		fail(err)
	}

	redirection, err := post(url)
	if err != nil {
		fail(err)
	}

	fmt.Println(redirection.ShortUrl)
}

func parseCli() {
	flag.Usage = func() {
		fmt.Println(*usage())
	}
	flag.Parse()
}

func usage() *string {
	usage := fmt.Sprintf("Usage: %s <URL>", os.Args[0])
	return &usage
}

func url() (*string, error) {
	url := strings.TrimSpace(flag.Arg(0))
	if url == "" {
		return nil, errors.New(*usage())
	}
	return &url, nil
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

func post(url *string) (*Redirection, error) {
	client := http.Client{Timeout: 10 * time.Second}
	req := newRequest(url)
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode == http.StatusOK {
		return jsonDecode(res.Body)
	}
	return nil, newPostError(res.StatusCode)
}

func newRequest(url *string) *http.Request {
	body := fmt.Appendf([]byte{}, `{"location": "%s"}`, *url)
	req, _ := http.NewRequest("POST", endpoint, bytes.NewBuffer(body))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	return req
}

func jsonDecode(body io.ReadCloser) (*Redirection, error) {
	redirection := &Redirection{}
	err := json.NewDecoder(body).Decode(redirection)
	if err != nil {
		return nil, err
	}
	return redirection, nil
}

func newPostError(statusCode int) error {
	switch statusCode {
	case http.StatusBadRequest:
		return errors.New("Invalid URL")
	case http.StatusTooManyRequests:
		return errors.New("Try again later")
	default:
		return fmt.Errorf("Unexpected error (%d)", statusCode)
	}
}
