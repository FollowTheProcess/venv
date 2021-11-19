// Package python provides functions that wrap external commands to create
// standard python virtual environments and install dependencies using the
// built in python tooling (e.g. pip, venv, setuptools)
package python

import (
	"fmt"
	"io"
	"os/exec"
	"path/filepath"
)

// pythonCommand is an internal reassignment of exec.Command
// used for mocking during tests
var (
	pythonCommand = exec.Command
)

// newPythonCmd returns an exec.Cmd configured with the parameters passed in
// pointing to whatever python is on $PATH
func newPythonCmd(cwd string, stdout, stderr io.Writer, args []string) *exec.Cmd {
	cmd := pythonCommand("python", args...)
	cmd.Dir = cwd
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	return cmd
}

// newVenvCmd returns an exec.Cmd configured with the parameters passed in
// pointing to the virtual environment's python
func newVenvCmd(cwd string, stdout, stderr io.Writer, args []string) *exec.Cmd {
	python := filepath.Join(cwd, ".venv/bin/python")
	cmd := pythonCommand(python, args...)
	cmd.Dir = cwd
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	return cmd
}

// CreateVenv will create a standard python virtual environment named ".venv" in the cwd
//
// The wrapped external command will be hooked up directly to stdout and stderr and
// will wait for the command to complete before returning
func CreateVenv(cwd string, stdout, stderr io.Writer) error {
	cmd := newPythonCmd(cwd, stdout, stderr, []string{"-m", "venv", ".venv"})
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("could not create virtual environment: %w", err)
	}

	return nil
}

// UpdateSeeds will use the virtual environment's python to update pip, setuptools and wheel
func UpdateSeeds(cwd string, stdout, stderr io.Writer) error {
	cmd := newVenvCmd(cwd, stdout, stderr, []string{"-m", "pip", "install", "--upgrade", "pip", "setuptools", "wheel"})
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("could not update seeds: %w", err)
	}

	return nil
}

// InstallRequirements will call pip to install into a virtual environment the dependencies
// specified in a requirements file given by `file`
func InstallRequirements(cwd string, stdout, stderr io.Writer, file string) error {
	cmd := newVenvCmd(cwd, stdout, stderr, []string{"-m", "pip", "install", "-r", file})
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("could not update seeds: %w", err)
	}

	return nil
}

// Install is a wrapper around the virtual environment's pip install
// installArgs are effectively passed to "python -m pip install ..."
func Install(cwd string, stdout, stderr io.Writer, installArgs []string) error {
	args := []string{"-m", "pip", "install"}
	args = append(args, installArgs...)
	cmd := newVenvCmd(cwd, stdout, stderr, args)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("could not install %v: %w", installArgs, err)
	}

	return nil
}
