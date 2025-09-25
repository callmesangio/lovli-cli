package app

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/callmesangio/lovli-cli/internal/lovli"
)

const version string = "unreleased"

type shortener interface {
	Shorten(longURL *string) (*lovli.Redirection, error)
}

type App struct {
	stdout, stderr io.Writer
	client         shortener
	flag           *flag.FlagSet
	version        bool
	help           bool
	longURL        string
}

func (a *App) Run(args []string) int {
	var err error
	if err = a.load(args); err != nil {
		return 2
	}
	if err = a.run(); err != nil {
		fmt.Fprintln(a.stderr, err)
		return 1
	}
	return 0
}

func (a *App) load(args []string) error {
	a.flag = flag.NewFlagSet("lovli", flag.ContinueOnError)
	a.flag.SetOutput(a.stderr)
	a.flag.BoolVar(&a.version, "v", false, "Print version")
	a.flag.BoolVar(&a.help, "h", false, "Print help")
	a.flag.StringVar(&a.longURL, "s", "", "Shorten the URL passed as argument")

	if err := a.flag.Parse(args); err != nil {
		return err
	}
	if a.flag.NFlag() == 0 {
		a.flag.Usage()
		return errors.New("no flag")
	}
	return nil
}

func (a *App) run() error {
	if a.version {
		return a.printVersion()
	}
	if a.help {
		return a.printHelp()
	}
	return a.shortenURL()
}

func (a *App) printVersion() error {
	fmt.Fprintln(a.stdout, version)
	return nil
}

func (a *App) printHelp() error {
	a.flag.SetOutput(a.stdout)
	a.flag.Usage()
	return nil
}

func (a *App) shortenURL() error {
	longURL := strings.TrimSpace(a.longURL)
	if longURL == "" {
		return errors.New("invalid URL")
	}
	redirection, err := a.client.Shorten(&longURL)
	if err != nil {
		return err
	}
	fmt.Fprintln(a.stdout, redirection.ShortURL)
	return nil
}

func Run(args []string) int {
	app := &App{
		stdout: os.Stdout,
		stderr: os.Stderr,
		client: lovli.NewClient(),
	}
	return app.Run(args)
}
