// Package flit provides wrappers around flit commands
package flit

import (
	"fmt"
	"io"
	"os/exec"

	"github.com/pelletier/go-toml/v2"
	"github.com/spf13/afero"
)

var flitCommand = exec.Command

// If the build-backend says this, it's a valid flit project
const flitMarker = "flit.buildapi"

type pyProjectTOML struct {
	BuildSystem struct {
		Requires     []string `toml:"requires"`
		BuildBackend string   `toml:"build-backend"`
	} `toml:"build-system"`
}

// newFlitCmd returns an exec.Cmd configured with the parameters passed in
func newFlitCommand(cwd string, stdout, stderr io.Writer, args []string) *exec.Cmd {
	cmd := flitCommand("flit", args...)
	cmd.Dir = cwd
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	return cmd
}

// Install calls flit install
func Install(cwd string, stdout, stderr io.Writer) error {
	cmd := newFlitCommand(cwd, stdout, stderr, []string{"install", "--deps", "develop", "--symlink", "--python", ".venv/bin/python"})
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("could not create flit environment: %w", err)
	}

	return nil
}

// IsFlitFile reads the contents of the toml file given by 'path' and
// determines if this is a valid flit pyproject.toml file
func IsFlitFile(af afero.Afero, path string) (bool, error) {
	var pyToml pyProjectTOML

	data, err := af.ReadFile(path)
	if err != nil {
		return false, fmt.Errorf("could not read %s: %w", path, err)
	}

	if err := toml.Unmarshal(data, &pyToml); err != nil {
		return false, fmt.Errorf("could not unmarshall toml data: %w", err)
	}

	return pyToml.BuildSystem.BuildBackend == flitMarker, nil
}
