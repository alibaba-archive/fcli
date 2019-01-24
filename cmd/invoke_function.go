package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"

	"github.com/aliyun/fc-go-sdk"
	"github.com/aliyun/fcli/util"
)

func init() {
	invokeFuncInput = fc.NewInvokeFunctionInput("", "")
	functionCmd.AddCommand(invokeFuncCmd)
}

var invokeFuncInput *fc.InvokeFunctionInput
var invocationType string
var eventStr string
var eventFile string
var invocationOutputFile string
var invkDebugEnabled bool
var qualifier string

// InvokeFunctionError ...
type InvokeFunctionError struct {
	ErrorCode    string
	ErrorMessage []byte
}

var invokeFuncCmd = &cobra.Command{
	Use:     "invoke",
	Aliases: []string{"i"},
	Short:   "Invoke function",
	Long: `
invoke function
EXAMPLE:
You can provide event with --event-file or --event-str, if both will use --event-str
fcli function invoke -s "service_name" -f "function_name"
		    --debug
		    --output "output_filename"
		    --invocation-type "Async|Sync"
		    --event-file "event file"
		    --event-str  "event_string"
	            --qualifier  "LATEST"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		resp, err := invokeFuncRun()
		if err != nil {
			return fmt.Errorf("%s\n", err)
		}

		var output []byte
		if !invkDebugEnabled {
			output = resp.Payload
		} else {
			output = []byte(resp.String())
			if resp.Payload != nil {
				var invokeError InvokeFunctionError
				json.Unmarshal(resp.Payload, &invokeError)
				output = append(output, '\n')
				output = append(output, []byte(string(invokeError.ErrorCode)+"\n")...)
				output = append(output, invokeError.ErrorMessage...)
			}
		}

		if invocationOutputFile == "" {
			fmt.Println(string(output))
			return nil
		}
		err = ioutil.WriteFile(invocationOutputFile, output, 0644)
		if err != nil {
			return fmt.Errorf("output error:%s\n", err)
		}
		return nil
	},
}

func invokeFuncRun() (*fc.InvokeFunctionOutput, error) {
	err := prepareInvokeFuncInput()
	if err != nil {
		return nil, err
	}
	client, err := util.NewFClient(gConfig)
	if err != nil {
		return nil, err
	}
	funcResponse, serviceError := client.InvokeFunction(invokeFuncInput)
	if serviceError != nil {
		return nil, serviceError
	}
	return funcResponse, nil
}

func prepareInvokeFuncInput() error {
	err := prepareCommon()
	if err != nil {
		return err
	}

	if eventStr != "" {
		invokeFuncInput.WithPayload([]byte(eventStr))
	} else {
		if eventFile != "" {
			bytes, err := ioutil.ReadFile(eventFile)
			if err != nil {
				return err
			}
			invokeFuncInput.WithPayload(bytes)
		}
	}
	invokeFuncInput.WithInvocationType(invocationType).
		WithHeader(HeaderInvocationCodeVersion, InvocationCodeVersionLatest)
	invokeFuncInput.ServiceName = &serviceName
	invokeFuncInput.FunctionName = &functionName
	invokeFuncInput.WithQualifier(qualifier)
	return nil
}

func init() {
	invokeFuncCmd.Flags().Bool("help", false, "Invoke function")
	invokeFuncCmd.Flags().StringVar(&invocationType, "invocation-type", "Sync", "invocation type")
	invokeFuncCmd.Flags().StringVarP(&serviceName, "service-name", "s", "", "service name")
	invokeFuncCmd.Flags().StringVarP(&functionName, "function-name", "f", "", "function name")
	invokeFuncCmd.Flags().StringVar(&eventStr, "event-str", "", "invoke event string")
	invokeFuncCmd.Flags().StringVar(&eventFile, "event-file", "", "invoke event in file with json format")
	invokeFuncCmd.Flags().StringVarP(&invocationOutputFile, "output", "o", "", "output filename")
	invokeFuncCmd.Flags().BoolVarP(&invkDebugEnabled, "debug", "d", false, "debug mode")
	invokeFuncCmd.Flags().StringVarP(&qualifier, "qualifier", "q", "", "service version or alias, optional")
}
