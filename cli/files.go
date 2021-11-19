package cli

// cwdHasFile returns whether or not the cwd has a file in it
// or an error if this could not be determined
func (a *App) cwdHasFile(path string) bool {
	a.logger.WithField("file", path).Debugln("Looking for file")
	exists, err := a.fs.Exists(path)
	if err != nil {
		return false
	}

	return exists
}

// cwdHasDir returns whether or not the cwd has a directory in it
// or an error if this could not be determined
func (a *App) cwdHasDir(path string) bool {
	a.logger.WithField("directory", path).Debugln("Looking for directory")
	exists, err := a.fs.DirExists(path)
	if err != nil {
		return false
	}

	return exists
}
