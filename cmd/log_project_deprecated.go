package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(logProjectDepCmd)
}

var logProjectDepCmd = &cobra.Command{
	Use:     "sls_project",
	Aliases: []string{"log_project"},
	Short:   "SLS project related operations",
	Long:    ``,
	Run: func(cmd *cobra.Command, args []string) {

	},
	Hidden:     true,
	Deprecated: "\b\b.\n[WARNNING] sls_project may be removed in future.\n",
}

func init() {
	logProjectDepCmd.Flags().Bool("help", true, "Print Usage")
}
