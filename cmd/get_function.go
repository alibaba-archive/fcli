package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/aliyun/fc-go-sdk"
	"github.com/aliyun/fcli/util"
)

func init() {
	functionCmd.AddCommand(getFuncCmd)

	getFuncCmd.Flags().Bool("help", false, "Print Usage")
	getFuncInput.ServiceName = getFuncCmd.Flags().StringP("service-name", "s", "", "the service name")
	getFuncInput.FunctionName = getFuncCmd.Flags().StringP("function-name", "f", "", "the function name")
	getFuncInput.Qualifier = getFuncCmd.Flags().StringP("qualifier", "q", "", "service version or alias, optional")
}

var getFuncInput fc.GetFunctionInput

var getFuncCmd = &cobra.Command{
	Use:     "get [option]",
	Aliases: []string{"g"},
	Short:   "Get function",
	Long:    ``,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := util.NewFClient(gConfig)
		if err != nil {
			fmt.Printf("Error: can not create fc client: %s\n", err)
			return
		}
		resp, err := client.GetFunction(&getFuncInput)
		if err == nil {
			fmt.Printf("%s\n", resp)
		} else {
			fmt.Printf("Error: %s\n", err)
		}
	},
}
