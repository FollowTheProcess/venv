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
	debugEnv      = "VENV_DEBUG"
	reqTxt        = "requirements.txt"
	reqDev        = "requirements_dev.txt"
	pyProjectTOML = "pyproject.toml"
	envYAML       = "environment.yml"
	venvDir       = "venv"
	dotVenvDir    = ".venv"
	setupCFG      = "setup.cfg"
	setupPy       = "setup.py"
	helpText      = `
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
	if debug := os.Getenv(debugEnv); debug != "" {
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

// Run is the entry point to the CLI, this is what gets run when
// you call `venv` on the terminal
func (a *App) Run() error {
	fmt.Fprintln(a.Stdout, "App.Run was called")

	a.Logger.Debugln("Looking for a virtual environment")
	switch {
	case a.cwdHasDir(dotVenvDir):
		a.Logger.WithField("venv directory", dotVenvDir).Debugln("virtual environment directory found")
		// Nothing to do, just say there's already a venv

	case a.cwdHasDir(venvDir):
		a.Logger.WithField("venv directory", venvDir).Debugln("virtual environment directory found")
		// Nothing to do, just say there's already a venv
	}

	a.Logger.Debugln("Looking for a requirements file")
	switch {
	case a.cwdHasFile(reqDev):
		a.Logger.WithField("requirements file", reqDev).Debugln("requirements file found")
		// Make a venv and install

	case a.cwdHasFile(reqTxt):
		a.Logger.WithField("requirements file", reqTxt).Debugln("requirements file found")
		// Make a venv and install
	}

	a.Logger.Debugln("Looking for python package files")
	if a.cwdHasFile(pyProjectTOML) {
		a.Logger.Debugln(fmt.Sprintf("%s found", pyProjectTOML))
		switch {
		case a.cwdHasFile(setupCFG):
			a.Logger.WithField("setuptools file", setupCFG).Debugln("found setuptools file")
			// Make a venv and install -e .[dev]
			// Maybe parse the file to check if has [dev], if yes use that, if not just -e .

		case a.cwdHasFile(setupPy):
			a.Logger.WithField("setuptools file", setupPy).Debugln("found setuptools file")
			// Same as above branch except parsing a [dev] equivalent might be hard
			// just do -e .

		default:
			a.Logger.Debugln("project not setuptools based")
			a.Logger.Debugln("checking whether it's poetry or flit")
			// Parse pyproject.toml to determine poetry or flit and make the call
			// should be an easy toml parse, look for [tool.poetry] or [tool.flit]
		}
	}

	a.Logger.Debugln("Looking for a conda environment")
	if a.cwdHasFile(envYAML) {
		a.Logger.Debugln(fmt.Sprintf("%s found", envYAML))
		// Get environment name from environment.yml
		// check output of conda env list to see if it exists on system
		// if it does, just say so and exit
		// if not, create it
	}

	a.Logger.Debugln("cannot detect environment for project")
	// User called `venv` so must want something doing
	// Prompt for tool to create new environment with and call it

	// We'll only get here if whatever logical branch was run was successful
	// so return nil
	return nil
}
