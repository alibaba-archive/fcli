package cmd

import (
	"fmt"

	"github.com/aliyun/aliyun-log-go-sdk"
	"github.com/spf13/cobra"
)

func init() {
	logStoreDepCmd.AddCommand(listLogStoreCmd)
	logStoreCmd.AddCommand(listLogStoreCmd)

	listLogStoreCmd.Flags().Bool("help", false, "list log stores")

	listLogStoreInput.logProjectName = listLogStoreCmd.Flags().StringP(
		"project-name",
		"p",
		"",
		"list the Stores belong to the specified Project",
	)
}

var listLogStoreCmd = &cobra.Command{
	Use:     "list [option]",
	Aliases: []string{"l"},
	Short:   "SLS store related operations",
	Long:    ``,
	Run: func(cmd *cobra.Command, args []string) {

		// Check for valid ProjectName because SLS SDK error is not very good
		if *listLogStoreInput.logProjectName == "" {
			fmt.Printf("Error: Project name is required but not provided\n")
			return
		}

		slsProject, err := sls.NewLogProject(
			*listLogStoreInput.logProjectName,
			gConfig.SLSEndpoint,
			gConfig.AccessKeyID,
			gConfig.AccessKeySecret,
		)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			return
		}

		// To use Security Token if available in global config
		if gConfig.SecurityToken != "" {
			fmt.Printf("Using SecurityToken instead of AccessKey.\n")
			slsProject.WithToken(gConfig.SecurityToken)
		}

		// Call Api to get List of names
		storeNameList, err := slsProject.ListLogStore()
		if err != nil {
			fmt.Printf("Error getting list of Stores: %s\n", err)
			return
		}

		// Print out list of names
		for _, storeName := range storeNameList {
			fmt.Printf("%s\n", storeName)
		}
	},
}

type listLogStoreInputType struct {
	logProjectName *string
}

var listLogStoreInput listLogStoreInputType
