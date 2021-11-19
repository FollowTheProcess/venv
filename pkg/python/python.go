// Package python provides functions that wrap external commands to create
// standard python virtual environments and install dependencies using the
// built in python tooling (e.g. pip, venv, setuptools)
package python

import (
	"fmt"
	"io"
	"os/exec"
)

// pythonCommand is an internal reassignment of exec.Command
// used for mocking during tests
var pythonCommand = exec.Command

// newPythonCmd returns an exec.Cmd configured with the parameters passed in
func newPythonCmd(cwd string, stdout, stderr io.Writer, args []string) *exec.Cmd {
	cmd := pythonCommand("python", args...)
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
