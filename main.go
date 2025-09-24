package main

import (
	"os"

	"github.com/callmesangio/lovli-cli/internal/app"
)

func main() {
	os.Exit(app.Run(os.Args[1:]))
}
