package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"

	"encoding/json"

	"github.com/aliyun/fcli/util"
)

type configInputType struct {
	endpoint        *string
	accessKeyID     *string
	accessKeySecret *string
	securityToken   *string
	apiVersion      *string
	timeout         *uint
	debug           *bool
	display         *bool
}

var configInput configInputType

func init() {
	RootCmd.AddCommand(configCmd)

	configCmd.Flags().Bool("help", false, "set configuration")
	configInput.endpoint = configCmd.Flags().String("endpoint", "", "fc endpoint")
	configInput.accessKeyID = configCmd.Flags().String("access-key-id", "", "access key id")
	configInput.accessKeySecret = configCmd.Flags().String("access-key-secret", "", "access key secret")
	configInput.securityToken = configCmd.Flags().String("security-token", "", "ram security token")
	configInput.timeout = configCmd.Flags().Uint("timeout", 60, "timeout in seconds")
	configInput.apiVersion = configCmd.Flags().String("api-version", "2016-08-15", "fc api version")
	configInput.debug = configCmd.Flags().Bool("debug", false, "enable debug or not")
	configInput.display = configCmd.Flags().Bool("display", false, "display the configuration")
}

var configCmd = &cobra.Command{
	Use:     "config",
	Aliases: []string{"c"},
	Short:   "Configure the fcli",
	Long:    ``,
	Run: func(cmd *cobra.Command, args []string) {
		if cmd.Flags().NFlag() == 0 {
			readConfig()
			return
		}
		if *configInput.display {
			fmt.Printf(displayConfig())
			return
		}
		config, err := getConfigAways()
		if err != nil {
			return
		}
		if cmd.Flags().Changed("debug") {
			config.Debug = *configInput.debug
		}
		if cmd.Flags().Changed("timeout") {
			config.Timeout = *configInput.timeout
		}
		if cmd.Flags().Changed("endpoint") {
			config.Endpoint = *configInput.endpoint
			config.Endpoint = strings.TrimSpace(config.Endpoint)
			config.SLSEndpoint = fmt.Sprintf(
				util.LogEndpointFmt, util.GetRegionNoForSLSEndpoint(config.Endpoint))
		}

		if cmd.Flags().Changed("access-key-id") {
			config.AccessKeyID = *configInput.accessKeyID
		}
		if cmd.Flags().Changed("access-key-secret") {
			config.AccessKeySecret = *configInput.accessKeySecret
		}
		if cmd.Flags().Changed("security-token") {
			config.SecurityToken = *configInput.securityToken
		}
		if cmd.Flags().Changed("api-version") {
			config.APIVersion = *configInput.apiVersion
		}

		data, err := yaml.Marshal(config)
		if err != nil {
			fmt.Printf("Failed to marshal config: %v. Error: %v\n", data, err)
			return
		}
		err = ioutil.WriteFile(gConfigPath, data, 0600)
		if err != nil {
			fmt.Printf("Failed to write file: %s. Error: %v\n", gConfigPath, err)
			return
		}
	},
}

func displayConfig() string {
	return displayConfigFile() + displayAllEnv()
}

func displayConfigFile() string {
	outputStr := fmt.Sprintf("Config file directory: %s\n", gConfigDir)
	config, _ := getConfigFromFile()
	if config != nil {
		content, _ := json.MarshalIndent(config, "", "  ")
		outputStr += fmt.Sprintf("%s\n", string(content))
	}
	return outputStr
}

func displayAllEnv() string {
	return fmt.Sprintf("Environment variables for fcli: \n") +
		displayEnv("ALIBABA_CLOUD_ACCESS_KEY_ID") +
		displayEnv("ALIBABA_CLOUD_ACCESS_KEY_SECRET") +
		displayEnv("ALIBABA_CLOUD_DEFAULT_REGION") +
		displayEnv("ALIBABA_CLOUD_ACCOUNT_ID") +
		fmt.Sprintf("Environment variables for fcli: (deprecated)\n") +
		displayEnv("ALIYUN_ACCESS_KEY_ID") +
		displayEnv("ALIYUN_ACCESS_KEY_SECRET")
}

func displayEnv(key string) string {
	if value := os.Getenv(key); value != "" {
		return fmt.Sprintf("  %s: %s\n", key, value)
	}
	return fmt.Sprintf("  %s is not set.\n", key)
}

func getConfigFromFile() (*util.GlobalConfig, error) {
	config := util.NewGlobalConfig()
	data, err := ioutil.ReadFile(gConfigPath)
	if err != nil {
		fmt.Printf("Config file does not yet exist: %s\n", gConfigPath)
		return nil, nil
	}
	err = yaml.Unmarshal(data, config)
	if err != nil {
		fmt.Printf("Failed to unmarshal config: %v. Error: %v\n", gConfigPath, err)
		return nil, err
	}
	return config, nil

}

// getConfigAways first try to parse the config object from the config file,
// if the config file does not exist, initialize a config object to return.
func getConfigAways() (*util.GlobalConfig, error) {
	config, err := getConfigFromFile()
	if err != nil {
		return nil, err
	}
	if config == nil {
		config = util.NewGlobalConfig()
		err := os.Mkdir(gConfigDir, 0755)
		if err != nil && os.IsNotExist(err) {
			fmt.Printf("Failed to mkdir: %s. Error: %v\n", gConfigDir, err)
			return nil, err
		}
	}
	return config, nil
}
