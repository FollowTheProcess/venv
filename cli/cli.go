// Package cli implements the CLI functionality
// main defers execution to the exported methods in this package
package cli

import (
	"fmt"
	"io"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/FollowTheProcess/msg"
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
	stdout  io.Writer      // Where to write "normal" CLI output
	stderr  io.Writer      // Debug logs and errors will write here
	logger  *logrus.Logger // The debug logger
	printer *msg.Printer   // The printer in charge of informing the user (will talk to 'stdout')
	fs      afero.Afero    // A filesystem, so we can mock out during tests
}

func New(stdout, stderr io.Writer, fs afero.Fs, printer *msg.Printer) *App {
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

	return &App{stdout: stdout, stderr: stderr, logger: log, fs: af, printer: printer}
}

// Help prints venv's help text
func (a *App) Help() {
	fmt.Fprintln(a.stdout, helpText)
}

// Version shows venv's version information
func (a *App) Version() {
	ver := color.CyanString("venv version")
	sha := color.CyanString("commit")

	fmt.Fprintf(a.stdout, "%s: %s\n", ver, version)
	fmt.Fprintf(a.stdout, "%s: %s\n", sha, commit)
}

// Run is the entry point to the CLI, this is what gets run when
// you call `venv` on the terminal
func (a *App) Run() error {
	fmt.Fprintln(a.stdout, "App.Run was called")

	a.logger.Debugln("Looking for a virtual environment")
	switch {
	case a.cwdHasDir(dotVenvDir):
		a.logger.WithField("venv directory", dotVenvDir).Debugln("virtual environment directory found")
		// Nothing to do, just say there's already a venv
		a.printer.Infof("There is already a virtual environment in this directory: %q", dotVenvDir)

	case a.cwdHasDir(venvDir):
		a.logger.WithField("venv directory", venvDir).Debugln("virtual environment directory found")
		// Nothing to do, just say there's already a venv
		a.printer.Infof("There is already a virtual environment in this directory: %q", venvDir)
	}

	a.logger.Debugln("Looking for a requirements file")
	switch {
	case a.cwdHasFile(reqDev):
		a.logger.WithField("requirements file", reqDev).Debugln("requirements file found")
		// Make a venv and install
		a.printer.Infof("Found %q. Creating virtual environment and installing requirements", reqDev)

	case a.cwdHasFile(reqTxt):
		a.logger.WithField("requirements file", reqTxt).Debugln("requirements file found")
		// Make a venv and install
		a.printer.Infof("Found %q. Creating virtual environment and installing requirements", reqDev)
	}

	a.logger.Debugln("Looking for python package files")
	if a.cwdHasFile(pyProjectTOML) {
		a.logger.Debugln(fmt.Sprintf("%s found", pyProjectTOML))
		switch {
		case a.cwdHasFile(setupCFG):
			a.logger.WithField("setuptools file", setupCFG).Debugln("found setuptools file")
			a.printer.Infof("Found %q with %q. Creating virtual environment and installing dependencies (setuptools)", pyProjectTOML, setupCFG)
			// Make a venv and install -e .[dev]
			// Maybe parse the file to check if has [dev], if yes use that, if not just -e .

		case a.cwdHasFile(setupPy):
			a.logger.WithField("setuptools file", setupPy).Debugln("found setuptools file")
			a.printer.Infof("Found %q with %q. Creating virtual environment and installing dependencies (setuptools)", pyProjectTOML, setupPy)
			// Same as above branch except parsing a [dev] equivalent might be hard
			// just do -e .

		default:
			a.logger.Debugln("project not setuptools based")
			a.logger.Debugln("checking whether it's poetry or flit")
			// Parse pyproject.toml to determine poetry or flit and make the call
			// should be an easy toml parse, look for [tool.poetry] or [tool.flit]
		}
	}

	a.logger.Debugln("Looking for a conda environment")
	if a.cwdHasFile(envYAML) {
		a.logger.Debugln(fmt.Sprintf("%s found", envYAML))
		// Get environment name from environment.yml
		// check output of conda env list to see if it exists on system
		// if it does, just say so and exit
		// if not, create it
	}

	a.logger.Debugln("cannot detect environment for project")
	a.printer.Warn("Cannot auto-detect project environment")
	// User called `venv` so must want something doing
	// Prompt for tool to create new environment with and call it

	next := ""
	prompt := &survey.Select{
		Message: "What next?",
		Options: []string{"Create (python)", "Create (flit)", "Create (poetry)", "Create (conda)", "Abort"},
	}
	if err := survey.AskOne(prompt, &next); err != nil {
		return fmt.Errorf("could not generate prompt: %w", err)
	}

	switch next {
	case "Create (python)":
		a.printer.Info("Creating a python virtual environment")

	case "Create (flit)":
		a.printer.Info("Creating a flit virtual environment")

	case "Create (poetry)":
		a.printer.Info("Creating a poetry virtual environment")

	case "Create (conda)":
		a.printer.Info("Creating a conda virtual environment")

	case "Abort":
		a.printer.Fail("Aborting!")

	default:
		// This should never happen
		return fmt.Errorf("somehow entered an unrecognised option in prompt: %s", next)
	}

	// We'll only get here if whatever logical branch was run was successful
	// so return nil
	return nil
}
