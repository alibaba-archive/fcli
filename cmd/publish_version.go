package cmd

import (
	"github.com/aliyun/fc-go-sdk"
	"github.com/aliyun/fcli/util"

	"fmt"

	"github.com/spf13/cobra"
)

var output *bool

func init() {
	serviceVersionDepCmd.AddCommand(publishVersionCmd)
	serviceVersionCmd.AddCommand(publishVersionCmd)

	publishVersionCmd.Flags().Bool("help", false, "publish service version")

	publishServiceVersionInput.ServiceName = publishVersionCmd.Flags().StringP(
		"service-name", "s", "", "the service name")
	publishServiceVersionInput.Description = publishVersionCmd.Flags().StringP(
		"description", "d", "", "version description")
	publishServiceVersionInput.IfMatch = publishVersionCmd.Flags().String(
		"etag", "", "provide etag to do the conditional publish. "+
			"If the specified etag does not match the service's, the publish will fail.")
	output = publishVersionCmd.Flags().Bool(
		"output", false, "print raw response body of API invoke.")
}

var publishServiceVersionInput fc.PublishServiceVersionInput

var publishVersionCmd = &cobra.Command{
	Use:     "publish [option]",
	Aliases: []string{"p"},
	Short:   "Publish service version",
	Long: `
publish service version
EXAMPLE:
fcli service version publish -s(--service-name)   service_name
				-d(--description) description
				--etag            a198ec37e2a1c2ababbb3717074f29ea
			`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := util.NewFClient(gConfig)
		if err != nil {
			fmt.Printf("Error: can not create fc client: %s\n", err)
			return
		}

		resp, err := client.PublishServiceVersion(&publishServiceVersionInput)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
		}

		if *output {
			fmt.Println(resp)
		}
	},
}
