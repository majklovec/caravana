/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/hashicorp/levant/helper"
	"github.com/hashicorp/levant/levant"
	"github.com/hashicorp/levant/levant/structs"
	"github.com/hashicorp/levant/template"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "A Golang CLI app to run Levant templates",
	Run: func(cmd *cobra.Command, args []string) {
		if err := runLevantTemplate(); err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
	},
}

func getEnvOrDefault(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}
func runLevantTemplate() error {
	config := &levant.DeployConfig{
		Client:   &structs.ClientConfig{},
		Deploy:   &structs.DeployConfig{},
		Plan:     &structs.PlanConfig{},
		Template: &structs.TemplateConfig{},
	}

	// Customize Levant config based on your needs
	config.Client.Addr = getEnvOrDefault("NOMAD_ADDR", "http://localhost:4646")

	// Set other Levant options as needed

	if templateFile != "" {
		config.Template.TemplateFile = templateFile
	} else {
		if config.Template.TemplateFile = helper.GetDefaultTmplFile(); config.Template.TemplateFile == "" {
			return fmt.Errorf("template_file missing and no default template found")
		}
	}
	v := make(map[string]interface{})
	v["SERVICE_ID"] = "gitea"
	var err error // Add this line to declare the 'err' variable
	config.Template.Job, err = template.RenderJob(config.Template.TemplateFile, config.Template.VariableFiles, config.Client.ConsulAddr, &v)
	if err != nil {
		return fmt.Errorf("error rendering Levant template: %v", err)
	}

	// Customize the Levant job as needed
	// for _, taskGroup := range config.Template.Job.TaskGroups {
	// 	for _, task := range taskGroup.Tasks {
	// 		// Customize task settings based on your requirements
	// 	}
	// }

	// Trigger deployment using Levant
	success := levant.TriggerDeployment(config, nil)
	if !success {
		return fmt.Errorf("unable to complete deployment")
	}

	fmt.Println("Deployment successful!")
	return nil
}

// Flags
var templateFile string

func init() {

	// Add Levant-related flags (customize based on your needs)
	startCmd.PersistentFlags().StringVar(&templateFile, "template-file", "", "Path to Levant template file")
	startCmd.PersistentFlags().Int("canary_auto_promote", 0, "Canary auto promote")

	rootCmd.AddCommand(startCmd)
}
