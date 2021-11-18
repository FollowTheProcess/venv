// Package cmd implements the venv CLI
package cmd

import (
	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"
)

var (
	version = "dev" // venv version, set at compile time by ldflags
	commit  = ""    // venv version's commit hash, set at compile time by ldflags
)

func BuildRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:           "venv <command> [flags]",
		Args:          cobra.NoArgs,
		SilenceUsage:  true,
		SilenceErrors: true,
		Short:         "CLI to take the pain out of python virtual environments.",
		Long: heredoc.Doc(`
		
		Longer description of your CLI.
		`),
		Example: heredoc.Doc(`

		$ venv hello

		$ venv version

		$ venv --help
		`),
	}

	// Attach child commands
	rootCmd.AddCommand(
		buildVersionCmd(),
		buildHelloCommand(),
	)

	return rootCmd
}
