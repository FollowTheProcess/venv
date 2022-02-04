// Package cli implements the CLI functionality
// main defers execution to the exported methods in this package
package cli

import (
	"fmt"
	"io"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/FollowTheProcess/msg"
	"github.com/FollowTheProcess/venv/pkg/flit"
	"github.com/FollowTheProcess/venv/pkg/poetry"
	"github.com/FollowTheProcess/venv/pkg/python"
	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

var (
	version = "dev" // venv's version, set at compile time with ldflags
	commit  = ""    // venv's commit hash, set at compile time with ldflags
)

const (
	debugEnv        = "VENV_DEBUG"
	reqTxt          = "requirements.txt"
	reqDev          = "requirements_dev.txt"
	pyProjectTOML   = "pyproject.toml"
	venvDir         = "venv"
	dotVenvDir      = ".venv"
	setupCFG        = "setup.cfg"
	setupPy         = "setup.py"
	createNewOption = "Create a new Environment"
	abortOption     = "Abort"
	helpText        = `
CLI to take the pain out of python virtual environments 🛠

venv aims to eliminate all the pain and hastle from creating and managing
python virtual environments as well as installing project dependencies.

It does this by trying to work out what it is you want it to do based on
context in the surrounding directory/project.

For the full logical flow followed by venv, see the README
at https://github.com/FollowTheProcess/venv

If it gets to the end of this flow without figuring out what to do
it will ask you!

You may also bypass the interactive prompt entirely by passing either
the '-c/--create' or the '-a/--abort' flag which are equivalent to answering
their interactive counterparts but bypass the interactive prompt so that
venv can be used in scripts without interruption.

Usage:

  venv [flags]

Examples:

# Let venv work everything out
$ venv

Flags:
  -h, --help      Help for venv
  -v, --version   Show venv's version info
  -c, --create    Bypass interactive prompt, telling it to create a new virtual environment
  -a, --abort     Bypass interactive prompt, telling it to abort and exit

Environment Variables:
  VENV_DEBUG   If set to anything will print debug information to stderr`
)

// App represents the venv CLI program
type App struct {
	stdout  io.Writer      // Where to write "normal" CLI output
	stderr  io.Writer      // Debug logs and errors will write here
	logger  *logrus.Logger // The debug logger
	printer *msg.Printer   // The printer in charge of informing the user (will talk to 'stdout')
	fs      afero.Afero    // A filesystem, so we can mock out during tests
}

// New creates and returns a new App configured with the filesystem, logger
// and printers
func New(stdout, stderr io.Writer, fs afero.Fs, printer *msg.Printer) *App {
	log := logrus.New()

	// If the VENV_DEBUG environment variable is set to anything
	// set logging level accoringly
	if debug := os.Getenv(debugEnv); debug != "" {
		log.Level = logrus.DebugLevel
	}

	log.Formatter = &logrus.TextFormatter{
		DisableLevelTruncation: true,
		DisableTimestamp:       true,
		ForceQuote:             true,
	}

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
func (a *App) Run(create, abort bool) error { // nolint: gocyclo
	// gocyclo moans because of too many switches but realistically this is the easiest
	// way of handling it and ensuring only one logical branch is executed
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("could not get cwd: %w", err)
	}

	if create && abort {
		// These flags are mutually exclusive
		return fmt.Errorf("--create and --abort are mutually exclusive")
	}

	switch {
	case a.cwdHasDir(dotVenvDir):
		// .venv found in cwd
		a.logger.WithField("venv directory", dotVenvDir).Debugln("virtual environment directory found")
		a.printer.Infof("There is already a virtual environment in this directory: %q", dotVenvDir)

	case a.cwdHasDir(venvDir):
		// venv found in cwd
		a.logger.WithField("venv directory", venvDir).Debugln("virtual environment directory found")
		a.printer.Infof("There is already a virtual environment in this directory: %q", venvDir)

	case a.cwdHasFile(reqDev):
		// requirements_dev.txt found in cwd
		a.logger.WithField("requirements file", reqDev).Debugln("requirements file found")
		a.printer.Infof("Found %q. Creating virtual environment and installing requirements", reqDev)
		if err := python.CreateVenv(cwd, a.stdout, a.stderr); err != nil {
			return fmt.Errorf("%w", err)
		}
		if err := python.UpdateSeeds(cwd, a.stdout, a.stderr); err != nil {
			return fmt.Errorf("%w", err)
		}
		if err := python.InstallRequirements(cwd, a.stdout, a.stderr, reqDev); err != nil {
			return fmt.Errorf("%w", err)
		}

	case a.cwdHasFile(reqTxt):
		// requirements.txt found in cwd
		a.logger.WithField("requirements file", reqTxt).Debugln("requirements file found")
		a.printer.Infof("Found %q. Creating virtual environment and installing requirements", reqTxt)
		if err := python.CreateVenv(cwd, a.stdout, a.stderr); err != nil {
			return fmt.Errorf("%w", err)
		}
		if err := python.UpdateSeeds(cwd, a.stdout, a.stderr); err != nil {
			return fmt.Errorf("%w", err)
		}
		if err := python.InstallRequirements(cwd, a.stdout, a.stderr, reqTxt); err != nil {
			return fmt.Errorf("%w", err)
		}

	case a.cwdHasFile(pyProjectTOML):
		// pyproject.toml found in cwd
		a.logger.Debugln(fmt.Sprintf("%s found", pyProjectTOML))
		switch {
		case a.cwdHasFile(setupCFG):
			// pyproject.toml + setup.cfg found in cwd
			a.logger.WithField("setuptools file", setupCFG).Debugln("found setuptools file")
			a.printer.Infof("Found %q with %q. Creating virtual environment and installing dependencies (setuptools)", pyProjectTOML, setupCFG)
			// Make a venv and install -e .[dev]
			// If the project does not define [dev] extras, pip will automatically fall back to -e . for us
			if err := python.CreateVenv(cwd, a.stdout, a.stderr); err != nil {
				return fmt.Errorf("%w", err)
			}
			if err := python.UpdateSeeds(cwd, a.stdout, a.stderr); err != nil {
				return fmt.Errorf("%w", err)
			}
			if err := python.Install(cwd, a.stdout, a.stderr, []string{"-e", ".[dev]"}); err != nil {
				return fmt.Errorf("%w", err)
			}

		case a.cwdHasFile(setupPy):
			// pyproject.toml + setup.py found in cwd
			a.logger.WithField("setuptools file", setupPy).Debugln("found setuptools file")
			a.printer.Infof("Found %q with %q. Creating virtual environment and installing dependencies (setuptools)", pyProjectTOML, setupPy)
			// Since parsing a python file to determine if it has a .[dev] might be tricky
			// just do a normal -e .
			if err := python.CreateVenv(cwd, a.stdout, a.stderr); err != nil {
				return fmt.Errorf("%w", err)
			}
			if err := python.UpdateSeeds(cwd, a.stdout, a.stderr); err != nil {
				return fmt.Errorf("%w", err)
			}
			if err := python.Install(cwd, a.stdout, a.stderr, []string{"-e", "."}); err != nil {
				return fmt.Errorf("%w", err)
			}

		default:
			// pyproject.toml found on it's own in cwd
			a.logger.Debugln("project not setuptools based")
			a.logger.Debugln("checking whether it's poetry or flit")
			// Parse pyproject.toml to determine poetry or flit and make the call
			poetryFile, err := poetry.IsPoetryFile(a.fs, pyProjectTOML)
			if err != nil {
				return fmt.Errorf("%w", err)
			}
			flitFile, err := flit.IsFlitFile(a.fs, pyProjectTOML)
			if err != nil {
				return fmt.Errorf("%w", err)
			}

			switch {
			case poetryFile:
				// pyproject.toml is a poetry spec
				a.logger.WithField("file", pyProjectTOML).Debugln("project file specifies poetry")
				a.printer.Infof("Found %q specifying poetry. Installing...", pyProjectTOML)
				if err := poetry.Install(cwd, a.stdout, a.stderr); err != nil {
					return fmt.Errorf("%w", err)
				}

			case flitFile:
				// pyproject.toml is a flit spec
				a.logger.WithField("file", pyProjectTOML).Debugln("project file specifies flit")
				a.printer.Infof("Found %q specifying flit. Installing...", pyProjectTOML)
				if err := flit.Install(cwd, a.stdout, a.stderr); err != nil {
					return fmt.Errorf("%w", err)
				}
			}
		}

	default:
		a.logger.Debugln("cannot detect environment for project")
		a.printer.Warn("Cannot auto-detect project environment")
		// User called `venv` so must want something doing
		// check create or abort flags or prompt for what to do next

		switch {
		case abort:
			// User passed --abort
			a.printer.Fail("Aborting!")
			return nil

		case create:
			// User passed --create
			a.printer.Info("Creating a new python virtual environment")
			if err := python.CreateVenv(cwd, a.stdout, a.stderr); err != nil {
				return fmt.Errorf("%w", err)
			}
			if err := python.UpdateSeeds(cwd, a.stdout, a.stderr); err != nil {
				return fmt.Errorf("%w", err)
			}

		default:
			// User didn't pass the flags so prompt for what to do next
			next := ""
			prompt := &survey.Select{
				Message: "What's next?",
				Options: []string{createNewOption, abortOption},
			}
			if err := survey.AskOne(prompt, &next); err != nil {
				return fmt.Errorf("could not generate prompt: %w", err)
			}

			switch next {
			case createNewOption:
				a.printer.Info("Creating a new python virtual environment")
				if err := python.CreateVenv(cwd, a.stdout, a.stderr); err != nil {
					return fmt.Errorf("%w", err)
				}
				if err := python.UpdateSeeds(cwd, a.stdout, a.stderr); err != nil {
					return fmt.Errorf("%w", err)
				}

			case abortOption:
				a.printer.Fail("Aborting!")
				return nil

			default:
				// This should never happen
				return fmt.Errorf("somehow entered an unrecognised option in prompt: %s", next)
			}

		}

	}

	// We'll only get here if whatever logical branch was run was successful
	// so return nil and a Done marker
	a.printer.Good("Done")
	return nil
}
