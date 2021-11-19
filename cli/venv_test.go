package cli

import (
	"bytes"
	"testing"

	"github.com/spf13/afero"
)

func TestApp_cwdHasFile(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	app := New(stdout, stderr, afero.NewMemMapFs())

	t.Run("returns true if exists", func(t *testing.T) {
		// Make a file
		file, err := app.FS.Create("a_file")
		if err != nil {
			t.Fatalf("could not create test file: %v", err)
		}
		file.Close()

		got, err := app.cwdHasFile(file.Name())
		if err != nil {
			t.Errorf("cwdHasFile returned an error: %v", err)
		}

		if got != true {
			t.Errorf("cwdHasFile said file %s, does not exist when it does", file.Name())
		}
	})

	t.Run("returns false if doesn't exist", func(t *testing.T) {
		got, err := app.cwdHasFile("im_not_here")
		if err != nil {
			t.Errorf("cwdHasFile returned an error: %v", err)
		}

		if got != false {
			t.Errorf("cwdHasFile said file %s, exists when it doesn't", "im_not_here")
		}
	})
}

func TestApp_cwdHasDir(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	app := New(stdout, stderr, afero.NewMemMapFs())

	t.Run("returns true if exists", func(t *testing.T) {
		// Make a dir
		err := app.FS.Mkdir("testdir", 0o755)
		if err != nil {
			t.Fatalf("could not create test dir: %v", err)
		}

		got, err := app.cwdHasDir("testdir")
		if err != nil {
			t.Errorf("cwdHasFile returned an error: %v", err)
		}

		if got != true {
			t.Errorf("cwdHasFile said dir %s, does not exist when it does", "testdir")
		}
	})

	t.Run("returns false if doesn't exist", func(t *testing.T) {
		got, err := app.cwdHasFile("im_not_here")
		if err != nil {
			t.Errorf("cwdHasFile returned an error: %v", err)
		}

		if got != false {
			t.Errorf("cwdHasFile said dir %s, exists when it doesn't", "im_not_here")
		}
	})
}
