package cmd

import (
	"github.com/aliyun/fc-go-sdk"
	"github.com/aliyun/fcli/util"

	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	serviceCmd.AddCommand(deleteServiceCmd)

	deleteServiceCmd.Flags().Bool(
		"help", false, "delete service")
	deleteServiceInput.serviceName = deleteServiceCmd.Flags().StringP(
		"service-name", "s", "", "the service name")
	deleteServiceInput.etag = deleteServiceCmd.Flags().String(
		"etag", "", "provide etag to do the conditional delete. "+
			"If the specified etag does not match the service's, the delete will fail.")
}

type deleteServiceInputType struct {
	serviceName *string
	etag        *string
}

var deleteServiceInput deleteServiceInputType

var deleteServiceCmd = &cobra.Command{
	Use:     "delete [option]",
	Aliases: []string{"d"},
	Short:   "Delete service",
	Long:    ``,

	Run: func(cmd *cobra.Command, args []string) {
		client, err := util.NewFClient(gConfig)
		if err != nil {
			fmt.Printf("Error: can not create fc client: %s\n", err)
			return
		}

		input := fc.NewDeleteServiceInput(*deleteServiceInput.serviceName)
		if cmd.Flags().Changed("etag") {
			input.WithIfMatch(*deleteServiceInput.etag)
		}

		_, err = client.DeleteService(input)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
		}
	},
}
