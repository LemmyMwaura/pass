package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "pass",
	Short: "A simple encrypted CLI password manager",
	Long: `pass stores credentials in an encrypted vault under ~/.pass/.

Use "pass start" for the Bubble Tea TUI, or "pass generate" to create a password.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {}
