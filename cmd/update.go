package cmd

import (
	"fmt"
	"net/http"

	"github.com/fatih/color"
	"github.com/fynelabs/selfupdate"
	"github.com/spf13/cobra"
	"golang.org/x/mod/semver"
)

func doUpdate(url string) error {
	fmt.Printf("Downloading: %s\n", url)
	// request the new file
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP request failed with status code: %d", resp.StatusCode)
	}
	if resp.Body == nil {
		return fmt.Errorf("HTTP request failed with empty body")
	}

	err = selfupdate.Apply(resp.Body, selfupdate.Options{})
	if err != nil {
		if rerr := selfupdate.RollbackError(err); rerr != nil {
			fmt.Printf("Failed to rollback from bad update: %s, %v", url, rerr)
		}
	}
	return err
}

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		repo := "majklovec/caravana"
		latest_version, _ := getLatestVersion(repo)
		fmt.Printf("latest version: %s\r\nactual version: v%s\r\n", latest_version, version)

		a := semver.Compare(latest_version, version)
		if a == 1 {
			err := doUpdate(fmt.Sprintf("https://github.com/%s/releases/download/%s/caravana-%s-%s", repo, latest_version, getEnvOrDefault("GOOS", "linux"), getEnvOrDefault("GOARCH", "amd64")))
			if err != nil {
				color.Red("Failed to update: %v\r\n", err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
