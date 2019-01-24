package cmd

import (
	"io/ioutil"

	"github.com/aliyun/fc-go-sdk"

	"fmt"

	"github.com/spf13/cobra"

	"github.com/aliyun/fcli/util"
)

type updateFuncInputType struct {
	serviceName           *string
	functionName          *string
	description           *string
	runtime               *string
	handler               *string
	initializer           *string
	codeDir               *string
	codeFile              *string
	codeOSSBucket         *string
	codeOSSObject         *string
	memory                *int32
	timeout               *int32
	initializationTimeout *int32
	etag                  *string
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
}

var updateFuncCmd = &cobra.Command{
	Use:     "update [option]",
	Aliases: []string{"u"},
	Short:   "Update function attributes",
	Long:    ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		input := fc.NewUpdateFunctionInput(*updateFuncInput.serviceName, *updateFuncInput.functionName)
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
				return fmt.Errorf("%v", err)
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
			return fmt.Errorf("can not create fc client: %s\n", err)
		}
		_, err = client.UpdateFunction(input)
		if err != nil {
			return fmt.Errorf("%s\n", err)
		}
		return nil
	},
}
