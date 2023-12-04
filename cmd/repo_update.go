package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
)

// isGitRepository checks if the specified directory contains a Git repository.
func isGitRepository(dir string) bool {
	_, err := git.PlainOpen(dir)
	return err == nil
}

// pullRepository pulls the latest changes from the Git repository.
func pullRepository(dir string) error {
	r, err := git.PlainOpen(dir)
	if err != nil {
		return err
	}

	w, err := r.Worktree()
	if err != nil {
		return err
	}

	fmt.Printf("Pulling changes for repository at '%s'\n", dir)

	err = w.Pull(&git.PullOptions{
		// RemoteName:    "origin",
		// ReferenceName: plumbing.ReferenceName("refs/heads/master"), // Change to the desired branch
	})
	if err != nil && err != git.NoErrAlreadyUpToDate {
		return err
	}

	return nil
}

func pullRepositories(rootDir string) error {
	entries, err := os.ReadDir(rootDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			dirPath := filepath.Join(rootDir, entry.Name())
			if isGitRepository(dirPath) {
				err := pullRepository(dirPath)
				if err != nil {
					fmt.Printf("Error pulling repository at '%s': %v\n", dirPath, err)
				}
			}
		}
	}

	return nil
}

// repoupdateCmd represents the update command
var repoupdateCmd = &cobra.Command{
	Use:   "update",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("update called")
		pullRepositories(templatesDir)
	},
	Args: cobra.ExactArgs(0),
}

func init() {
	repoCmd.AddCommand(repoupdateCmd)
}
