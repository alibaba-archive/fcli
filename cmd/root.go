package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/AlecAivazis/survey.v1"
	"gopkg.in/yaml.v2"

	"github.com/aliyun/fcli/util"
)

type answer struct {
	AccountID       string
	Region          string
	AccessKeyID     string
	AccessKeySecret string
}

func (a *answer) build(config *util.GlobalConfig) {
	config.AccessKeyID = a.AccessKeyID
	config.AccessKeySecret = a.AccessKeySecret
	config.Endpoint = fmt.Sprintf(util.EndpointFmt, a.AccountID, a.Region)
	config.SLSEndpoint = fmt.Sprintf(util.LogEndpointFmt, util.GetRegionNoForSLSEndpoint(config.Endpoint))
}

// TODO: Replace all other config with gConfig.
var gConfig *util.GlobalConfig
var gConfigDir string
var gConfigPath string

//RootCmd is the root, which is root of all the command variable names
var RootCmd = &cobra.Command{
	Use:   "fcli",
	Short: "fcli: function compute command line tools",
	Long:  `fcli: function compute command line tools`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if !checkConfigCommandWithoutFlags(cmd) && !checkConfigRequredExist() {
			readConfig()
			initConfig()
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			return
		}
	},
}

//Execute method is entrance of this program package
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}

func init() {
	initConfig()
}

func readConfig() {
	config, err := getConfigAways()
	if err != nil {
		return
	}

	accountID, _ := util.GetUIDFromEndpoint(config.Endpoint)

	qs := []*survey.Question{
		{
			Name: "AccountID",
			Prompt: &survey.Input{
				Message: "Alibaba Cloud Account ID",
				Default: accountID,
			},
			Validate: survey.Required,
			Transform: survey.TransformString(func(ans string) string {
				return strings.TrimSpace(ans)
			}),
		},
		{
			Name: "AccessKeyID",
			Prompt: &survey.Input{
				Message: "Alibaba Cloud Access Key ID",
				Default: mark(config.AccessKeyID),
			},
			Validate: survey.Required,
			Transform: survey.TransformString(func(ans string) string {
				if ans == mark(config.AccessKeyID) {
					return config.AccessKeyID
				}
				return strings.TrimSpace(ans)
			}),
		},
		{
			Name: "AccessKeySecret",
			Prompt: &survey.Input{
				Message: "Alibaba Cloud Access Key Secret",
				Default: mark(config.AccessKeySecret),
			},
			Validate: survey.Required,
			Transform: survey.TransformString(func(ans string) string {
				if ans == mark(config.AccessKeySecret) {
					return config.AccessKeySecret
				}
				return strings.TrimSpace(ans)
			}),
		},
		{
			Name: "Region",
			Prompt: &survey.Select{
				Message: "Default region name",
				Options: util.GetRegions(),
				Default: util.GetRegionNoForEndpoint(config.Endpoint),
			},
			Validate: survey.Required,
		},
	}

	ans := &answer{}

	err = survey.Ask(qs, ans)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	ans.build(config)

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

	fmt.Printf("Store the configuration in: %s\n", gConfigDir)
}

func initConfig() {
	gConfig = util.NewGlobalConfig()

	home := os.Getenv("HOME")
	gConfigDir = path.Join(home, ".fcli")
	gConfigPath = path.Join(gConfigDir, "config.yaml")

	pickupConfigFromConfigFile()
	pickupConfigFromEnv()
	pickupConfigFromOldEnv()
}

func checkConfigCommandWithoutFlags(cmd *cobra.Command) bool {
	if cmd.Use == "config" && cmd.Flags().NFlag() == 0 {
		return true
	}
	return false
}

func checkConfigRequredExist() bool {
	if gConfig.AccessKeyID == "" || gConfig.AccessKeySecret == "" || gConfig.Endpoint == "" {
		return false
	}
	return true
}

func pickupConfigFromConfigFile() {
	data, err := ioutil.ReadFile(gConfigPath)
	if err == nil {
		err = yaml.Unmarshal(data, gConfig)
		if err != nil {
			fmt.Printf("Failed to unmarshal config: %v. Error: %v\n", gConfigPath, err)
		}
	}
}

func pickupConfigFromEnv() {
	if accessKey := os.Getenv("ALIBABA_CLOUD_ACCESS_KEY_ID"); accessKey != "" {
		gConfig.AccessKeyID = accessKey
	}
	if accessKeySecret := os.Getenv("ALIBABA_CLOUD_ACCESS_KEY_SECRET"); accessKeySecret != "" {
		gConfig.AccessKeySecret = accessKeySecret
	}

	if regionID := os.Getenv("ALIBABA_CLOUD_DEFAULT_REGION"); regionID != "" {
		if accountID := os.Getenv("ALIBABA_CLOUD_ACCOUNT_ID"); accountID != "" {
			gConfig.Endpoint = fmt.Sprintf(util.EndpointFmt, accountID, regionID)
			gConfig.SLSEndpoint = fmt.Sprintf(util.LogEndpointFmt, util.GetRegionNoForSLSEndpoint(gConfig.Endpoint))

		}
	}
}

func pickupConfigFromOldEnv() {
	if accessKey := os.Getenv("ALIYUN_ACCESS_KEY_ID"); accessKey != "" {
		gConfig.AccessKeyID = accessKey
	}
	if accessKeySecret := os.Getenv("ALIYUN_ACCESS_KEY_SECRET"); accessKeySecret != "" {
		gConfig.AccessKeySecret = accessKeySecret
	}
}

func mark(source string) string {
	if source == "" {
		return ""
	}
	subStr := source
	sourceLen := len(source)
	if sourceLen >= 4 {
		subStr = source[sourceLen-4 : sourceLen]
	}
	return fmt.Sprintf("***********%s", subStr)
}
