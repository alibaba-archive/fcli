package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(functionCmd)
}

var functionName, serviceName string

var functionCmd = &cobra.Command{
	Use:     "function",
	Aliases: []string{"f"},
	Short:   "function related operation",
	Long:    ``,

	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() {
	functionCmd.Flags().Bool("help", true, "Print Usage")
}
