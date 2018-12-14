package cmd

import (
	"github.com/aliyun/fc-go-sdk"
	"github.com/aliyun/fcli/util"

	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	serviceCmd.AddCommand(listServiceCmd)

	listServiceCmd.Flags().Bool("help", false, "list functions")

	listServiceInput.prefix = listServiceCmd.Flags().StringP(
		"prefix", "p", "", "list the services whose names contain the specified prefix")
	listServiceInput.startKey = listServiceCmd.Flags().StringP(
		"start-key", "k", "", "start key is where you want to start listing from")
	listServiceInput.nextToken = listServiceCmd.Flags().StringP(
		"next-token", "t", "", "continue listing the functions from the previous point")
	listServiceInput.limit = listServiceCmd.Flags().Int32P(
		"limit", "l", 100, "the max number of the returned services")
	listServiceInput.nameOnly = listServiceCmd.Flags().Bool(
		"name-only", true, "display service name only")
}

type listServiceInputType struct {
	prefix    *string
	startKey  *string
	nextToken *string
	limit     *int32
	nameOnly  *bool
}

var listServiceInput listServiceInputType

var listServiceCmd = &cobra.Command{
	Use:     "list [option]",
	Aliases: []string{"l"},
	Short:   "List services of the current account",
	Long:    ``,

	Run: func(cmd *cobra.Command, args []string) {
		client, err := util.NewFClient(gConfig)
		if err != nil {
			fmt.Printf("Error: can not create fc client: %s\n", err)
			return
		}

		input := fc.NewListServicesInput().
			WithPrefix(*listServiceInput.prefix).
			WithStartKey(*listServiceInput.startKey).
			WithNextToken(*listServiceInput.nextToken).
			WithLimit(*listServiceInput.limit)

		resp, err := client.ListServices(input)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			return
		}

		if *listServiceInput.nameOnly {
			type listServiceOutputType struct {
				Services  []*string
				NextToken *string
			}
			output := listServiceOutputType{
				NextToken: resp.NextToken,
			}
			for _, f := range resp.Services {
				output.Services = append(output.Services, f.ServiceName)
			}
			ret, _ := json.MarshalIndent(output, "", "  ")
			fmt.Printf("%s\n", string(ret))
		} else {
			ret, _ := json.MarshalIndent(resp, "", "  ")
			fmt.Printf("%s\n", string(ret))
		}
	},
}
