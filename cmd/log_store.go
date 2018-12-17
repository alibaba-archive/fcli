package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	logCmd.AddCommand(logStoreCmd)

	logStoreCmd.Flags().Bool("help", true, "Print Usage")
}

var logStoreCmd = &cobra.Command{
	Use:     "store",
	Aliases: []string{"store"},
	Short:   "SLS store related operations",
	Long:    ``,
	Run: func(cmd *cobra.Command, args []string) {

	},
}
