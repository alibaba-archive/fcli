package cmd

import (
	"fmt"
	"time"

	"github.com/aliyun/fc-go-sdk"
	"github.com/aliyun/fcli/util"

	"github.com/aliyun/aliyun-log-go-sdk"
	"github.com/spf13/cobra"
)

func init() {
	functionCmd.AddCommand(functionLogsCmd)
}

// LogParam defines log cmd parameters
type LogParam struct {
	startTime *string
	endTime   *string
}

var logParams LogParam

var functionLogsCmd = &cobra.Command{
	Use:     "logs [option]",
	Aliases: []string{"l"},
	Short:   "fetch the logs of a function",
	Long: `
function logs
Example
monitor the function logs like tail
   fcli function logs -s "service name"
                     -f "function name"

fetch all the function logs within [start, now)
   fcli function logs -s "service name"
                     -f "function name"
                     --start "start time"
					 
fetch all the function logs within [start, end) 
   fcli function logs -s "service name"
                     -f "function name"
                     --start "start time"  --end "end time"

time format is UTC RFC3339, such as 2017-01-01T01:02:03Z
   		  `,

	Run: func(cmd *cobra.Command, args []string) {
		err := functionLogsRun(cmd)
		if err != nil {
			fmt.Println(err)
		}
	},
}

func functionLogsRun(cmd *cobra.Command) error {
	err := prepareCommon()
	if err != nil {
		return err
	}

	client, err = util.NewFClient(gConfig)
	if err != nil {
		return err
	}

	smeta, err := client.GetService(fc.NewGetServiceInput(serviceName))
	if err != nil {
		return err
	}

	enableLogging := false
	if smeta.LogConfig != nil {
		if !isEmpty(smeta.LogConfig.Project) && !isEmpty(smeta.LogConfig.Logstore) && !isEmpty(smeta.Role) {
			enableLogging = true
		}
	}
	if enableLogging {
		//endpoint := fmt.Sprintf(util.LogEndpointFmt, util.GetRegionNo(basicConfig.Endpoint))
		project := *smeta.LogConfig.Project
		logstore := *smeta.LogConfig.Logstore
		slsProject, err := sls.NewLogProject(project, gConfig.SLSEndpoint, gConfig.AccessKeyID, gConfig.AccessKeySecret)
		if err != nil {
			return err
		}
		slsLogstore, err := slsProject.GetLogStore(logstore)
		if err != nil {
			return err
		}

		// Get function logs like tail
		if !cmd.Flags().Changed("start") && !cmd.Flags().Changed("end") {
			err := tailFunctionLogs(slsLogstore, serviceName, functionName)
			if err != nil {
				return err
			}
		}
		// Get function logs within specific timerange
		if cmd.Flags().Changed("start") {
			startTime, err := time.Parse(util.TimeLayoutInLogs, *logParams.startTime)
			if err != nil {
				return fmt.Errorf("start time format error, expect:%s, actual:%s", util.TimeLayoutInLogs, *logParams.startTime)
			}
			startTimestamp := startTime.Unix()
			var endTimestamp int64
			if !cmd.Flags().Changed("end") {
				endTimestamp = time.Now().Unix()
			} else {
				endTime, err := time.Parse(util.TimeLayoutInLogs, *logParams.endTime)
				if err != nil {
					return fmt.Errorf("end time format error, expect:%s, actual:%s", util.TimeLayoutInLogs, *logParams.endTime)
				}
				endTimestamp = endTime.Unix()
			}
			err = util.GetAllLogsWithinTimeRange(slsLogstore, serviceName, functionName, startTimestamp, endTimestamp)
			if err != nil {
				return err
			}
		}

	} else {
		return fmt.Errorf("function logging was disabled, please update service to give valid service role/logConfig parameters")
	}
	return nil
}

func tailFunctionLogs(slsLogstore *sls.LogStore, serviceName, functionName string) error {
	now := time.Now().Unix()
	start := now - util.GetLogsIntervalTimeInSecs
	end := now
	for {
		err := util.GetAllLogsWithinTimeRange(slsLogstore, serviceName, functionName, start, end)
		if err != nil {
			return err
		}

		start = end
		sleepTimeInSec := util.GetLogsIntervalTimeInSecs - (time.Now().Unix() - start)
		time.Sleep(time.Duration(sleepTimeInSec) * time.Second)
		end = time.Now().Unix()
	}
}

func init() {
	functionLogsCmd.Flags().BoolP("help", "h", false, "Print Usage")
	functionLogsCmd.Flags().StringVarP(&serviceName, "service-name", "s", "", "service name")
	functionLogsCmd.Flags().StringVarP(&functionName, "function-name", "f", "", "function name")
	logParams.startTime = functionLogsCmd.Flags().String("start", "", "start time, such as 2017-01-01T01:02:03Z")
	logParams.endTime = functionLogsCmd.Flags().String("end", "", "end time,   such as 2017-01-01T02:02:03Z")
}