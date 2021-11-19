package cli

import (
	"bytes"
	"io"
	"testing"

	"github.com/FollowTheProcess/msg"
	"github.com/spf13/afero"
)

func newTestPrinter(stdout io.Writer) *msg.Printer {
	return &msg.Printer{Out: stdout}
}

func TestApp_cwdHasFile(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	app := New(stdout, stderr, afero.NewMemMapFs(), newTestPrinter(stdout))

	t.Run("returns true if exists", func(t *testing.T) {
		// Make a file
		file, err := app.fs.Create("a_file")
		if err != nil {
			t.Fatalf("could not create test file: %v", err)
		}
		file.Close()

		got := app.cwdHasFile(file.Name())

		if got != true {
			t.Errorf("cwdHasFile said file %s, does not exist when it does", file.Name())
		}
	})

	t.Run("returns false if doesn't exist", func(t *testing.T) {
		got := app.cwdHasFile("im_not_here")

		if got != false {
			t.Errorf("cwdHasFile said file %s, exists when it doesn't", "im_not_here")
		}
	})
}

func TestApp_cwdHasDir(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	app := New(stdout, stderr, afero.NewMemMapFs(), newTestPrinter(stdout))

	t.Run("returns true if exists", func(t *testing.T) {
		// Make a dir
		err := app.fs.Mkdir("testdir", 0o755)
		if err != nil {
			t.Fatalf("could not create test dir: %v", err)
		}

		got := app.cwdHasDir("testdir")

		if got != true {
			t.Errorf("cwdHasFile said dir %s, does not exist when it does", "testdir")
		}
	})

	t.Run("returns false if doesn't exist", func(t *testing.T) {
		got := app.cwdHasFile("im_not_here")

		if got != false {
			t.Errorf("cwdHasFile said dir %s, exists when it doesn't", "im_not_here")
		}
	})
}
