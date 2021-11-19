package cli

// cwdHasFile returns whether or not the cwd has a file in it
// or an error if this could not be determined
func (a *App) cwdHasFile(path string) bool {
	exists, err := a.FS.Exists(path)
	if err != nil {
		return false
	}

	return exists
}

// cwdHasDir returns whether or not the cwd has a directory in it
// or an error if this could not be determined
func (a *App) cwdHasDir(path string) bool {
	exists, err := a.FS.DirExists(path)
	if err != nil {
		return false
	}

	return exists
}
