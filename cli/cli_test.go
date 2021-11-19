package cli

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/spf13/afero"
)

func TestApp_Help(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	app := New(stdout, stderr, afero.NewMemMapFs())

	want := fmt.Sprintf("%s\n", helpText)

	// Call help
	app.Help()

	if got := stdout.String(); got != want {
		t.Errorf("got %#v, wanted %#v", got, want)
	}
}

func TestApp_Version(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	app := New(stdout, stderr, afero.NewMemMapFs())

	// Call version
	app.Version()

	if !strings.Contains(stdout.String(), "venv version") {
		t.Errorf("version string did not contain version: %s", stdout.String())
	}

	if !strings.Contains(stdout.String(), "commit") {
		t.Errorf("version string did not contain commit: %s", stdout.String())
	}
}
