package cmd

import (
	"github.com/aliyun/fc-go-sdk"
	"github.com/aliyun/fcli/util"

	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	aliasCmd.AddCommand(createAliasCmd)

	createAliasCmd.Flags().Bool("help", false, "create alias")

	createAliasInput.ServiceName = createAliasCmd.Flags().StringP(
		"service-name", "s", "", "the service name")
	createAliasInput.AliasName = createAliasCmd.Flags().StringP(
		"alias-name", "a", "", "the alias name")
	createAliasInput.VersionID = createAliasCmd.Flags().StringP(
		"version-id", "v", "", "the version id you want the alias pointing to")
	createAliasInput.Description = createAliasCmd.Flags().StringP(
		"description", "d", "", "version description, optional")
	createRoutes = createAliasCmd.Flags().StringArrayP(
		"route", "r", []string{}, "additional version weight for dark launch purpose, optional")
}

var createRoutes = new([]string)

var createAliasInput fc.CreateAliasInput

var createAliasCmd = &cobra.Command{
	Use:     "create [option]",
	Aliases: []string{"c"},
	Short:   "Create alias",
	Long: `
create alias
EXAMPLE:
fcli alias create -s(--service-name)      service_name
			-a(--alias-name)  alias_name
			-v(--version-id)  1
			-d(--description) description
			-r(--route)       2=0.05
			`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := util.NewFClient(gConfig)
		if err != nil {
			fmt.Printf("Error: can not create fc client: %s\n", err)
			return
		}
		createAliasInput.AdditionalVersionWeight = util.ParseAdditionalVersionWeight(*createRoutes)
		_, err = client.CreateAlias(&createAliasInput)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
		}
	},
}
