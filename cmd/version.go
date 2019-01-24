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

	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("fcli version: %s\n", version.Version)
		return nil
	},
}
