// Package cli implements the CLI functionality
// main defers execution to the exported methods in this package
package cli

import (
	"fmt"
	"io"
	"os"

	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

var (
	version = "dev" // venv's version, set at compile time with ldflags
	commit  = ""    // venv's commit hash, set at compile time with ldflags
)

const (
	helpText = `
venv

CLI to take the pain out of python virtual environments 🛠

venv aims to eliminate all the pain and hastle from creating, installing,
and managing python virtual environments

It does this by trying to work out what it is you want it to do based on
context in the surrounding directory/project.

Usage:

  venv [flags]

Examples:

# Let venv work everything out
$ venv

Flags:
  -h, --help      Help for py
  -v, --version   Show venv's version info

Environment Variables:
  VENV_DEBUG   If set to anything will print debug information to stderr
`
)

// App represents the venv CLI program
type App struct {
	Stdout io.Writer      // Where to write "normal" CLI output
	Stderr io.Writer      // Debug logs and errors will write here
	Logger *logrus.Logger // The debug logger
	FS     afero.Afero
}

func New(stdout, stderr io.Writer, fs afero.Fs) *App {
	log := logrus.New()

	// If the VENV_DEBUG environment variable is set to anything
	// set logging level accoringly
	if debug := os.Getenv("VENV_DEBUG"); debug != "" {
		log.Level = logrus.DebugLevel
	}

	log.Formatter = &logrus.TextFormatter{DisableLevelTruncation: true, DisableTimestamp: true}
	log.Out = stderr

	// Create the afero type and give it the filesystem
	af := afero.Afero{Fs: fs}

	return &App{Stdout: stdout, Stderr: stderr, Logger: log, FS: af}
}

// Help prints venv's help text
func (a *App) Help() {
	fmt.Fprintln(a.Stdout, helpText)
}

// Version shows venv's version information
func (a *App) Version() {
	ver := color.CyanString("venv version")
	sha := color.CyanString("commit")

	fmt.Fprintf(a.Stdout, "%s: %s\n", ver, version)
	fmt.Fprintf(a.Stdout, "%s: %s\n", sha, commit)
}

func (a *App) Run() error {
	fmt.Fprintln(a.Stdout, "App.Run was called")

	exists, err := a.cwdHasFile("justfile")
	if err != nil {
		return err
	}

	fmt.Printf("justfile exists: %v\n", exists)

	return nil
}
