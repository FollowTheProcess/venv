package cli

import (
	"fmt"
)

// cwdHasFile returns whether or not the cwd has a file in it
// or an error if this could not be determined
func (a *App) cwdHasFile(path string) (bool, error) {
	exists, err := a.FS.Exists(path)
	if err != nil {
		return false, fmt.Errorf("could not determine if %s exists: %w", path, err)
	}

	return exists, nil
}

// cwdHasDir returns whether or not the cwd has a directory in it
// or an error if this could not be determined
func (a *App) cwdHasDir(path string) (bool, error) {
	exists, err := a.FS.DirExists(path)
	if err != nil {
		return false, fmt.Errorf("could not determine if %s exists: %w", path, err)
	}

	return exists, nil
}
