package cmd

import (
	"encoding/json"
	"github.com/aliyun/fc-go-sdk"
	"github.com/aliyun/fcli/util"
	"github.com/spf13/cobra"
	"net/http"
)

type triggerCliOutputDisplay struct {
	Header           http.Header
	HTTPTriggerURL   *string     `json:"httpUrl,omitempty"`
	TriggerName      *string     `json:"triggerName"`
	SourceARN        *string     `json:"sourceArn"`
	TriggerType      *string     `json:"triggerType"`
	InvocationRole   *string     `json:"invocationRole"`
	Qualifier        *string     `json:"qualifier"`
	TriggerConfig    interface{} `json:"triggerConfig"`
	CreatedTime      *string     `json:"createdTime"`
	LastModifiedTime *string     `json:"lastModifiedTime"`
}

// GetTriggerCliOutput is an envelope struct to decorate fc api response with additional information
type GetTriggerCliOutput struct {
	triggerCliOutputDecorate
	fc.GetTriggerOutput
}

func (o GetTriggerCliOutput) String() string {
	b, err := json.MarshalIndent(o, "", "  ")
	if err != nil {
		return ""
	}
	return string(b)
}

// MarshalJSON defines how GetTriggerCliOutput should be displayed in JSON format
func (o GetTriggerCliOutput) MarshalJSON() ([]byte, error) {
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

func init() {
	triggerCmd.AddCommand(getTriggerCmd)
}

var getTriggerCmd = &cobra.Command{
	Use:     "get",
	Aliases: []string{"c"},
	Short:   "get trigger",
	Long: `
get trigger
EXAMPLE:
fcli trigger get -s(--service-name)  service_name
				-f(--function-name) function_name 
				-t(--trigger-name)  trigger_name
			`,
	Run: func(cmd *cobra.Command, args []string) {
		prettyPrint(getTriggerRun(cmd))
	},
}

func getTriggerRun(cmd *cobra.Command) (*GetTriggerCliOutput, error) {
	err := prepareCommon()
	if err != nil {
		return nil, err
	}
	client, err := util.NewFClient(gConfig)
	if err != nil {
		return nil, err
	}
	getTriggerInput := fc.NewGetTriggerInput(serviceName, functionName, triggerName)
	response, serviceError := client.GetTrigger(getTriggerInput)
	if serviceError != nil {
		return nil, serviceError
	}
	output := &GetTriggerCliOutput{GetTriggerOutput: *response}
	decorateTriggerOutput(response.TriggerType, &output.triggerCliOutputDecorate)

	return output, nil
}

func init() {
	getTriggerCmd.Flags().Bool("help", false, "create trigger")
}
