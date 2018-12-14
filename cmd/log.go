package cmd

import "github.com/spf13/cobra"

func init() {
	RootCmd.AddCommand(logCmd)
}

var logCmd = &cobra.Command{
	Use:     "sls",
	Aliases: []string{"log"},
	Short:   "sls related operation",
	Long:    ``,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() {
	logCmd.Flags().Bool("help", true, "Print Usage")
}
