package cmd

import (
	"strings"

	"github.com/aliyun/fcli/ram"

	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	//RootCmd.AddCommand(roleCmd)
}

const (
	ramEndpoint       = "https://ram.aliyuncs.com"
	roleType          = "Service"
	rolePrincipal     = "fc.aliyuncs.com"
	postLogAction     = "PostLogStoreLogs"
	getlogStoreAction = "GetLogStore"
	getObjectAction   = "GetObject"
	loggingPolicyFmt  = "fc-%s-logging-%v"
	ossPolicyFmt      = "fc-%s-oss-%v"
	roleArnFmt        = "acs:ram::%s:role/%s"
)

var roleCmd = &cobra.Command{
	Use:     "role",
	Aliases: []string{"r"},
	Short:   "role related operation",
	Long: `role related operation
	
EXAMPLE:
  fcli role config ......
	`,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() {
	roleCmd.Flags().BoolP("help", "h", true, "Print Usage")
}

// Some util functions for role
func getRAMClient(accessKeyID, accessKeySecret string) (*ram.Client, error) {
	client, err := ram.NewClient(ramEndpoint, accessKeyID, accessKeySecret)
	if err != nil {
		return nil, fmt.Errorf("get ram client err: %s", err)
	}
	return client, nil
}

func createOSSPolicy(client *ram.Client, uid, accessKeyID, accessKeySecret, policyName, bucket string) error {
	// acs:oss:*:uid:bucketName/*
	temp := `{"Version": "1", "Statement": [{"Effect": "Allow", "Action": ["oss:%s"], "Resource": ["acs:oss:*:%s:%s/*"]}]}`
	doc := fmt.Sprintf(temp, getObjectAction, uid, bucket)
	_, err := client.CreatePolicy(policyName, doc, "fc cli OSS copy code policy")
	if err != nil {
		return fmt.Errorf("failed to create policy %s due to %s", policyName, err)
	}
	return nil
}

func createLoggingPolicy(client *ram.Client, accessKeyID, accessKeySecret, policyName, uid, project, logstore string) error {
	// return policyName
	// acs:log:*:uid:project/projectName/logstore/*
	// acs:log:*:uid:project/projectName/logstore/logstoreName
	if logstore == "" {
		logstore = "*"
	}
	temp := `{"Version": "1", "Statement": [{"Effect": "Allow", "Action": ["log:%s"], "Resource": ["acs:log:*:%s:project/%s/logstore/%s"]}]}`
	doc := fmt.Sprintf(temp, postLogAction, uid, project, logstore)

	_, err := client.CreatePolicy(policyName, doc, "fc cli logging policy")
	if err != nil {
		return fmt.Errorf("failed to create policy %s due to :%s", policyName, err)
	}
	return nil
}

func isRoleExist(client *ram.Client, roleName string) (bool, error) {
	_, err := client.GetRole(roleName)
	if err != nil {
		if value, ok := err.(ram.ServiceError); ok {
			if value.ErrorCode == "EntityNotExist.Role" {
				return false, nil
			}
			return false, err
		}
		return false, err
	}
	return true, nil
}

func createRole(client *ram.Client, accessKeyID, accessKeySecret, roleName string) (string, error) {
	// return roleArn
	temp := `{"Statement": [{"Action": "sts:AssumeRole", "Effect": "Allow", "Principal": { "%s": ["%s"]}}], "Version": "1"}`
	doc := fmt.Sprintf(temp, roleType, rolePrincipal)
	resp, err := client.CreateRole(roleName, doc, "fc cli role")
	if err != nil {
		return "", fmt.Errorf("failed to create role %s due to:%s", roleName, err)
	}
	return resp.Role.Arn, nil
}

func attachPolicyToRole(client *ram.Client, policyName string, roleName string) error {
	_, err := client.AttachPolicyToRole("Custom", policyName, roleName)
	if err != nil {
		return fmt.Errorf("failed to attach policy %s to %s, due to %s", policyName, roleName, err)
	}
	return nil
}

func isEmpty(s *string) bool {
	if s == nil {
		return true
	}
	if len(strings.TrimSpace(*s)) == 0 {
		return true
	}
	return false
}

func getUserID(endpoint string) string {
	tmp := strings.SplitN(endpoint, ".", 2)
	uid := tmp[0]
	uid = strings.TrimPrefix(uid, "http://")
	uid = strings.TrimPrefix(uid, "https://")
	return uid
}

func getRoleArn(uid, roleName string) string {
	return fmt.Sprintf(roleArnFmt, uid, roleName)
}
