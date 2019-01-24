package cmd

import (
	"github.com/aliyun/fc-go-sdk"
	"github.com/aliyun/fcli/util"

	"fmt"

	"encoding/json"

	"github.com/spf13/cobra"
)

func init() {
	functionCmd.AddCommand(listFuncCmd)

	listFuncCmd.Flags().Bool("help", false, "list functions")

	listFuncInput.serviceName = listFuncCmd.Flags().StringP(
		"service-name", "s", "", "list the functions belong to the specified service")
	listFuncInput.prefix = listFuncCmd.Flags().StringP(
		"prefix", "p", "", "list the functions whose names contain the specified prefix")
	listFuncInput.startKey = listFuncCmd.Flags().StringP(
		"start-key", "k", "", "start key is where you want to start listing from")
	listFuncInput.nextToken = listFuncCmd.Flags().StringP(
		"next-token", "t", "", "continue listing the functions from the previous point")
	listFuncInput.limit = listFuncCmd.Flags().Int32P(
		"limit", "l", 100, "the max number of the returned functions")
	listFuncInput.nameOnly = listFuncCmd.Flags().Bool(
		"name-only", true, "display function name only")
	listFuncInput.qualifier = listFuncCmd.Flags().StringP(
		"qualifier", "q", "", "service version or alias, optional")
}

type listFuncInputType struct {
	serviceName *string
	prefix      *string
	startKey    *string
	nextToken   *string
	limit       *int32
	nameOnly    *bool
	qualifier   *string
}

var listFuncInput listFuncInputType

var listFuncCmd = &cobra.Command{
	Use:     "list [option]",
	Aliases: []string{"l"},
	Short:   "List functions belong to the specified service",
	Long:    ``,

	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := util.NewFClient(gConfig)
		if err != nil {
			return fmt.Errorf("can not create fc client: %s\n", err)
		}

		input := fc.NewListFunctionsInput(*listFuncInput.serviceName).
			WithPrefix(*listFuncInput.prefix).
			WithStartKey(*listFuncInput.startKey).
			WithNextToken(*listFuncInput.nextToken).
			WithLimit(*listFuncInput.limit).
			WithQualifier(*listFuncInput.qualifier)

		resp, err := client.ListFunctions(input)
		if err != nil {
			return fmt.Errorf("%s\n", err)
		}

		if *listFuncInput.nameOnly {
			type listFuncOutputType struct {
				Functions []*string
				NextToken *string
			}
			output := listFuncOutputType{
				NextToken: resp.NextToken,
			}
			for _, f := range resp.Functions {
				output.Functions = append(output.Functions, f.FunctionName)
			}
			ret, _ := json.MarshalIndent(output, "", "  ")
			fmt.Printf("%s\n", string(ret))
		} else {
			ret, _ := json.MarshalIndent(resp, "", "  ")
			fmt.Printf("%s\n", string(ret))
		}
		return nil
	},
}
