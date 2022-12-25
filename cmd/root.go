package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:           "kvmtoggle",
	Short:         "Daemon to handle toggling between KVM inputs",
	SilenceUsage:  true,
	SilenceErrors: true,
}

// Execute function is the entrypoint for the CLI
func Execute() error {
	return rootCmd.Execute()
}
