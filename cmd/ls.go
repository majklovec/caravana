package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/hashicorp/nomad/api"
	"github.com/ryanuber/columnize"
	"github.com/spf13/cobra"
)

func listRunningJobs() ([]*api.JobListStub, error) {

	jobListOptions := &api.JobListOptions{
		Fields: &api.JobListFields{
			Meta: true,
		},
	}

	queryOptions := &api.QueryOptions{
		AllowStale: true,
	}

	// Query all jobs with the specified options
	jobs, _, err := client.Jobs().ListOptions(jobListOptions, queryOptions)
	if err != nil {
		return nil, err
	}

	return jobs, nil
}

func isRunning(configName string, runningJobs []*api.JobListStub) bool {
	found := false
	for _, job := range runningJobs {
		if configName == job.Meta["CARAVANA_CONFIG"] {
			found = true
			break
		}
	}
	return found
}

func listConfigs(root string, runningJobs []*api.JobListStub) ([]string, error) {
	out := make([]string, 0)
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories and list only files
		if !info.IsDir() {
			filename := info.Name()

			configName := filepath.Base(filename[:len(filename)-len(filepath.Ext(filename))])

			if !isRunning(configName, runningJobs) {
				out = append(out, fmt.Sprintf("%s | not started", configName))
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}
	return out, nil
}

// lsCmd represents the ls command
var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List configs and running jobs",
	Long:  `List configs and running jobs`,
	Run: func(cmd *cobra.Command, args []string) {
		green := color.New(color.FgGreen).SprintFunc()
		red := color.New(color.FgRed).SprintFunc()
		// blue := color.New(color.FgBlue).SprintFunc()
		yellow := color.New(color.FgYellow).SprintFunc()

		runningJobs, err := listRunningJobs()
		if err != nil {
			log.Fatal("Error:", err)
		}

		// Filter running jobs
		out := make([]string, 0)
		for _, job := range runningJobs {
			if job.Meta["CARAVANA_CONFIG"] != "" {
				// runningJobs = append(runningJobs, job)

				if job.Status == "running" {
					out = append(out, fmt.Sprintf("%s | %s", job.Meta["CARAVANA_CONFIG"], green(job.Status)))
					continue
				}
				if job.Status == "dead" {
					out = append(out, fmt.Sprintf("%s | %s", job.Meta["CARAVANA_CONFIG"], red(job.Status)))
					continue
				}
				if job.Status == "pending" {
					out = append(out, fmt.Sprintf("%s | %s", job.Meta["CARAVANA_CONFIG"], yellow(job.Status)))
					continue
				}
				out = append(out, fmt.Sprintf("%s | %s", job.Meta["CARAVANA_CONFIG"], job.Status))

			}
		}

		notStarted, err := listConfigs("configs", runningJobs)
		if err != nil {
			fmt.Println("Error:", err)
		}

		out = append(out, notStarted...)

		// fmt.Printf("%+v\n", out)
		result := columnize.SimpleFormat(out)
		fmt.Printf("%s\r\n\r\n", result)

	},
	Args: cobra.ExactArgs(0),
}

func init() {
	rootCmd.AddCommand(lsCmd)

}
