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

	switch {
	case a.cwdHasDir(".venv"):
		a.Logger.WithField("venv directory", ".venv").Debugln("virtual environment directory found")
		// Nothing to do, just say there's already a venv

	case a.cwdHasDir("venv"):
		a.Logger.WithField("venv directory", "venv").Debugln("virtual environment directory found")
		// Nothing to do, just say there's already a venv

	case a.cwdHasFile("requirements_dev.txt"):
		a.Logger.WithField("requirements file", "requirements_dev.txt").Debugln("requirements file found")
		// Make a venv and install

	case a.cwdHasFile("requirements.txt"):
		a.Logger.WithField("requirements file", "requirements.txt").Debugln("requirements file found")
		// Make a venv and install

	case a.cwdHasFile("pyproject.toml"):
		a.Logger.Debugln("pyproject.toml found")
		switch {
		case a.cwdHasFile("setup.cfg"):
			a.Logger.WithField("setuptools file", "setup.cfg").Debugln("found setuptools file")
			// Make a venv and install -e .[dev]
			// Maybe parse the file to check if has [dev], if yes use that, if not just -e .

		case a.cwdHasFile("setup.py"):
			a.Logger.WithField("setuptools file", "setup.py").Debugln("found setuptools file")
			// Same as above branch except parsing a [dev] equivalent might be hard
			// just do -e .

		default:
			a.Logger.Debugln("project not setuptools based")
			a.Logger.Debugln("checking whether it's poetry or flit")
			// Parse pyproject.toml to determine poetry or flit and make the call
			// should be an easy toml parse, look for [tool.poetry] or [tool.flit]
		}

	case a.cwdHasFile("environment.yml"):
		a.Logger.Debugln("environment.yml found")
		// Get environment name from environment.yml
		// check output of conda env list to see if it exists on system
		// if it does, just say so and exit
		// if not, create it

	default:
		a.Logger.Debugln("cannot detect environment for project")
		// User called `venv` so must want something doing
		// Prompt for tool to create new environment with and call it

	}

	// We'll only get here if whatever branch was run was successful
	// so return nil
	return nil
}
