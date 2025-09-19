package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/callmesangio/lovli-cli/lovli"
)

func main() {
	parseCli()

	url, err := url()
	if err != nil {
		fail(err)
	}

	redirection, err := lovli.NewClient().Shorten(url)
	if err != nil {
		fail(err)
	}

	fmt.Println(redirection.ShortURL)
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
