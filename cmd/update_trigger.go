package cmd

import (
	"encoding/json"
	"github.com/aliyun/fc-go-sdk"
	"github.com/aliyun/fcli/util"
	"github.com/spf13/cobra"
)

// TriggerUpdateParam ...
type TriggerUpdateParam struct {
	Etag              *string
	InvocationRole    *string
	Qualifier         *string
	TriggerConfigFile *string
}

// UpdateTriggerCliOutput a wrapper around UpdateTriggerOutput in order to add in more information
type UpdateTriggerCliOutput struct {
	triggerCliOutputDecorate
	fc.UpdateTriggerOutput
}

func (o UpdateTriggerCliOutput) String() string {
	b, err := json.MarshalIndent(o, "", "  ")
	if err != nil {
		return ""
	}
	return string(b)
}

// MarshalJSON defines how UpdateTriggerCliOutput should be displayed in JSON format
func (o UpdateTriggerCliOutput) MarshalJSON() ([]byte, error) {
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

var updateParam TriggerUpdateParam

func init() {
	triggerCmd.AddCommand(updateTriggerCmd)
}

var updateTriggerCmd = &cobra.Command{
	Use:     "update",
	Aliases: []string{"c"},
	Short:   "update trigger",
	Long: `
update trigger
EXAMPLE:

fcli trigger update -s(service-name)  demo_service
				   -f(function-name) demo_function
				   -t(trigger-name)  demo_trigger
				   --etag            a198ec37e2a1c2ababbb3717074f29ea      
				   --invocation-role acs:ram::123456:role/role-v1
				   --trigger-config  oss_trigger_sample.yaml
				   --q(qualifier)    LATEST

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
		prettyPrint(updateTriggerRun(cmd))
	},
}

func updateTriggerRun(cmd *cobra.Command) (*UpdateTriggerCliOutput, error) {
	err := prepareCommon()
	if err != nil {
		return nil, err
	}
	client, err := util.NewFClient(gConfig)
	if err != nil {
		return nil, err
	}
	updateTriggerInput := fc.NewUpdateTriggerInput(serviceName, functionName, triggerName)
	if cmd.Flags().Changed("etag") {
		updateTriggerInput.WithIfMatch(*updateParam.Etag)
	}
	if cmd.Flags().Changed("invocation-role") {
		updateTriggerInput.WithInvocationRole(*updateParam.InvocationRole)
	}

	triggerType := ""
	if cmd.Flags().Changed("trigger-config") {

		getTriggerOutput, err := client.GetTrigger(fc.NewGetTriggerInput(serviceName, functionName, triggerName))
		if err != nil {
			return nil, err
		}
		triggerType = *getTriggerOutput.TriggerType

		triggerConfig, err := util.GetTriggerConfig(triggerType, *updateParam.TriggerConfigFile)
		if err != nil {
			return nil, err
		}
		updateTriggerInput.WithTriggerConfig(triggerConfig)

	}
	if cmd.Flags().Changed("qualifier") {
		if updateParam.Qualifier != nil && *updateParam.Qualifier != "" {
			updateTriggerInput.WithQualifier(*updateParam.Qualifier)
		}
	}

	resp, err := client.UpdateTrigger(updateTriggerInput)
	if err != nil {
		return nil, err
	}
	output := &UpdateTriggerCliOutput{
		UpdateTriggerOutput: *resp,
	}
	decorateTriggerOutput(&triggerType, &output.triggerCliOutputDecorate)

	return output, err
}

func init() {
	updateTriggerCmd.Flags().Bool("help", false, "update trigger")
	updateParam.Etag = updateTriggerCmd.Flags().String("etag", "", "update with etag, you can get etag from GetTrigger call")
	updateParam.InvocationRole = updateTriggerCmd.Flags().String("invocation-role", "", "trigger invocation role")
	updateParam.TriggerConfigFile = updateTriggerCmd.Flags().String("trigger-config", "", "trigger config file")
	updateParam.Qualifier = updateTriggerCmd.Flags().StringP("qualifier", "q", "", "service version or alias, optional")
}
