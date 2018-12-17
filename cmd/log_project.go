package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	logCmd.AddCommand(logProjectCmd)
}

var logProjectCmd = &cobra.Command{
	Use:     "project",
	Aliases: []string{"project"},
	Short:   "SLS project related operations",
	Long:    ``,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() {
	logProjectCmd.Flags().Bool("help", true, "Print Usage")
}
