package cmd

import (
	"github.com/aliyun/fc-go-sdk"
	"github.com/aliyun/fcli/util"

	"fmt"

	"github.com/spf13/cobra"
)

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
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := util.NewFClient(gConfig)
		if err != nil {
			return fmt.Errorf("can not create fc client: %s\n", err)
		}
		_, err = client.PublishServiceVersion(&publishServiceVersionInput)
		if err != nil {
			return fmt.Errorf("%s\n", err)
		}
		return nil
	},
}
