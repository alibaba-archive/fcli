package cmd

import (
	"github.com/aliyun/fc-go-sdk"
	"github.com/aliyun/fcli/util"

	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	aliasCmd.AddCommand(updateAliasCmd)

	updateAliasCmd.Flags().Bool("help", false, "update alias")

	updateAliasInput.ServiceName = updateAliasCmd.Flags().StringP(
		"service-name", "s", "", "the service name")
	updateAliasInput.AliasName = updateAliasCmd.Flags().StringP(
		"alias-name", "a", "", "the alias name")
	updateAliasInput.VersionID = updateAliasCmd.Flags().StringP(
		"version-id", "v", "", "the version id you want the alias pointing to")
	updateAliasInput.Description = updateAliasCmd.Flags().StringP(
		"description", "d", "", "version description, optional")
	updateRoutes = updateAliasCmd.Flags().StringArrayP(
		"route", "r", []string{}, "additional version weight for dark launch purpose, optional")
	updateAliasInput.IfMatch = updateAliasCmd.Flags().String(
		"etag", "", "provide etag to do the conditional update. "+
			"If the specified etag does not match the alias's, the update will fail.")
}

var updateRoutes = new([]string)

var updateAliasInput fc.UpdateAliasInput

var updateAliasCmd = &cobra.Command{
	Use:     "update [option]",
	Aliases: []string{"u"},
	Short:   "Update alias",
	Long: `
update alias
EXAMPLE:
fcli alias update -s(--service-name)      service_name
			-a(--alias-name) alias_name
			-v(--version-id)  1
			-d(--description) description
			-r(--route)       2=0.05
			--etag            a198ec37e2a1c2ababbb3717074f29ea
			`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := util.NewFClient(gConfig)
		if err != nil {
			return fmt.Errorf("can not create fc client: %s\n", err)
		}
		updateAliasInput.AdditionalVersionWeight = util.ParseAdditionalVersionWeight(*updateRoutes)
		_, err = client.UpdateAlias(&updateAliasInput)
		if err != nil {
			return fmt.Errorf("%s\n", err)
		}
		return nil
	},
}
