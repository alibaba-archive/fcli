package cmd

import (
	"github.com/aliyun/fc-go-sdk"
	"github.com/aliyun/fcli/util"

	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	serviceCmd.AddCommand(getServiceCmd)

	getServiceCmd.Flags().Bool("help", false, "Print Usage")
	getServiceInput.ServiceName = getServiceCmd.Flags().StringP("service-name", "s", "", "the service name")
	getServiceInput.Qualifier = getServiceCmd.Flags().StringP("qualifier", "q", "", "service version or alias, optional")
}

var getServiceInput fc.GetServiceInput

var getServiceCmd = &cobra.Command{
	Use:     "get [option]",
	Aliases: []string{"g"},
	Short:   "Get the information of service",
	Long:    ``,

	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := util.NewFClient(gConfig)
		if err != nil {
			return fmt.Errorf("can not create fc client: %s\n", err)
		}
		resp, err := client.GetService(&getServiceInput)
		if err == nil {
			fmt.Printf("%s\n", resp)
		} else {
			return fmt.Errorf("%s\n", err)
		}
		return nil
	},
}
