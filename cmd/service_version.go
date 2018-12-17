package cmd

import "github.com/spf13/cobra"

func init() {
	serviceCmd.AddCommand(serviceVersionCmd)
}

var serviceVersionCmd = &cobra.Command{
	Use:     "version",
	Aliases: []string{"version"},
	Short:   "service version related operation",
	Long:    ``,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() {
	serviceVersionCmd.Flags().Bool("help", true, "Print Usage")
}
