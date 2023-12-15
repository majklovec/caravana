/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	cmd2 "github.com/jsiebens/hashi-up/cmd"
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("create called")
		// hashi-up nomad install \
		// --ssh-target-addr $AGENT_1_IP \
		// --ssh-target-user ubuntu \
		// --client \
		// --retry-join $SERVER_IP
		fmt.Println(cmd.Args)
		nomadmd := cmd2.InstallNomadCommand()
		nomadmd.SetArgs([]string{
			"install",
			"--ssh-target-addr",
			"192.168.1.10",
			"--ssh-target-user",
			"ubuntu",
			"--client",
			"--retry-join",
			"192.168.1.10"})
		fmt.Println(nomadmd)
		nomadmd.Execute()
	},
}

func init() {
	clusterCmd.AddCommand(createCmd)
}
