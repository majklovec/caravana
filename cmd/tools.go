package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/hashicorp/levant/levant"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

func mergeConfigs(dest map[string]interface{}, srcs ...map[string]interface{}) {
	for _, src := range srcs {
		for key, val := range src {
			dest[key] = val
		}
	}
}
func directoryExists(directoryPath string) bool {
	_, err := os.Stat(directoryPath)
	return err == nil
}

func loadConfig(fileName string) (map[string]interface{}, error) {
	fmt.Printf("Loading config %s\n", fileName)
	configFile, err := os.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %v", err)
	}

	var config map[string]interface{}
	err = yaml.Unmarshal(configFile, &config)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling config: %v", err)
	}

	return config, nil
}

func ensureDirectoryExists(directoryPath string, config levant.DeployConfig) error {
	nomad_addr := viper.GetString("NOMAD_ADDR")
	if nomad_addr == "" {
		// Attempt to create the directory
		err := os.MkdirAll(directoryPath, os.ModePerm)
		if err != nil {
			fmt.Printf("Failed to create directory %s: %v\n", directoryPath, err)
			return err
		}
		color.Yellow("Created local directory %s\n", directoryPath)
	} else {
		color.Blue("mkdir -p %s\n", directoryPath)
		// config.Template.Job.TaskGroups = append(config.Template.Job.TaskGroups, &api.TaskGroup{
		// 	Name: StringPtr("create-dir" + sanitizeServiceID(directoryPath, '-')),
		// 	Tasks: []*api.Task{
		// 		{
		// 			Name:   "create-dir",
		// 			Driver: "raw_exec",
		// 			Config: map[string]interface{}{
		// 				"command": fmt.Sprintf("mkdir -p %s", directoryPath),
		// 			},
		// 		},
		// 	},
		// })

	}

	return nil
}

func sanitizeServiceID(serviceID string, replace rune) string {
	return strings.Map(func(r rune) rune {
		if ('a' <= r && r <= 'z') || ('A' <= r && r <= 'Z') || ('0' <= r && r <= '9') || r == '-' {
			return r
		}
		return replace
	}, serviceID)
}

func getEnvOrDefault(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}

func saveConfig(filePath string, data deployment) error {
	yamlData, err := yaml.Marshal(data)
	if err != nil {
		return err
	}

	err = os.WriteFile(filePath, yamlData, 0644)
	if err != nil {
		return err
	}

	fmt.Printf("Config %s has been created.\n", filePath)
	return nil
}

func extractRepoName(gitURL string) string {
	// Assuming Git URL is in the format: https://github.com/username/repo.git
	parts := strings.Split(gitURL, "/")
	repoWithGitSuffix := parts[len(parts)-1]
	repoName := strings.TrimSuffix(repoWithGitSuffix, ".git")
	return strings.TrimPrefix(repoName, "caravana-")
}

func getLatestVersion(repo string) string {
	url := fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", repo)

	// create a request with basic-auth
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("User-Agent", "metal3d-go-client")
	req.Header.Add("Accept", "application/vnd.github.v3.text-match+json")
	req.Header.Add("Accept", "application/vnd.github.moondragon+json")

	// call github
	client := http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.Fatal("Error while making request", err)
	}

	// status in <200 or >299
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		log.Fatalf("Error: %d %s", resp.StatusCode, resp.Status)
	}

	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading response", err)
	}
	result := make(map[string]interface{})
	json.Unmarshal(bodyText, &result)

	version := result["tag_name"].(string)

	return version
}
