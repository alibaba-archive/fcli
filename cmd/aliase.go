package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(aliasCmd)
}

var aliasCmd = &cobra.Command{
	Use:     "alias",
	Aliases: []string{"a"},
	Short:   "alias related operation",
	Long:    ``,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() {
	aliasCmd.Flags().Bool("help", true, "Print Usage")
}
