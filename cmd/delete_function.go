package cmd

import (
	"github.com/aliyun/fc-go-sdk"
	"github.com/aliyun/fcli/util"

	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	functionCmd.AddCommand(deleteFuncCmd)

	deleteFuncInput.serviceName = deleteFuncCmd.Flags().StringP("service-name", "s", "", "the service name")
	deleteFuncInput.functionName = deleteFuncCmd.Flags().StringP("function-name", "f", "", "the function name")
	deleteFuncInput.etag = deleteFuncCmd.Flags().String(
		"etag", "", "provide etag to do the conditional delete. "+
			"If the specified etag does not match the function's, the delete will fail.")
}

type deleteFuncInputType struct {
	serviceName  *string
	functionName *string
	etag         *string
}

var deleteFuncInput deleteFuncInputType

var deleteFuncCmd = &cobra.Command{
	Use:     "delete [option]",
	Aliases: []string{"d"},
	Short:   "Delete funtion",
	Long:    ``,

	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := util.NewFClient(gConfig)
		if err != nil {
			return fmt.Errorf("can not create fc client: %s\n", err)
		}

		input := fc.NewDeleteFunctionInput(*deleteFuncInput.serviceName, *deleteFuncInput.functionName)
		if cmd.Flags().Changed("etag") {
			input.WithIfMatch(*deleteFuncInput.etag)
		}

		_, err = client.DeleteFunction(input)
		if err != nil {
			return fmt.Errorf("%s\n", err)
		}
		return nil
	},
}
