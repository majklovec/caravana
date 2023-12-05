package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var repoCmd = &cobra.Command{
	Use:   "repo",
	Short: "Operations with repositories",
	Long: `Operations with repositories
	- add
	- del
	- update`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			os.Exit(0)
		}
	},
	// ValidArgs: validArgs,
}

func init() {
	rootCmd.AddCommand(repoCmd)
}
