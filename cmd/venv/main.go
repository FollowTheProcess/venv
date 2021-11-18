package main

import (
	"os"

	"github.com/FollowTheProcess/msg"
	"github.com/FollowTheProcess/venv/cli"
)

func main() {
	app := cli.New(os.Stdout, os.Stderr)

	if err := app.Run(os.Args[1:]); err != nil {
		prefix := msg.Sfail("Error:")
		msg.Textf("%s %s", prefix, err)
		os.Exit(1)
	}
}
