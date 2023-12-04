package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/hashicorp/nomad/api"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var version = "dev"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "caravana",
	Version: version,
	Short:   "Deployment tool for HashiCorp Nomad",
	Long:    `Caravana is a deployment tool for HashiCorp Nomad.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		directoryPath := "configs"
		exists := directoryExists(directoryPath)
		if !exists {
			color.Red(`No configs present! 
			
1) Are you in the right directory? 
2) If it's your first run, create ./config directory and put your yaml config files there
  example:			

  cat > configs/example.yaml
  SERVICE: test-repo/gitea
  DOMAIN: git.vondracek.dev
			
			`)
			cmd.Help()
			os.Exit(1)
		}
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var client *api.Client
var templatesDir string = "templates"
var configDir string = "configs"

func init() {
	config := api.DefaultConfig()
	var err error
	client, err = api.NewClient(config)
	if err != nil {
		fmt.Printf("failed to create Nomad client: %v", err)
	}

	viper.AutomaticEnv()
}
