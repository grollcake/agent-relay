package main

import (
	"fmt"
	"os"

	"github.com/grollcake/agent-relay/internal/relay"
)

func main() {
	app, err := relay.Discover()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if err := app.Run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
