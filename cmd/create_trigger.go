package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/aliyun/fc-go-sdk"
	"github.com/aliyun/fcli/util"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type createTriggerInputType struct {
	serviceName       *string
	functionName      *string
	triggerName       *string
	triggerType       *string
	sourceARN         *string
	invocationRole    *string
	qualifier         *string
	triggerConfigFile *string
}

// CreateTriggerCliOutput a wrapper around CreateTriggerOutput in order to add in more information
type CreateTriggerCliOutput struct {
	triggerCliOutputDecorate
	fc.CreateTriggerOutput
}

func (o CreateTriggerCliOutput) String() string {
	b, err := json.MarshalIndent(o, "", "  ")
	if err != nil {
		return ""
	}
	return string(b)
}

// MarshalJSON define how CreateTriggerCliOutput should be displayed in JSON format
func (o CreateTriggerCliOutput) MarshalJSON() ([]byte, error) {
	return json.Marshal(triggerCliOutputDisplay{
		HTTPTriggerURL:   o.HTTPTriggerURL,
		Header:           o.Header,
		TriggerName:      o.TriggerName,
		SourceARN:        o.SourceARN,
		TriggerType:      o.TriggerType,
		InvocationRole:   o.InvocationRole,
		Qualifier:        o.Qualifier,
		TriggerConfig:    o.TriggerConfig,
		CreatedTime:      o.CreatedTime,
		LastModifiedTime: o.LastModifiedTime,
	})
}

var createTriggerInput createTriggerInputType

func init() {
	triggerCmd.AddCommand(createTriggerCmd)

	createTriggerCmd.Flags().Bool(
		"help", false, "create trigger")

	createTriggerInput.serviceName = createTriggerCmd.Flags().StringP(
		"service-name", "s", "", "the service name")
	createTriggerInput.functionName = createTriggerCmd.Flags().StringP(
		"function-name", "f", "", "the function name")
	createTriggerInput.triggerName = createTriggerCmd.Flags().StringP(
		"trigger-name", "t", "", "the trigger name")
	createTriggerInput.triggerType = createTriggerCmd.Flags().String(
		"type", "", "trigger type, support oss, log, timer, http, cdn_events, mns_topic")
	createTriggerInput.sourceARN = createTriggerCmd.Flags().StringP(
		"source-arn", "a", "", "event source arn,for example, acs:oss:cn-hangzhou:123456:bucket1.timer trigger optional")
	createTriggerInput.invocationRole = createTriggerCmd.Flags().StringP(
		"role", "r", "", "invocation role,  timer trigger optional")
	createTriggerInput.triggerConfigFile = createTriggerCmd.Flags().StringP(
		"config", "c", "",
		`trigger config file, support json and yaml format.
		Below is a oss trigger config file with yaml format.

			OSSTriggerConfig:
    			    events:
        		        - oss:ObjectCreated:CopyObject
        		        - oss:ObjectCreated:PostObject
    			    filter:
        		        key:
            		            prefix: foo
            		            suffix: bar
}		`)
	createTriggerInput.qualifier = createTriggerCmd.Flags().StringP("qualifier", "q", "", "service version or alias, optional")
}

var createTriggerCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"c"},
	Short:   "create trigger",
	Long: `
create trigger
EXAMPLE:

fcli trigger create -s(service-name)  demo_service
				   -f(function-name) demo_function
				   -t(trigger-name)  demo_trigger
				   -type             oss
				   -a(source-arn)    acs:oss:cn-hangzhou:123456:bucket1
				   -r(role)          acs:ram::123456:role/role-v1
				   -c(config)        oss_trigger_sample.yaml
				   -q(qualifier)     LATEST

oss_trigger_sample.yaml example:
triggerConfig:
    events: 
        - oss:ObjectCreated:CopyObject
        - oss:ObjectCreated:PostObject
    filter:
        key:
            prefix: foo
            suffix: bar
			`,
	Run: func(cmd *cobra.Command, args []string) {
		prettyPrint(createTriggerRun(cmd))
	},
}

func createTriggerRun(cmd *cobra.Command) (*CreateTriggerCliOutput, error) {
	err := prepareCommon()
	if err != nil {
		return nil, err
	}
	client, err := util.NewFClient(gConfig)
	if err != nil {
		return nil, err
	}
	createTriggerInput, err := prepareCreateTriggerInput()
	if err != nil {
		return nil, err
	}
	response, serviceError := client.CreateTrigger(createTriggerInput)
	if serviceError != nil {
		return nil, serviceError
	}

	output := &CreateTriggerCliOutput{
		CreateTriggerOutput: *response,
	}
	decorateTriggerOutput(createTriggerInput.TriggerType, &output.triggerCliOutputDecorate)

	return output, nil
}

func prepareCreateTriggerInput() (*fc.CreateTriggerInput, error) {
	viper.SetConfigFile(*createTriggerInput.triggerConfigFile)
	err := viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s due to %v", *createTriggerInput.triggerConfigFile, err)
	}
	triggerConfig, err := util.GetTriggerConfig(*createTriggerInput.triggerType, *createTriggerInput.triggerConfigFile)
	if err != nil {
		return nil, err
	}

	input := fc.NewCreateTriggerInput(*createTriggerInput.serviceName, *createTriggerInput.functionName).
		WithTriggerName(*createTriggerInput.triggerName).
		WithTriggerType(*createTriggerInput.triggerType).
		WithTriggerConfig(triggerConfig)

	if *createTriggerInput.sourceARN != "" {
		input.WithSourceARN(*createTriggerInput.sourceARN)
	}

	if *createTriggerInput.invocationRole != "" {
		input.WithInvocationRole(*createTriggerInput.invocationRole)
	}

	if *createTriggerInput.qualifier != "" {
		input.WithQualifier(*createTriggerInput.qualifier)
	}

	return input, nil
}

func init() {
}
