package cmd

import (
	"github.com/aliyun/fc-go-sdk"
	"github.com/aliyun/fcli/util"

	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	aliasCmd.AddCommand(getAliasCmd)

	getAliasCmd.Flags().Bool("help", false, "get alias")

	getAliasInput.ServiceName = getAliasCmd.Flags().StringP(
		"service-name", "s", "", "the service name")
	getAliasInput.AliasName = getAliasCmd.Flags().StringP(
		"alias-name", "a", "", "the alias name")
}

var getAliasInput fc.GetAliasInput

var getAliasCmd = &cobra.Command{
	Use:     "get [option]",
	Aliases: []string{"g"},
	Short:   "Get alias",
	Long: `
get alias
EXAMPLE:
fcli alias get -s(--service-name)        service_name
			-a(--alias-name) alias_name
			`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := util.NewFClient(gConfig)
		if err != nil {
			return fmt.Errorf("can not create fc client: %s\n", err)
		}
		resp, err := client.GetAlias(&getAliasInput)
		if err != nil {
			return fmt.Errorf("%s\n", err)
		} else {
			fmt.Println(resp)
		}
		return nil
	},
}
