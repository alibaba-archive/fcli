package cmd

import (
	"fmt"
	"github.com/aliyun/fc-go-sdk"
	"github.com/aliyun/fcli/util"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	serviceVersionDepCmd.AddCommand(listVersionsCmd)
	serviceVersionCmd.AddCommand(listVersionsCmd)

	listVersionsCmd.Flags().Bool("help", false, "list service versions")

	listServiceVersionsInput.ServiceName = listVersionsCmd.Flags().StringP(
		"service-name", "s", "", "the service name")
	listServiceVersionsInput.StartKey = listVersionsCmd.Flags().StringP(
		"start-key", "k", "", "start key is where you want to start listing from, optional")
	listServiceVersionsInput.NextToken = listVersionsCmd.Flags().StringP(
		"next-token", "t", "", "continue listing the versions from the previous point, optional")
	listServiceVersionsInput.Limit = listVersionsCmd.Flags().Int32P(
		"limit", "l", 100, "the max number of the returned versions, optional")
	listServiceVersionsInput.Direction = listVersionsCmd.Flags().StringP(
		"direction", "d", "BACKWARD", "listing direction, BACKWARD or FORWARD, optional")
}

var listServiceVersionsInput fc.ListServiceVersionsInput

var listVersionsCmd = &cobra.Command{
	Use:     "list [option]",
	Aliases: []string{"l"},
	Short:   "List service versions",
	Long: `
list service version
EXAMPLE:
fcli service version list -s(--service-name)      service_name
				 -k(--start-key)  "list alias start key"
			     	 -n(--next-token) "next token"
				 -l(--limit)      100
				 -d(--direction)  BACKWARD
				`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := util.NewFClient(gConfig)
		if err != nil {
			fmt.Printf("Error: can not create fc client: %s\n", err)
			return
		}
		*listServiceVersionsInput.Direction = strings.ToUpper(*listServiceVersionsInput.Direction)
		resp, err := client.ListServiceVersions(&listServiceVersionsInput)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
		} else {
			fmt.Println(resp)
		}
	},
}
