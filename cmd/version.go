package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/aliyun/fcli/version"
)

func init() {
	RootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:     "version",
	Aliases: []string{"v"},
	Short:   "fcli version information",
	Long:    `fcli version information`,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("fcli version: %s\n", version.Version)
	},
}
