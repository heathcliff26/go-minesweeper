package main

import (
	"flag"

	"github.com/heathcliff26/go-minesweeper/pkg/app"
)

func main() {
	flag.Parse()
	initializeLogger()

	app := app.New()
	app.Run()
}
