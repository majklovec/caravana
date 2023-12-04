package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
)

func plainClone(url string) {
	dir := filepath.Join(templatesDir, extractRepoName(url))
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		fmt.Printf("Failed to create directory %s: %v\n", dir, err)
	}

	// Clones the repository into the given dir, just as a normal git clone does
	r, err := git.PlainClone(dir, false, &git.CloneOptions{
		URL:          url,
		SingleBranch: true,
		Depth:        1,
		Progress:     os.Stdout,
	})

	if err != nil {
		fmt.Printf("Failed to clone %s: %v\n", dir, err)
	}
	fmt.Printf("Cloned %s: %v\n", dir, r)
}

// repoaddCmd represents the add command
var repoaddCmd = &cobra.Command{
	Use:   "add [url/dir]",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("add called")
		plainClone(args[0])

	},
	Args: cobra.ExactArgs(1),
}

func init() {
	repoCmd.AddCommand(repoaddCmd)
}
