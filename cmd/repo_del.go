package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// repodelCmd represents the del command
var repodelCmd = &cobra.Command{
	Use:   "del [url/dir]",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		dir := extractRepoName(args[0])
		err := os.RemoveAll(filepath.Join(templatesDir, dir))
		if err != nil {
			fmt.Printf("failed to delete directory '%s': %v", dir, err)
			os.Exit(1)
		}
	},
	Args: cobra.ExactArgs(1),
}

func init() {
	repoCmd.AddCommand(repodelCmd)
}
