package cmd

import (
	"github.com/spf13/cobra"

	"github.com/aliyun/fc-go-sdk"
	"github.com/aliyun/fcli/util"
)

func init() {
	triggerCmd.AddCommand(listTriggerCmd)
}

// QueryParam query trigger parameters
type QueryParam struct {
	prefix    *string
	startKey  *string
	nextToken *string
	limit     *int32
}

var triggerQueryParam QueryParam
var isShowAll bool
var onlyNames bool

var listTriggerCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"c"},
	Short:   "list trigger",
	Long: `
list trigger
EXAMPLE:

fcli trigger list -s(--service-name)  service_name
				 -f(--function-name) function_name 
				 -p(--prefix) "trigger prefix"
				 -k(--start-key) "list trigger start key"
			     -n(--next-token) "next token"
				 -l(--limit) 100
				 --all false
				 --only-names false
			`,
	Run: func(cmd *cobra.Command, args []string) {
		prettyPrint(listTriggerRun(cmd))
	},
}

func listTriggerRun(cmd *cobra.Command) (interface{}, error) {
	err := prepareCommon()
	if err != nil {
		return nil, err
	}
	client, err := util.NewFClient(gConfig)
	if err != nil {
		return nil, err
	}
	listTriggerInput := fc.NewListTriggersInput(serviceName, functionName).
		WithLimit(*triggerQueryParam.limit).
		WithNextToken(*triggerQueryParam.nextToken).
		WithPrefix(*triggerQueryParam.prefix).
		WithStartKey(*triggerQueryParam.startKey)

	// return all trigger list
	if isShowAll {
		triggerList := []fc.ListTriggersOutput{}
		triggerNameList := []string{}
		for {
			listTriggerResponse, serviceError := client.ListTriggers(listTriggerInput)
			if serviceError != nil {
				return nil, serviceError
			}
			triggerList = append(triggerList, *listTriggerResponse)
			for _, v := range listTriggerResponse.Triggers {
				triggerNameList = append(triggerNameList, *v.TriggerName)
			}
			if listTriggerResponse.NextToken == nil {
				break
			}
			listTriggerInput.NextToken = listTriggerResponse.NextToken
		}
		if onlyNames {
			return triggerNameList, nil
		}
		return triggerList, nil
	}

	// return trigger list query with limit
	listTriggerResp, err := client.ListTriggers(listTriggerInput)
	if err != nil {
		return nil, err
	}
	if onlyNames {
		triggerNameList := []string{}
		for _, v := range listTriggerResp.Triggers {
			triggerNameList = append(triggerNameList, *v.TriggerName)
		}
		return triggerNameList, nil
	}
	return listTriggerResp, nil
}

func init() {
	listTriggerCmd.Flags().Bool("help", false, "list trigger")
	triggerQueryParam.prefix = listTriggerCmd.Flags().StringP("prefix", "p", "", "trigger prefix")
	triggerQueryParam.startKey = listTriggerCmd.Flags().StringP("start-key", "k", "", "trigger start key")
	triggerQueryParam.nextToken = listTriggerCmd.Flags().StringP("next-token", "n", "", "list trigger next token")
	triggerQueryParam.limit = listTriggerCmd.Flags().Int32P("limit", "l", 100, "limit number")
	listTriggerCmd.Flags().BoolVar(&isShowAll, "all", false, "get all the trigger list")
	listTriggerCmd.Flags().BoolVar(&onlyNames, "only-names", false, "get all the trigger list but only show names")
}
