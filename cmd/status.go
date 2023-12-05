package cmd

import (
	"fmt"
	"strings"

	"github.com/hashicorp/nomad/api"
	"github.com/ryanuber/columnize"
	"github.com/spf13/cobra"
)

func formatJobInfo(job *api.Job) string {
	// Format the job info
	info := []string{
		fmt.Sprintf("ID: | %s", *job.ID),
		fmt.Sprintf("Name: | %s", *job.Name),
		fmt.Sprintf("Datacenters: | %s", strings.Join(job.Datacenters, ",")),
		fmt.Sprintf("Namespace: | %s", *job.Namespace),
		fmt.Sprintf("STATUS: | %s", *job.Status),
	}

	return columnize.SimpleFormat(info)
}

func getJobInfo(jobID string) error {
	q := &api.QueryOptions{}
	job, _, err := client.Jobs().Info(jobID, q)
	if err != nil {
		return fmt.Errorf("error querying job: %s", err)
	}

	fmt.Println(formatJobInfo(job))

	return nil
}

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := getJobInfo(args[0])
		if err != nil {
			fmt.Printf("Error: %s\n", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
