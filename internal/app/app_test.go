package app

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/callmesangio/lovli-cli/internal/lovli"
)

var usage = strings.Join(
	[]string{
		"Usage of lovli:",
		"  -h\tPrint help",
		"  -s string",
		"    \tShorten the URL passed as argument",
		"  -v\tPrint version",
	},
	"\n",
)

type testShortener func() (*lovli.Redirection, error)

func (f testShortener) Shorten(longURL *string) (*lovli.Redirection, error) {
	return f()
}

func nullShortener() testShortener {
	return func() (*lovli.Redirection, error) {
		return &lovli.Redirection{}, nil
	}
}

func TestUnexpectedFlag(t *testing.T) {
	stdout, stderr := &strings.Builder{}, &strings.Builder{}
	app := &App{stdout: stdout, stderr: stderr, client: nullShortener()}

	status := app.Run([]string{"-x"})

	if status != 2 {
		t.Errorf("Unexpected exit status: %d", status)
	}
	if stdout.String() != "" {
		t.Error("Unexpected stdout")
	}
	if stderr.String() != fmt.Sprintf(
		"flag provided but not defined: -x\n%s\n", usage,
	) {
		t.Error("Unexpected stderr")
	}
}

func TestNoFlag(t *testing.T) {
	stdout, stderr := &strings.Builder{}, &strings.Builder{}
	app := &App{stdout: stdout, stderr: stderr, client: nullShortener()}

	status := app.Run([]string{})

	if status != 2 {
		t.Errorf("Unexpected exit status: %d", status)
	}
	if stdout.String() != "" {
		t.Error("Unexpected stdout")
	}
	if stderr.String() != fmt.Sprintf("%s\n", usage) {
		t.Error("Unexpected stderr")
	}
}

func TestVersion(t *testing.T) {
	stdout, stderr := &strings.Builder{}, &strings.Builder{}
	app := &App{stdout: stdout, stderr: stderr, client: nullShortener()}

	status := app.Run([]string{"-v"})

	if status != 0 {
		t.Errorf("Unexpected exit status: %d", status)
	}
	if stdout.String() != fmt.Sprintf("%s\n", version) {
		t.Error("Unexpected stdout")
	}
	if stderr.String() != "" {
		t.Error("Unexpected stderr")
	}
}

func TestHelp(t *testing.T) {
	stdout, stderr := &strings.Builder{}, &strings.Builder{}
	app := &App{stdout: stdout, stderr: stderr, client: nullShortener()}

	status := app.Run([]string{"-h"})

	if status != 0 {
		t.Errorf("Unexpected exit status: %d", status)
	}
	if stdout.String() != fmt.Sprintf("%s\n", usage) {
		t.Error("Unexpected stdout")
	}
	if stderr.String() != "" {
		t.Error("Unexpected stderr")
	}
}

func TestShortenMissingArgument(t *testing.T) {
	stdout, stderr := &strings.Builder{}, &strings.Builder{}
	app := &App{stdout: stdout, stderr: stderr, client: nullShortener()}

	status := app.Run([]string{"-s"})

	if status != 2 {
		t.Errorf("Unexpected exit status: %d", status)
	}
	if stdout.String() != "" {
		t.Error("Unexpected stdout")
	}
	if stderr.String() != fmt.Sprintf(
		"flag needs an argument: -s\n%s\n", usage,
	) {
		t.Error("Unexpected stderr")
	}
}

func TestShortenInvalidArgument(t *testing.T) {
	stdout, stderr := &strings.Builder{}, &strings.Builder{}
	app := &App{stdout: stdout, stderr: stderr, client: nullShortener()}

	for _, args := range [][]string{
		{"-s", ""},
		{"-s", "\t "},
	} {
		status := app.Run(args)

		if status != 1 {
			t.Errorf("Unexpected exit status: %d", status)
		}
		if stdout.String() != "" {
			t.Error("Unexpected stdout")
		}
		if stderr.String() != "invalid URL\n" {
			t.Error("Unexpected stderr")
		}

		stdout.Reset()
		stderr.Reset()
	}
}

func TestShortenFailure(t *testing.T) {
	stdout, stderr := &strings.Builder{}, &strings.Builder{}
	var client testShortener = func() (*lovli.Redirection, error) {
		return nil, errors.New("an error occurred")
	}
	app := &App{stdout: stdout, stderr: stderr, client: client}

	status := app.Run([]string{"-s", "https://long.url.example.com"})

	if status != 1 {
		t.Errorf("Unexpected exit status: %d", status)
	}
	if stdout.String() != "" {
		t.Error("Unexpected stdout")
	}
	if stderr.String() != "an error occurred\n" {
		t.Error("Unexpected stderr")
	}
}

func TestShortenSuccess(t *testing.T) {
	stdout, stderr := &strings.Builder{}, &strings.Builder{}
	var client testShortener = func() (*lovli.Redirection, error) {
		return &lovli.Redirection{ShortURL: "https://example.com/abcd"}, nil
	}
	app := &App{stdout: stdout, stderr: stderr, client: client}

	status := app.Run([]string{"-s", "https://long.url.example.com"})

	if status != 0 {
		t.Errorf("Unexpected exit status: %d", status)
	}
	if stdout.String() != "https://example.com/abcd\n" {
		t.Error("Unexpected stdout")
	}
	if stderr.String() != "" {
		t.Error("Unexpected stderr")
	}
}

func TestFlagPriority(t *testing.T) {
	stdout, stderr := &strings.Builder{}, &strings.Builder{}
	var client testShortener = func() (*lovli.Redirection, error) {
		return &lovli.Redirection{ShortURL: "https://example.com/abcd"}, nil
	}
	app := &App{stdout: stdout, stderr: stderr, client: client}

	app.Run([]string{"-s", "https://long.url.example.com", "-v", "-h"})
	if stdout.String() != fmt.Sprintf("%s\n", version) {
		t.Error("Unexpected stdout")
	}
	stdout.Reset()

	app.Run([]string{"-s", "https://long.url.example.com", "-h"})
	if stdout.String() != fmt.Sprintf("%s\n", usage) {
		t.Error("Unexpected stdout")
	}
	stdout.Reset()

	app.Run([]string{"-s", "-v", "-h"})
	if stdout.String() != fmt.Sprintf("%s\n", usage) {
		t.Error("Unexpected stdout")
	}
}
