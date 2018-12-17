package cmd

import (
	"github.com/aliyun/fc-go-sdk"
	"github.com/aliyun/fcli/util"

	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	aliasCmd.AddCommand(listAliasesCmd)

	listAliasesCmd.Flags().Bool("help", false, "list aliases")

	listAliasesInput.ServiceName = listAliasesCmd.Flags().StringP(
		"service-name", "s", "", "the service name")
	listAliasesInput.Prefix = listAliasesCmd.Flags().StringP(
		"prefix", "p", "", "list the aliases whose names contain the specified prefix, optional")
	listAliasesInput.StartKey = listAliasesCmd.Flags().StringP(
		"start-key", "k", "", "start key is where you want to start listing from, optional")
	listAliasesInput.NextToken = listAliasesCmd.Flags().StringP(
		"next-token", "t", "", "continue listing the aliases from the previous point, optional")
	listAliasesInput.Limit = listAliasesCmd.Flags().Int32P(
		"limit", "l", 100, "the max number of the returned aliases, optional")
}

var listAliasesInput fc.ListAliasesInput

var listAliasesCmd = &cobra.Command{
	Use:     "list [option]",
	Aliases: []string{"l"},
	Short:   "List alias",
	Long: `
list alias
EXAMPLE:
fcli alias list -s(--service-name)       service_name
			-p(--prefix)     "alias prefix"
			-k(--start-key)  "list alias start key"
			-n(--next-token) "next token"
			-l(--limit)      100
			`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := util.NewFClient(gConfig)
		if err != nil {
			fmt.Printf("Error: can not create fc client: %s\n", err)
			return
		}
		resp, err := client.ListAliases(&listAliasesInput)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
		} else {
			fmt.Println(resp)
		}
	},
}
