package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/FollowTheProcess/msg"
	"github.com/FollowTheProcess/venv/cli"
	"github.com/spf13/afero"
)

var (
	help    bool // The --help flag
	version bool // The --version flag
	create  bool // The --create flag to bypass the interactive prompt
	abort   bool // The --abort flag to bypass the interactive prompt
)

func main() {
	// Set up flags
	flag.BoolVar(&help, "help", false, "--help")
	flag.BoolVar(&version, "version", false, "--version")
	flag.BoolVar(&create, "create", false, "--create")
	flag.BoolVar(&abort, "abort", false, "--abort")

	app := cli.New(os.Stdout, os.Stderr, afero.NewOsFs(), msg.Default())

	flag.Usage = app.Help

	flag.Parse()

	// venv accepts no arguments (for now)
	if flag.NArg() != 0 {
		message := fmt.Sprintf("venv accepts no command line arguments, got: %v", flag.Args())
		msg.Fail(message)
		os.Exit(1)
	}

	switch {
	case help:
		app.Help()
	case version:
		app.Version()
	default:
		// Run the actual program
		if err := app.Run(create, abort); err != nil {
			msg.Failf("%s", err)
			os.Exit(1)
		}
	}
}
