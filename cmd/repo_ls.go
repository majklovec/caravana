package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/ryanuber/columnize"
	"github.com/spf13/cobra"
)

func splitTitleAndDescription(filename string) (string, string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var title, description string

	for scanner.Scan() {
		line := scanner.Text()

		// Trim spaces from the beginning and end of the line
		line = strings.TrimSpace(line)

		// Check if the line starts with '#'
		if strings.HasPrefix(line, "#") {
			// Extract title by removing the '#' and trimming spaces
			title = strings.TrimSpace(strings.TrimPrefix(line, "#"))
		} else {
			// Set the description until the first new line
			if line != "" {
				description = line
				break
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return "", "", err
	}

	// Trim spaces from the beginning and end of the description
	description = "| " + strings.TrimSpace(description)
	title = "| " + strings.TrimSpace(title)
	return title, description, nil
}
func listDirectories(root string, maxDepth int) ([]string, error) {
	blue := color.New(color.FgBlue).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	out := make([]string, 0)
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip the root directory and.git directory
		if path == root || info.Name() == ".git" {
			return nil
		}

		depth := strings.Count(strings.TrimPrefix(path, root), string(os.PathSeparator))

		if depth <= maxDepth && info.IsDir() {
			if depth == maxDepth {

				name, description, err := splitTitleAndDescription(path + "/" + "README.md")
				if err != nil {
					description = ""
				}
				out = append(out, fmt.Sprintf("%s %s %s", blue(strings.TrimPrefix(path, root)[1:]), yellow(name), description))

			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}
	return out, nil
}

// repolsCmd represents the ls command
var repolsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List repos and available service templates",
	Long:  `List repos and available service templates`,
	Run: func(cmd *cobra.Command, args []string) {
		maxDepth := 2 // Replace with the maximum depth you desire

		output, err := listDirectories(templatesDir, maxDepth)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

		result := columnize.SimpleFormat(output)
		fmt.Println(result)
	},
	Args: cobra.ExactArgs(0),
}

func init() {
	repoCmd.AddCommand(repolsCmd)
}
