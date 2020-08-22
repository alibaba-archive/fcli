package cmd

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/aliyun/fc-go-sdk"
	"github.com/aliyun/fcli/util"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type createFuncInputType struct {
	serviceName            string
	functionName           string
	description            string
	runtime                string
	handler                string
	initializer            string
	codeDir                string
	codeFile               string
	codeOSSBucket          string
	codeOSSObject          string
	customContainerImage   string
	customContainerCommand string
	customContainerArgs    string
	caPort                 int32
	memory                 int32
	timeout                int32
	initializationTimeout  int32
	environmentVariables   []string
	environmentConfigFiles []string
}

// Use a unique name to avoid global variable confliction.
var createFuncInput createFuncInputType

func init() {
	functionCmd.AddCommand(createFuncCmd)

	createFuncCmd.Flags().Bool("help", false, "")
	createFuncCmd.Flags().StringVarP(&createFuncInput.serviceName, "service-name", "s", "", "the service name")
	createFuncCmd.Flags().StringVarP(&createFuncInput.functionName, "function-name", "f", "", "the function name")
	createFuncCmd.Flags().StringVar(&createFuncInput.description, "description", "", "brief description")
	createFuncCmd.Flags().Int32VarP(&createFuncInput.memory, "memory", "m", 128, "memory size in MB")
	createFuncCmd.Flags().Int32Var(&createFuncInput.timeout, "timeout", 30, "timeout in seconds")
	createFuncCmd.Flags().Int32VarP(&createFuncInput.initializationTimeout, "initializationTimeout", "e", 30, "timeout in seconds")
	createFuncCmd.Flags().StringVarP(&createFuncInput.codeOSSBucket, "code-bucket", "b", "", "oss bucket of the code")
	createFuncCmd.Flags().StringVarP(&createFuncInput.codeOSSObject, "code-object", "o", "", "oss object of the code")
	createFuncCmd.Flags().StringVarP(&createFuncInput.customContainerImage, "custom-container-image", "g", "", "custom container config image")
	createFuncCmd.Flags().StringVarP(&createFuncInput.customContainerCommand, "custom-container-command", "n", "", "custom container config command, e.g. [\"node\"]")
	createFuncCmd.Flags().StringVarP(&createFuncInput.customContainerArgs, "custom-container-args", "r", "", "custom container config args, e.g. [\"server.js\"]")
	createFuncCmd.Flags().Int32Var(&createFuncInput.caPort, "ca-port", 0, "args of custom container config")

	createFuncCmd.Flags().StringVarP(
		&createFuncInput.codeDir, "code-dir", "d", "",
		"function code directory. If both code-file and code-dir are provided, code-file will be used.")
	createFuncCmd.Flags().StringVar(
		&createFuncInput.codeFile, "code-file", "",
		"zipped code file. If both code-file and code-dir are provided, code-file will be used.")
	createFuncCmd.Flags().StringVarP(&createFuncInput.runtime, "runtime", "t", "", "function runtime")
	createFuncCmd.Flags().StringVarP(
		&createFuncInput.handler, "handler", "h", "", "handler is the entrypoint for the function execution")
	createFuncCmd.Flags().StringVarP(
		&createFuncInput.initializer, "initializer", "i", "", "initializer is the entrypoint for the initializer execution")
	createFuncCmd.Flags().StringArrayVar(&createFuncInput.environmentVariables, "env", []string{}, "set environment variables. e.g. --env VAR1=val1 --env VAR2=val2")
	createFuncCmd.Flags().StringArrayVar(&createFuncInput.environmentConfigFiles, "env-file", []string{}, "read in a file of environment variables. e.g. --env-file FILE1 --env-file FILE2")
}

var createFuncCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"c"},
	Short:   "Create function",
	Long:    ``,
	Run: func(cmd *cobra.Command, args []string) {
		input := fc.NewCreateFunctionInput(createFuncInput.serviceName).
			WithFunctionName(createFuncInput.functionName).
			WithDescription(createFuncInput.description).
			WithMemorySize(createFuncInput.memory).
			WithTimeout(createFuncInput.timeout).
			WithInitializationTimeout(createFuncInput.initializationTimeout).
			WithHandler(createFuncInput.handler).
			WithInitializer(createFuncInput.initializer).
			WithRuntime(createFuncInput.runtime)

		envMap := make(map[string]string)

		if cmd.Flags().Changed("env-file") {
			for _, envFilePath := range createFuncInput.environmentConfigFiles {
				_, err := util.GetEnvSetting(envMap, envFilePath)
				if err != nil {
					fmt.Printf("Error: %s\n", err)
					return
				}
			}
			input.WithEnvironmentVariables(envMap)
		}

		if cmd.Flags().Changed("env") {
			for _, envVar := range createFuncInput.environmentVariables {
				config := strings.Split(envVar, "=")
				if len(config) == 2 {
					envMap[config[0]] = config[1]
				}
			}
			input.WithEnvironmentVariables(envMap)
		}

		if createFuncInput.runtime == util.RuntimeCustomContainer {
			input = createFunctionInputWithCustomContainerConfig(cmd.Flags(), input, createFuncInput.customContainerImage,
				createFuncInput.customContainerCommand, createFuncInput.customContainerArgs)
		} else {
			if createFuncInput.codeFile != "" {
				var data []byte
				data, err := ioutil.ReadFile(createFuncInput.codeFile)
				if err != nil {
					fmt.Printf("Error: %s\n", err)
					return
				}
				input.WithCode(fc.NewCode().WithZipFile(data))
			} else if createFuncInput.codeDir != "" {
				input.WithCode(fc.NewCode().WithDir(createFuncInput.codeDir))
			} else {
				input.WithCode(fc.NewCode().
					WithOSSBucketName(createFuncInput.codeOSSBucket).
					WithOSSObjectName(createFuncInput.codeOSSObject))
			}
		}

		if cmd.Flags().Changed("ca-port") {
			input.WithCAPort(createFuncInput.caPort)
		}

		client, err := util.NewFClient(gConfig)
		if err != nil {
			fmt.Printf("Error: can not create fc client: %s\n", err)
			return
		}
		_, err = client.CreateFunction(input)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
		}
	},
}

func createFunctionInputWithCustomContainerConfig(flags *pflag.FlagSet, input *fc.CreateFunctionInput,
	image, command, args string) *fc.CreateFunctionInput {
	if flags.Changed("custom-container-image") {
		input.WithCustomContainerConfig(fc.NewCustomContainerConfig().
			WithImage(image))
		if flags.Changed("custom-container-command") {
			input.CustomContainerConfig.WithCommand(command)
		}
		if flags.Changed("custom-container-args") {
			input.CustomContainerConfig.WithArgs(args)
		}
	}

	return input
}
