// Package poetry implements wrapper functions around poetry commands
package poetry

import (
	"fmt"
	"io"
	"os/exec"
)

var poetryCommand = exec.Command

// newPoetryCmd returns an exec.Cmd configured with the parameters passed in
func newPoetryCommand(cwd string, stdout, stderr io.Writer, args []string) *exec.Cmd {
	cmd := poetryCommand("poetry", args...)
	cmd.Dir = cwd
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	return cmd
}

// Install calls poetry install
func Install(cwd string, stdout, stderr io.Writer) error {
	cmd := newPoetryCommand(cwd, stdout, stderr, []string{"install"})
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("could not create poetry environment: %w", err)
	}

	return nil
}
