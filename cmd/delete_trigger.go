package cmd

import (
	"github.com/aliyun/fc-go-sdk"
	"github.com/aliyun/fcli/util"

	"github.com/spf13/cobra"
)

func init() {
	triggerCmd.AddCommand(deleteTriggerCmd)
}

var deleteTriggerEtag *string
var deleteTriggerCmd = &cobra.Command{
	Use:     "delete",
	Aliases: []string{"c"},
	Short:   "delete trigger",
	Long: `
delete trigger
EXAMPLE:

fcli trigger delete -s(--service-name)  service_name
				   -f(--function-name) function_name 
				   -t(--trigger-name)  trigger_name
				   --etag a198ec37e2a1c2ababbb3717074f29ea
				   
You can get etag from trigger get cmd, then delete with condition
			`,
	Run: func(cmd *cobra.Command, args []string) {
		prettyPrint(deleteTriggerRun(cmd))
	},
}

func deleteTriggerRun(cmd *cobra.Command) (*fc.DeleteTriggerOutput, error) {
	err := prepareCommon()
	if err != nil {
		return nil, err
	}
	client, err := util.NewFClient(gConfig)
	if err != nil {
		return nil, err
	}
	deleteTriggerInput := fc.NewDeleteTriggerInput(serviceName, functionName, triggerName)
	if cmd.Flags().Changed("etag") {
		deleteTriggerInput.WithIfMatch(*deleteTriggerEtag)
	}

	response, serviceError := client.DeleteTrigger(deleteTriggerInput)
	if serviceError != nil {
		return nil, serviceError
	}
	return response, nil
}

func init() {
	deleteTriggerCmd.Flags().Bool("help", false, "delete trigger")
	deleteTriggerEtag = deleteTriggerCmd.Flags().String("etag", "", "delete with etag, you can get etag from trigger get cmd")
}
