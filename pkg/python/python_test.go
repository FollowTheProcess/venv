package python

import (
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"testing"
)

// testCase is used as an env var to pass around so our test helper
// knows what condition to test for
var testCase string

// extractCmdArgs is a helper for TestHelperProcess which teases out the desired
// external command arguments from the special ones required to make go test use the
// helper process
func extractCmdArgs(args []string) []string {
	for len(args) > 0 {
		if args[0] == "--" {
			args = args[1:]
			break
		}
		args = args[1:]
	}
	return args
}

// sssertCorrectArgs compares external command arguments to verify correctness
// designed to be used inside the external command TestHelperProcess
func assertCorrectArgs(expected, args []string) {
	if !reflect.DeepEqual(args, expected) {
		fmt.Fprintf(os.Stderr, "Error: expected cmd %#v, got %#v", expected, args)
		os.Exit(1)
	}
}

// fakeExecCommand is a helper that creates a fake external command
// It does some clever magic and uses the way go test runs to insert itself
// during a test in place of an actual command
// it's used in the std lib to test exec
// see: https://npf.io/2015/06/testing-exec-command/
func fakeExecCommand(command string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestPythonHelperProcess", "--", command}
	cs = append(cs, args...)

	cmd := exec.Command(os.Args[0], cs...)
	// By passing env variables like this, we can control the behaviour of our
	// mocked command
	// For example, have it return a non-zero exit code etc.
	tc := "PYTHON_TEST_CASE=" + testCase
	cmd.Env = []string{"GO_WANT_PYTHON_HELPER_PROCESS=1", tc}
	return cmd
}

func setUp(testcase string) {
	pythonCommand = fakeExecCommand
	testCase = testcase
}

func tearDown() {
	pythonCommand = exec.Command
}

// This is the main helper process for external command tests. It first checks whether or not go test wants to use it
// by looking for the GO_WANT_HELPER_PROCESS env var (which is set by our faked external command)
// it will then separate out the arguments required to get go test to insert it from our actual
// external command arguments.
//
// It will then switch on the value of the TEST_CASE env var which each test sets individually so that it
// knows what to do
// i.e. return a 0 exit code and a success message to verify our happy path, or a non-zero exit code
// and a message to stderr to test our error handling
func TestPythonHelperProcess(t *testing.T) {
	// Tell go test to use this helper if env var is set
	if os.Getenv("GO_WANT_PYTHON_HELPER_PROCESS") != "1" {
		return
	}

	// First separate the go test args from what we actually want
	args := extractCmdArgs(os.Args)

	switch os.Getenv("PYTHON_TEST_CASE") {
	case "create_venv_success":
		expectedArgs := []string{"python", "-m", "venv", ".venv"}
		assertCorrectArgs(expectedArgs, args)

	case "create_venv_error":
		expectedArgs := []string{"python", "-m", "venv", ".venv"}
		assertCorrectArgs(expectedArgs, args)
		// Simulate failure by printing to stderr and exit 1
		fmt.Fprintf(os.Stderr, "something wrong")
		os.Exit(1)
	}
}

func TestCreateVenv(t *testing.T) {
	tests := []struct {
		testcase string
		wantErr  bool
	}{
		{
			testcase: "create_venv_success",
			wantErr:  false,
		},
		{
			testcase: "create_venv_error",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testcase, func(t *testing.T) {
			setUp(tt.testcase)
			defer tearDown()

			if err := CreateVenv(".", os.Stdout, os.Stderr); (err != nil) != tt.wantErr {
				t.Errorf("CreateVenv() error = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}
