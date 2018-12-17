package cmd

import (
	"fmt"
	"time"

	"github.com/aliyun/fcli/ram"

	"github.com/spf13/cobra"
)

func init() {
	roleCmd.AddCommand(roleConfigCmd)
}

// RoleConfig defines the role config cmd input parameters
type RoleConfig struct {
	project  *string
	logstore *string
	bucket   *string
	roleName *string
}

var roleConfig RoleConfig

var roleConfigCmd = &cobra.Command{
	Use:     "config",
	Aliases: []string{"c"},
	Short:   "role config for logging and oss code copy",
	Long: `
role config
首先确保当前的AK有创建role和policy的权限，如果没有请到RAM控制台上进行授权
fcli role config -r(--role-name) "user defined role name, it is the suffix of role arn"
				-p(--project)   "loghub project"
                -l(--logstore)  "loghub logstore"
				-b(--bucket)    "oss bucket"
example:
roleName: demo-role
roleArn:  acs:ram::1234567:role/demo-role 

-p(--project):   logging授权的必选参数
-l(--logstore):  logging授权的可选参数，不设置时授权project的所有logstore
-b(--bucket):    oss授权的必选参数

-r(--role-name): 两种授权场景的必选参数，如果role存在则在role上增加新的权限，如果role不存在则创建role并增加新的权限
每次执行命令时根据场景(logging和oss）生成新的policy

默认生成的policyName的格式为:
fc-{roleName}-logging-{timestamp}
fc-{roleName}-oss-{timestamp}
		`,
	Run: func(cmd *cobra.Command, args []string) {
		roleArn, policyNameList, err := roleConfigRun()
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("roleArn:", roleArn)
			fmt.Println("policyNameList:")
			for _, v := range policyNameList {
				fmt.Println(v)
			}
		}
	},
}

func roleConfigRun() (string, []string, error) {
	roleArn := ""
	policyNameList := []string{}

	err := prepareCommon()
	if err != nil {
		return "", nil, err
	}
	if isEmpty(roleConfig.roleName) {
		return "", nil, fmt.Errorf("-r(--role-name) is required")
	}

	uid := getUserID(gConfig.Endpoint)
	accessKeyID := gConfig.AccessKeyID
	accessKeySecret := gConfig.AccessKeySecret
	roleName := *roleConfig.roleName
	client, err := getRAMClient(accessKeyID, accessKeySecret)
	if err != nil {
		return "", nil, fmt.Errorf("get ram client failed due to %s", err)
	}

	roleExist, err := isRoleExist(client, roleName)
	if err != nil {
		return "", nil, fmt.Errorf("check role %s exist fail due to %s", roleName, err)
	}

	if !roleExist {
		roleArn, err = createRole(client, accessKeyID, accessKeySecret, roleName)
		if err != nil {
			return "", nil, fmt.Errorf("create role %s fail due to %s", roleName, err)
		}
	} else {
		// role已经存在
		roleArn = getRoleArn(uid, roleName)
	}

	// logging授权
	if !isEmpty(roleConfig.project) {
		loggingPolicy, err := loggingRoleConfig(client, uid, accessKeyID, accessKeySecret, *roleConfig.project, *roleConfig.logstore, roleName)
		if err != nil {
			return "", nil, err
		}
		policyNameList = append(policyNameList, loggingPolicy)
	}

	// oss授权
	if !isEmpty(roleConfig.bucket) {
		ossCopyCodePolicy, err := ossRoleConfig(client, uid, accessKeyID, accessKeySecret, *roleConfig.bucket, roleName)
		if err != nil {
			return "", nil, err
		}
		policyNameList = append(policyNameList, ossCopyCodePolicy)
	}

	return roleArn, policyNameList, nil
}

func ossRoleConfig(client *ram.Client, uid, accessKeyID, accessKeySecret, bucket, roleName string) (string, error) {
	ossPolicyName := fmt.Sprintf(ossPolicyFmt, roleName, time.Now().Unix())
	err := createOSSPolicy(client, uid, accessKeyID, accessKeySecret, ossPolicyName, bucket)
	if err != nil {
		return "", err
	}

	err = attachPolicyToRole(client, ossPolicyName, roleName)
	if err != nil {
		return "", err
	}
	return ossPolicyName, nil
}

func loggingRoleConfig(client *ram.Client, uid, accessKeyID, accessKeySecret, project, logstore, roleName string) (string, error) {
	loggingPolicyName := fmt.Sprintf(loggingPolicyFmt, roleName, time.Now().Unix())
	err := createLoggingPolicy(client, accessKeyID, accessKeySecret, loggingPolicyName, uid, project, logstore)
	if err != nil {
		return "", err
	}

	err = attachPolicyToRole(client, loggingPolicyName, roleName)
	if err != nil {
		return "", err
	}

	return loggingPolicyName, nil
}

func init() {
	roleConfigCmd.Flags().Bool("help", false, "role config for fc")
	roleConfig.project = roleConfigCmd.Flags().StringP("project", "p", "", "loghub project name")
	roleConfig.logstore = roleConfigCmd.Flags().StringP("logstore", "l", "", "loghub logstore name")
	roleConfig.bucket = roleConfigCmd.Flags().StringP("bucket", "b", "", "oss bucket name")
	roleConfig.roleName = roleConfigCmd.Flags().StringP("role-name", "r", "", "user defined role name")
}
