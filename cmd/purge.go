package cmd

import (
	"fmt"

	"github.com/hashicorp/nomad/api"
	"github.com/spf13/cobra"
)

func purgeNomadJob(jobName string) error {

	id, _, err := client.Jobs().Deregister(jobName, true, &api.WriteOptions{})

	if id == "" {
		return fmt.Errorf("job does not exist")
	}
	if err != nil {
		return fmt.Errorf("failed to purge job: %v", err)
	}
	fmt.Printf("Job %s (%s) purged\n", jobName, id)
	return nil
}

// purgeCmd represents the purge command
var purgeCmd = &cobra.Command{
	Use:   "purge",
	Short: "Will stop and purge given job",
	Long:  `Will stop and purge given job`,
	Run: func(cmd *cobra.Command, args []string) {
		err := purgeNomadJob(args[0])
		if err != nil {
			fmt.Println(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(purgeCmd)
}
