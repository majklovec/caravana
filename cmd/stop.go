package cmd

import (
	"fmt"

	"github.com/hashicorp/nomad/api"
	"github.com/spf13/cobra"
)

func stopNomadJob(jobName string) error {
	id, _, err := client.Jobs().Deregister(jobName, false, &api.WriteOptions{})

	if id == "" {
		return fmt.Errorf("job does not exist")
	}
	if err != nil {
		return fmt.Errorf("failed to stop job: %v", err)
	}
	fmt.Printf("Job %s (%s) stopped\n", jobName, id)
	return nil
}

// stopCmd represents the stop command
var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Will stop given job",
	Long:  `Will stop given job`,

	Run: func(cmd *cobra.Command, args []string) {
		err := stopNomadJob(args[0])
		if err != nil {
			fmt.Println(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)
}
