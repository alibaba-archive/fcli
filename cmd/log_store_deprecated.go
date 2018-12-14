package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(logStoreDepCmd)

	logStoreDepCmd.Flags().Bool("help", true, "Print Usage")
}

var logStoreDepCmd = &cobra.Command{
	Use:     "sls_store",
	Aliases: []string{"log_store"},
	Short:   "SLS store related operations",
	Long:    ``,
	Run: func(cmd *cobra.Command, args []string) {

	},
	Hidden:     true,
	Deprecated: "\b\b.\n[WARNNING] sls_store may be removed in future.\n",
}
