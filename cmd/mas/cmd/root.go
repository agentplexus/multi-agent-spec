// Package cmd implements the mas CLI commands.
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const version = "0.8.0"

var rootCmd = &cobra.Command{
	Use:   "mas",
	Short: "Multi-Agent Spec CLI",
	Long: `mas is the command-line interface for multi-agent-spec.

It provides tools for rendering, validating, and working with
multi-agent team reports and specifications.`,
	SilenceUsage: true,
}

// Execute runs the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Fprintf(os.Stdout, "mas version %s\n", version)
	},
}
