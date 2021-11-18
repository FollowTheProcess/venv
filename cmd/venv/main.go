package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/FollowTheProcess/msg"
	"github.com/FollowTheProcess/venv/cli"
)

var (
	help    bool // The --help flag
	version bool // The --version flag
)

func main() {
	// Set up flags
	flag.BoolVar(&help, "help", false, "--help")
	flag.BoolVar(&version, "version", false, "--version")

	app := cli.New(os.Stdout, os.Stderr)

	flag.Usage = app.Help

	flag.Parse()

	// venv accepts no arguments (for now)
	if flag.NArg() != 0 {
		prefix := msg.Sfail("Error:")
		message := fmt.Sprintf("venv accepts no command line arguments, got: %v", flag.Args())
		msg.Textf("%s %s", prefix, message)
		os.Exit(1)
	}

	switch {
	case help:
		app.Help()
	case version:
		app.Version()
	default:
		// Run the actual program
		if err := app.Run(); err != nil {
			prefix := msg.Sfail("Error:")
			msg.Textf("%s %s", prefix, err)
			os.Exit(1)
		}
	}
}
