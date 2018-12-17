package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(triggerCmd)
}

var triggerName string

var triggerCmd = &cobra.Command{
	Use:     "trigger",
	Aliases: []string{"t"},
	Short:   "trigger related operation",
	Long:    ``,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() {
	triggerCmd.Flags().BoolP("help", "h", true, "Print Usage")
	triggerCmd.PersistentFlags().StringVarP(&serviceName, "service-name", "s", "", "service name")
	triggerCmd.PersistentFlags().StringVarP(&functionName, "function-name", "f", "", "function name")
	triggerCmd.PersistentFlags().StringVarP(&triggerName, "trigger-name", "t", "", "trigger name")
}

func getTriggerHelp(operation, value string) string {
	return `
ERROR:
   service name or function name or trigger name is null
EXAMPLE:
   ` + operation + ` -s(--service-name) "service_name" -f(--function-name) "function_name" -t(--trigger-name) "trigger_name"` + value + `
HELP: 
   ` + operation + ` --help
`
}

func printTrigger() {
	fmt.Println("trigger  name:" + triggerName)
	fmt.Println("function name:" + functionName)
	fmt.Println("service  name:" + serviceName)
}
