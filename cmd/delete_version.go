package cmd

import (
	"github.com/aliyun/fc-go-sdk"
	"github.com/aliyun/fcli/util"

	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	serviceVersionDepCmd.AddCommand(deleteVersionCmd)
	serviceVersionCmd.AddCommand(deleteVersionCmd)

	deleteVersionCmd.Flags().Bool("help", false, "delete service version")

	deleteServiceVersionInput.ServiceName = deleteVersionCmd.Flags().StringP(
		"service-name", "s", "", "the service name")
	deleteServiceVersionInput.VersionID = deleteVersionCmd.Flags().StringP(
		"version-id", "v", "", "version id")
}

var deleteServiceVersionInput fc.DeleteServiceVersionInput

var deleteVersionCmd = &cobra.Command{
	Use:     "delete [option]",
	Aliases: []string{"d"},
	Short:   "Delete service version",
	Long: `
delete service version
EXAMPLE:
fcli service version delete -s(--service-name)   service_name
				-v(--version-id) 1
				`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := util.NewFClient(gConfig)
		if err != nil {
			return fmt.Errorf("can not create fc client: %s\n", err)
		}
		_, err = client.DeleteServiceVersion(&deleteServiceVersionInput)
		if err != nil {
			return fmt.Errorf("%s\n", err)
		}
		return nil
	},
}
