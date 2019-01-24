package cmd

import (
	"github.com/aliyun/fc-go-sdk"
	"github.com/aliyun/fcli/util"

	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	aliasCmd.AddCommand(deleteAliasCmd)

	deleteAliasCmd.Flags().Bool("help", false, "delete alias")

	deleteAliasInput.ServiceName = deleteAliasCmd.Flags().StringP(
		"service-name", "s", "", "the service name")
	deleteAliasInput.AliasName = deleteAliasCmd.Flags().StringP(
		"alias-name", "a", "", "the alias name")
	deleteAliasInput.IfMatch = deleteAliasCmd.Flags().String(
		"etag", "", "provide etag to do the conditional delete. "+
			"If the specified etag does not match the alias's, the delete will fail.")
}

var deleteAliasInput fc.DeleteAliasInput

var deleteAliasCmd = &cobra.Command{
	Use:     "delete [option]",
	Aliases: []string{"d"},
	Short:   "Delete alias",
	Long: `
delete alias
EXAMPLE:
fcli alias delete -s(--service-name)      service_name
			-a(--alias-name)  alias_name
			--etag            a198ec37e2a1c2ababbb3717074f29ea
			`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := util.NewFClient(gConfig)
		if err != nil {
			return fmt.Errorf("can not create fc client: %s\n", err)
		}
		_, err = client.DeleteAlias(&deleteAliasInput)
		if err != nil {
			return fmt.Errorf("%s\n", err)
		}
		return nil
	},
}
