// Package poetry implements wrapper functions around poetry commands
package poetry

import (
	"fmt"
	"io"
	"os/exec"

	"github.com/pelletier/go-toml/v2"
	"github.com/spf13/afero"
)

var poetryCommand = exec.Command

// If the build-backend says this, it's a valid poetry project
const poetryMarker = "poetry.core.masonry.api"

type pyProjectTOML struct {
	BuildSystem struct {
		Requires     []string `toml:"requires"`
		BuildBackend string   `toml:"build-backend"`
	} `toml:"build-system"`
}

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

// IsPoetryFile reads the contents of the toml file given by 'path' and
// determines if this is a valid poetry pyproject.toml file
func IsPoetryFile(af afero.Afero, path string) (bool, error) {
	var pyToml pyProjectTOML

	data, err := af.ReadFile(path)
	if err != nil {
		return false, fmt.Errorf("could not read %s: %w", path, err)
	}

	if err := toml.Unmarshal(data, &pyToml); err != nil {
		return false, fmt.Errorf("could not unmarshall toml data: %w", err)
	}

	return pyToml.BuildSystem.BuildBackend == poetryMarker, nil
}
