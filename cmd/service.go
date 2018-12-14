package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(serviceCmd)
}

var serviceCmd = &cobra.Command{
	Use:     "service",
	Aliases: []string{"s"},
	Short:   "service related operation",
	Long:    ``,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() {
	serviceCmd.Flags().Bool("help", true, "Print Usage")
}
