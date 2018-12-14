package cmd

import (
	"fmt"

	"github.com/aliyun/aliyun-log-go-sdk"
	"github.com/spf13/cobra"
)

func init() {
	logProjectDepCmd.AddCommand(listLogProjectCmd)
	logProjectCmd.AddCommand(listLogProjectCmd)

	listLogProjectCmd.Flags().Bool("help", false, "list log projects")
}

var listLogProjectCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l"},
	Short:   "List log projects belong to the configured account",
	Long:    ``,
	Run: func(cmd *cobra.Command, args []string) {
		err := getListLogProject(cmd)
		if err != nil {
			fmt.Println(err)
		}
	},
}

func getListLogProject(cmd *cobra.Command) error {
	slsClient := sls.CreateNormalInterface(
		gConfig.SLSEndpoint,
		gConfig.AccessKeyID,
		gConfig.AccessKeySecret,
		gConfig.SecurityToken,
	)

	projectNameList, err := slsClient.ListProject()
	if err != nil {
		return err
	}

	if len(projectNameList) == 0 {
		fmt.Printf("No project found.\n")
		fmt.Printf("Please make sure you have the correct SLS Endpoint setup in ~/.fcli/config.yml\n")
		fmt.Printf("For more information: https://www.alibabacloud.com/help/doc-detail/29008.htm\n")
		return nil
	}

	for _, projectName := range projectNameList {
		fmt.Printf("%s\n", projectName)
	}

	return nil
}
