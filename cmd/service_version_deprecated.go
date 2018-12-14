package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(serviceVersionDepCmd)
}

var serviceVersionDepCmd = &cobra.Command{
	Use:     "service_version",
	Aliases: []string{"sv"},
	Short:   "service version related operation",
	Long:    ``,
	Run: func(cmd *cobra.Command, args []string) {

	},
	Hidden:     true,
	Deprecated: "\b\b.\n[WARNNING] service_version may be removed in future.\n",
}

func init() {
	serviceVersionDepCmd.Flags().Bool("help", true, "Print Usage")
}
