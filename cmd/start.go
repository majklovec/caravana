package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/hashicorp/levant/levant"
	"github.com/hashicorp/levant/levant/structs"
	"github.com/hashicorp/levant/logging"

	"github.com/hashicorp/levant/template"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var Dry bool = false

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start [config]",
	Short: "Start service",
	Args:  cobra.RangeArgs(0, 1),
	Run: func(cmd *cobra.Command, args []string) {
		var configFile string
		if len(args) > 0 {
			configFile = args[0]
		}

		// viper.AutomaticEnv()

		// zerolog.SetGlobalLevel(zerolog.NoLevel)
		logging.SetupLogger("debug", "human")

		if err := processConfig(configFile); err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
	},
}

func init() {
	startCmd.Flags().BoolVarP(&Dry, "dry", "", false, "Do not touch the cluster, only simulate")
	rootCmd.AddCommand(startCmd)
}

func processConfig(configFile string) error {

	v := make(map[string]interface{})
	globalConfig, err := loadConfig("config.yaml")
	if err == nil {
		for key, val := range globalConfig {
			v[key] = val
		}
	}

	var template string

	config, err := loadConfig(filepath.Join(configDir, configFile) + ".yaml")
	if err == nil {
		if _, ok := config["TEMPLATE"]; ok {
			template = config["TEMPLATE"].(string)
		}
	}

	if viper.GetString("TEMPLATE") != "" {
		template = viper.GetString("TEMPLATE")
	}
	if viper.GetString("DOMAIN") != "" {

		configFile = viper.GetString("DOMAIN")
	}
	if template == "" {
		return fmt.Errorf("template %s does not exists", filepath.Join(configDir, configFile)+".yaml")
	}

	templateConfig, err := loadConfig(filepath.Join(templatesDir, template, "config.yaml"))
	if err == nil {
		mergeConfigs(v, templateConfig)
	}

	mergeConfigs(v, config)

	checkIfDirectoryExists := directoryExists(filepath.Join(templatesDir, template))
	if !checkIfDirectoryExists {
		return fmt.Errorf("service template %s does not exists", filepath.Join(templatesDir, template))
	}

	filepath.Walk(filepath.Join(templatesDir, template), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && filepath.Ext(path) == ".nomad" {
			if err := processTemplate(configFile, path, v); err != nil {
				fmt.Printf("Error processing %s: %v\n", path, err)
			}
		}

		return nil
	})

	return nil
}

func processTemplate(configFilePath string, templateFilePath string, templateConfig map[string]interface{}) error {
	config := &levant.DeployConfig{
		Client:   &structs.ClientConfig{},
		Deploy:   &structs.DeployConfig{},
		Plan:     &structs.PlanConfig{},
		Template: &structs.TemplateConfig{},
	}
	config.Template.TemplateFile = templateFilePath

	variables := make(map[string]interface{})
	variables["SERVICE_ID"] = sanitizeServiceID(configFilePath, '-')
	variables["DOMAIN"] = configFilePath

	envs := viper.AllSettings()

	mergeConfigs(variables, templateConfig, envs)

	// parse levant template
	var err error
	config.Template.Job, err = template.RenderJob(config.Template.TemplateFile, config.Template.VariableFiles, config.Client.ConsulAddr, &variables)
	if err != nil {
		return fmt.Errorf("error rendering Levant template: %v", err)
	}

	// output levant template
	if Dry {
		tpl, err := template.RenderTemplate(config.Template.TemplateFile, config.Template.VariableFiles, config.Client.ConsulAddr, &variables)
		if err != nil {
			fmt.Printf("[ERROR] levant/command: %v", err)
			return fmt.Errorf("error rendering Levant template: %v", err)
		}
		if *config.Template.Job.Name != configFilePath {
			return fmt.Errorf("job has to be named using 'job \"[[.DOMAIN]]\"' name should be %s not: %s", configFilePath, *config.Template.Job.Name)
		}
		fmt.Println(tpl)
		return nil
	}

	// ensure that host dirs exists for binding
	for _, taskGroup := range config.Template.Job.TaskGroups {
		for _, task := range taskGroup.Tasks {
			if _, ok := task.Config["mount"]; ok {
				mounts := task.Config["mount"].([]map[string]interface{})
				for _, mount := range mounts {
					if _, ok := mount["source"]; ok {
						ensureDirectoryExists(mount["source"].(string), *config)
					}
				}
			}
		}
	}

	// mark jobs by metadata to find em later
	config.Template.Job.Meta = make(map[string]string)
	config.Template.Job.Meta["CARAVANA_TEMPLATE"] = templateFilePath
	config.Template.Job.Meta["CARAVANA_CONFIG"] = configFilePath

	success := levant.TriggerDeployment(config, client)
	if !success {
		return fmt.Errorf("unable to complete deployment")
	}

	fmt.Println("Deployment successful for Nomad file:", templateFilePath)
	return nil
}
