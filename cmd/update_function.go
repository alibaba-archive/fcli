package cmd

import (
	"github.com/aliyun/fc-go-sdk"
	"io/ioutil"
	"strings"

	"fmt"

	"github.com/spf13/cobra"

	"github.com/aliyun/fcli/util"
)

type updateFuncInputType struct {
	serviceName            *string
	functionName           *string
	description            *string
	runtime                *string
	handler                *string
	initializer            *string
	codeDir                *string
	codeFile               *string
	codeOSSBucket          *string
	codeOSSObject          *string
	memory                 *int32
	timeout                *int32
	initializationTimeout  *int32
	etag                   *string
	environmentVariables   *[]string
	environmentConfigFiles *[]string
}

var updateFuncInput updateFuncInputType

func init() {
	functionCmd.AddCommand(updateFuncCmd)

	updateFuncCmd.Flags().Bool("help", false, "Print Usage")
	updateFuncInput.serviceName = updateFuncCmd.Flags().StringP("service-name", "s", "", "the service name")
	updateFuncInput.functionName = updateFuncCmd.Flags().StringP("function-name", "f", "", "the function name")
	updateFuncInput.description = updateFuncCmd.Flags().StringP("description", "d", "", "function description")
	updateFuncInput.runtime = updateFuncCmd.Flags().StringP("runtime", "t", "", "function runtime")
	updateFuncInput.timeout = updateFuncCmd.Flags().Int32("timeout", 0, "function timeout in seconds")
	updateFuncInput.initializationTimeout = updateFuncCmd.Flags().Int32P("initializationTimeout", "e", 0, "initializer timeout in seconds")
	updateFuncInput.handler = updateFuncCmd.Flags().StringP("handler", "h", "", "function handler")
	updateFuncInput.initializer = updateFuncCmd.Flags().StringP("initializer", "i", "", "function initializer")
	updateFuncInput.memory = updateFuncCmd.Flags().Int32P("memory", "m", 0, "memory size in MB")
	updateFuncInput.codeOSSBucket = updateFuncCmd.Flags().StringP("bucket", "b", "", "oss code bucket")
	updateFuncInput.codeOSSObject = updateFuncCmd.Flags().StringP("object", "o", "", "oss code object")
	updateFuncInput.codeDir = updateFuncCmd.Flags().String(
		"code-dir", "", "function code directory. If both code-file and code-dir are provided, "+
			"code-file will be used.")
	updateFuncInput.codeFile = updateFuncCmd.Flags().String(
		"code-file", "", "zipped code file. If both code-file and code-dir are provided, "+
			"code-file will be used.")
	updateFuncInput.etag = updateFuncCmd.Flags().String(
		"etag", "", "provide etag to do the conditional update. "+
			"If the specified etag does not match the function's, the update will fail.")
	updateFuncInput.environmentVariables = updateFuncCmd.Flags().StringArray("env", []string{}, "set environment variables. e.g. --env VAR1=val1 --env VAR2=val2")
	updateFuncInput.environmentConfigFiles = updateFuncCmd.Flags().StringArray("env-file", []string{}, "read in a file of environment variables. e.g. --env-file FILE1 --env-file FILE2")
}

var updateFuncCmd = &cobra.Command{
	Use:     "update [option]",
	Aliases: []string{"u"},
	Short:   "Update function attributes",
	Long:    ``,
	Run: func(cmd *cobra.Command, args []string) {
		input := fc.NewUpdateFunctionInput(*updateFuncInput.serviceName, *updateFuncInput.functionName)

		envMap := make(map[string]string)

		if cmd.Flags().Changed("env-file") {
			for _, envFilePath := range *updateFuncInput.environmentConfigFiles {
				_, err := util.GetEnvSetting(envMap, envFilePath)
				if err != nil {
					fmt.Printf("Error: %s\n", err)
					return
				}
			}
			input.WithEnvironmentVariables(envMap)
		}

		if cmd.Flags().Changed("env") {
			for _, envVar := range *updateFuncInput.environmentVariables {
				config := strings.Split(envVar, "=")
				if len(config) == 2 {
					envMap[config[0]] = config[1]
				}
			}
			input.WithEnvironmentVariables(envMap)
		}

		if cmd.Flags().Changed("description") {
			input.WithDescription(*updateFuncInput.description)
		}
		if cmd.Flags().Changed("etag") {
			input.WithIfMatch(*updateFuncInput.etag)
		}
		if cmd.Flags().Changed("memory") {
			input.WithMemorySize(*updateFuncInput.memory)
		}
		if cmd.Flags().Changed("timeout") {
			input.WithTimeout(*updateFuncInput.timeout)
		}
		if cmd.Flags().Changed("initializationTimeout") {
			input.WithInitializationTimeout(*updateFuncInput.initializationTimeout)
		}
		if cmd.Flags().Changed("handler") {
			input.WithHandler(*updateFuncInput.handler)
		}
		if cmd.Flags().Changed("initializer") {
			input.WithInitializer(*updateFuncInput.initializer)
		}
		if cmd.Flags().Changed("runtime") {
			input.WithRuntime(*updateFuncInput.runtime)
		}
		if cmd.Flags().Changed("code-file") {
			data, err := ioutil.ReadFile(*updateFuncInput.codeFile)
			if err != nil {
				fmt.Printf("Error: %v", err)
				return
			}
			input.WithCode(fc.NewCode().WithZipFile(data))
		} else if cmd.Flags().Changed("code-dir") {
			input.WithCode(fc.NewCode().WithDir(*updateFuncInput.codeDir))
		} else if cmd.Flags().Changed("bucket") && cmd.Flags().Changed("object") {
			input.WithCode(fc.NewCode().
				WithOSSBucketName(*updateFuncInput.codeOSSBucket).
				WithOSSObjectName(*updateFuncInput.codeOSSObject))
		}

		client, err := util.NewFClient(gConfig)
		if err != nil {
			fmt.Printf("Error: can not create fc client: %s\n", err)
			return
		}
		_, err = client.UpdateFunction(input)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
		}
	},
}
