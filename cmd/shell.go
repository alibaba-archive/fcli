package cmd

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v2"

	"github.com/abiosoft/ishell"
	"github.com/aliyun/aliyun-log-go-sdk"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"strconv"

	"github.com/aliyun/fc-go-sdk"
	"github.com/aliyun/fcli/util"
	"github.com/aliyun/fcli/version"
)

type shellState struct {
	resrcType    string // user, service, function, role
	resrcAbsPath string
	resrcName    string
}

// consts ...
const (
	dockerRuntimeImageTag = "build"
	dockerRunParameter    = "run -a stdout -a stderr -a stdin --rm -i -t -v %s:/code %s /bin/bash"
	//use this to handle request parameters, not for local env path.
	filepathSeparator string = "/"

	HeaderInvocationCodeVersion = "X-Fc-Invocation-Code-Version"
	InvocationCodeVersionLatest = "Latest"
)

func userConfigString() string {
	var dir string
	// only display work dir, if error happens, ignore
	workDir, err := getWorkDir()
	if err == nil {
		dir += fmt.Sprintf("Work Directory:\n  %s\n", workDir)
	}

	return dir + displayConfig()
}

func findFirstArg(args []string) string {
	flag := false
	for _, arg := range args {
		if strings.HasPrefix(arg, "-") {
			flag = true
		} else if flag {
			flag = false
		} else {
			return arg
		}
	}
	return ""
}

// getWorkDir ...
func getWorkDir() (string, error) {
	return filepath.Abs(filepath.Dir(os.Args[0]))
}

func isAbs(path string) bool {
	// All the resource path is followed as unix style.
	return strings.HasPrefix(filepath.ToSlash(path), "/")
}

// parsePath handle the "." and "..".
func parseAbsPath(path string) []string {
	absPath := filepath.ToSlash(path)
	if !isAbs(absPath) {
		panic(fmt.Errorf("The input path is not absolute path: %s", absPath))
	}
	v := strings.Split(absPath, "/")
	var ret []string
	ret = append(ret, "/")
	for _, entry := range v {
		if entry == "" || entry == "." {
			continue
		} else if entry == ".." {
			// Do not go up since we reach the root directory.
			if len(ret) == 1 {
				continue
			} else {
				ret = ret[0 : len(ret)-1]
			}
		} else {
			ret = append(ret, entry)
		}
	}
	return ret
}

func getFCResrc(client *fc.Client, absPath string) (interface{}, error) {
	v := parseAbsPath(absPath)
	if len(v) == 2 {
		return nil, nil
	} else if len(v) == 3 {
		service := v[2]
		input := fc.NewGetServiceInput(service)
		return client.GetService(input)
	} else if len(v) == 4 {
		service := v[2]
		function := v[3]
		input := fc.NewGetFunctionInput(service, function)
		return client.GetFunction(input)
	} else if len(v) == 5 {
		service := v[2]
		function := v[3]
		trigger := v[4]
		input := fc.NewGetTriggerInput(service, function, trigger)
		return client.GetTrigger(input)
	} else {
		return nil, fmt.Errorf("invalid resource: %s", absPath)
	}
}

var shellCmd = &cobra.Command{
	Use:     "shell",
	Aliases: []string{"s"},
	Short:   "interactive shell",
	Long:    `interactive shell, with tab completion and etc`,
	Run: func(cmd *cobra.Command, args []string) {
		const serviceRolePrincipal = "fc.aliyuncs.com"
		const invocationRolePrincipal = "oss.aliyuncs.com"
		rootDir := string(filepathSeparator)
		fcRootDir := path.Join(string(filepathSeparator), "fc")
		ramRootDir := path.Join(string(filepathSeparator), "ram")
		slsRootDir := path.Join(string(filepathSeparator), "sls")

		client, err := util.NewFClient(gConfig)
		if err != nil {
			fmt.Printf("Can not create fc client: %s\n", err)
			return
		}
		ramCli, err := util.NewRAMClient(gConfig.AccessKeyID, gConfig.AccessKeySecret)
		if err != nil {
			fmt.Printf("Can not create ram client: %s\n", err)
		}
		slsCli := util.NewSLSClient(gConfig)
		state := shellState{
			resrcName:    "user",
			resrcAbsPath: fcRootDir,
		}

		upsertService := func(args []string, op string) error {
			currPath := findFirstArg(args)
			if !isAbs(currPath) {
				currPath = path.Join(state.resrcAbsPath, currPath)
			}
			flags := pflag.NewFlagSet("upsert-service", pflag.ContinueOnError)
			desc := flags.StringP("description", "d", "", "service description")
			internetAccess := flags.Bool("internet-access", true, "service internet access")
			role := flags.StringP(
				"role", "r", "",
				"role arn for oss code copy, function execution and logging")
			logProj := flags.StringP("log-project", "p", "", "loghub project for logging")
			logStore := flags.StringP("log-store", "l", "", "loghub logstore for logging")
			etag := flags.String("etag", "", "service etag for update")
			vpcID := flags.StringP("vpc-id", "", "", "vpc id is required to enable the vpc access")
			vSwitchIDs := flags.StringArrayP("v-switch-ids", "", []string{},
				"at least one vswitch id is required to enable the vpc access")
			securityGroupID := flags.StringP("security-group-id", "", "",
				"security group id is required to enable the vpc access")
			nasUserID := flags.Int32P("nas-userid", "u", -1, "user id to access NAS volume")
			nasGroupID := flags.Int32P("nas-groupid", "g", -1, "group id to access NAS volume")
			nasServer := flags.StringArrayP("nas-server-addr", "", []string{},
				"at least one nas server is required to enable the NAS access")
			nasMount := flags.StringArrayP("nas-mount-dir", "", []string{},
				"at least one nas dir is required to enable the NAS access")

			help := flags.Bool("help", false, "")
			err := flags.Parse(args)
			if err != nil {
				return err
			}
			if *help {
				fmt.Println(flags.FlagUsages())
				return nil
			}
			if !strings.HasPrefix(currPath, fcRootDir) {
				return fmt.Errorf(
					"invalid path: %s. The root directory of FC resources must be /fc", currPath)
			}
			resrcList := parseAbsPath(currPath)
			if len(resrcList) < 3 {
				return fmt.Errorf("missing service parameter")
			}
			if len(*nasServer) != len(*nasMount) {
				return fmt.Errorf("nas server array length must match nas dir array length")
			}

			name := resrcList[2]
			if op == "CreateService" {
				input := fc.NewCreateServiceInput().
					WithServiceName(name).
					WithDescription(*desc).
					WithInternetAccess(*internetAccess).
					WithRole(*role).
					WithLogConfig(
						fc.NewLogConfig().
							WithProject(*logProj).
							WithLogstore(*logStore)).
					WithVPCConfig(
						fc.NewVPCConfig().
							WithVPCID(*vpcID).
							WithVSwitchIDs(*vSwitchIDs).
							WithSecurityGroupID(*securityGroupID))
				if len(*nasServer) > 0 {

					mountPoints := []fc.NASMountConfig{}
					for i, addr := range *nasServer {
						mountPoints = append(mountPoints, fc.NASMountConfig{
							ServerAddr: addr,
							MountDir:   (*nasMount)[i],
						})
					}
					input.WithNASConfig(
						fc.NewNASConfig().
							WithUserID(*nasUserID).
							WithGroupID(*nasGroupID).
							WithMountPoints(mountPoints))
				}
				_, err = client.CreateService(input)
			} else {
				input := fc.NewUpdateServiceInput(name)
				if flags.Changed("description") {
					input.WithDescription(*desc)
				}
				if flags.Changed("internet-access") {
					input.WithInternetAccess(*internetAccess)
				}
				if flags.Changed("role") {
					input.WithRole(*role)
				}
				if flags.Changed("etag") {
					input.WithIfMatch(*etag)
				}
				if flags.Changed("log-project") && flags.Changed("log-store") {
					input.WithLogConfig(
						fc.NewLogConfig().WithProject(*logProj).WithLogstore(*logStore))
				} else if !flags.Changed("log-project") && !flags.Changed("log-store") {
					// Do nothing
				} else {
					return fmt.Errorf("both log project and store need be provided")
				}
				if flags.Changed("vpc-id") {
					input.WithVPCConfig(
						fc.NewVPCConfig().
							WithVPCID(*vpcID).
							WithVSwitchIDs(*vSwitchIDs).
							WithSecurityGroupID(*securityGroupID))
				}
				nasConfig := fc.NewNASConfig()
				if flags.Changed("nas-userid") {
					nasConfig.WithUserID(*nasUserID)
				}
				if flags.Changed("nas-groupid") {
					nasConfig.WithUserID(*nasGroupID)
				}
				if flags.Changed("nas-server-addr") || flags.Changed("nas-mount-dir") {
					mountPoints := []fc.NASMountConfig{}
					for i, addr := range *nasServer {
						mountPoints = append(mountPoints, fc.NASMountConfig{
							ServerAddr: addr,
							MountDir:   (*nasMount)[i],
						})
					}
					nasConfig.WithMountPoints(mountPoints)
				}
				input.WithNASConfig(nasConfig)
				_, err = client.UpdateService(input)
			}
			return err
		}

		mks := &ishell.Cmd{
			Name:     "mks",
			Help:     "create the service",
			LongHelp: "",
			Func: func(c *ishell.Context) {
				if len(c.Args) < 1 {
					c.Err(fmt.Errorf("mks service [flags]"))
					return
				}
				err := upsertService(c.Args, "CreateService")
				c.Err(err)
			},
		}

		ups := &ishell.Cmd{
			Name:     "ups",
			Help:     "update the service",
			LongHelp: "",
			Func: func(c *ishell.Context) {
				if len(c.Args) < 1 {
					c.Err(fmt.Errorf("ups service [flags]"))
					return
				}
				err := upsertService(c.Args, "UpdateService")
				c.Err(err)
			},
		}

		upsertFunction := func(args []string, op string) error {
			currPath := findFirstArg(args)
			if !isAbs(currPath) {
				currPath = path.Join(state.resrcAbsPath, currPath)
			}

			fmt.Print(len(parseAbsPath(currPath)))

			// create function
			flags := pflag.NewFlagSet("upsert-funciton", pflag.ContinueOnError)
			help := flags.Bool("help", false, "")
			desc := flags.String("description", "", "brief description")
			memory := flags.Int32P("memory", "m", 128, "memory size in MB")
			timeout := flags.Int32("timeout", 30, "timeout in seconds")
			initializationTimeout := flags.Int32P("initializationTimeout", "e", 30, "initializer timeout in seconds")
			ossBucket := flags.StringP("code-bucket", "b", "", "oss bucket of the code")
			ossObject := flags.StringP("code-object", "o", "", "oss object of the code")
			codeDir := flags.StringP("code-dir", "d", "", "local code directory")
			codeFile := flags.StringP("code-file", "f", "", "zipped code file")
			runtime := flags.StringP("runtime", "t", "", "function runtime")
			handler := flags.StringP("handler", "h", "", "handler is the entrypoint for the function execution")
			initializer := flags.StringP("initializer", "i", "", "initializer is the entrypoint for the initializer execution")
			etag := flags.String("etag", "", "function etag for update")
			err := flags.Parse(args)
			if err != nil {
				return err
			}
			if *help {
				fmt.Println(flags.FlagUsages())
				return nil
			}
			if op == "CreateFunction" && *handler == "" {
				return fmt.Errorf("please specify the handler parameter")
			}
			if op == "CreateFunction" && *runtime == "" {
				return fmt.Errorf("please specify the runtime parameter")
			}
			if !strings.HasPrefix(currPath, fcRootDir) {
				return fmt.Errorf(
					"invalid path: %s. The root directory of FC resources must be /fc", currPath)
			}
			resrcList := parseAbsPath(currPath)
			if len(resrcList) != 4 {
				return fmt.Errorf(
					"invalid arguments: %s. Function must be created under service", currPath)
			}
			serviceName := resrcList[2]
			functionName := resrcList[3]

			if op == "CreateFunction" {
				input := fc.NewCreateFunctionInput(serviceName).
					WithFunctionName(functionName).
					WithDescription(*desc).
					WithMemorySize(*memory).
					WithTimeout(*timeout).
					WithInitializationTimeout(*initializationTimeout).
					WithHandler(*handler).
					WithInitializer(*initializer).
					WithRuntime(*runtime)
				if *codeFile != "" {
					var data []byte
					data, err = ioutil.ReadFile(*codeFile)
					if err != nil {
						return err
					}
					input.WithCode(fc.NewCode().WithZipFile(data))
				} else if *codeDir != "" {
					if err == nil {
						input.WithCode(fc.NewCode().WithDir(*codeDir))

					}
				} else {
					input.WithCode(fc.NewCode().
						WithOSSBucketName(*ossBucket).
						WithOSSObjectName(*ossObject))
				}
				_, err = client.CreateFunction(input)
			} else {
				input := fc.NewUpdateFunctionInput(serviceName, functionName)
				if flags.Changed("description") {
					input.WithDescription(*desc)
				}
				if flags.Changed("etag") {
					input.WithIfMatch(*etag)
				}
				if flags.Changed("memory") {
					input.WithMemorySize(*memory)
				}
				if flags.Changed("timeout") {
					input.WithTimeout(*timeout)
				}
				if flags.Changed("initializationTimeout") {
					input.WithInitializationTimeout(*initializationTimeout)
				}
				if flags.Changed("handler") {
					input.WithHandler(*handler)
				}
				if flags.Changed("initializer") {
					input.WithInitializer(*initializer)
				}
				if flags.Changed("runtime") {
					input.WithRuntime(*runtime)
				}
				if flags.Changed("code-file") {
					data, err := ioutil.ReadFile(*codeFile)
					if err != nil {
						return err
					}
					input.WithCode(fc.NewCode().WithZipFile(data))
				} else if flags.Changed("code-dir") {
					if err == nil {
						input.WithCode(fc.NewCode().WithDir(*codeDir))
					}
				} else if flags.Changed("code-bucket") && flags.Changed("code-object") {
					input.WithCode(fc.NewCode().
						WithOSSBucketName(*ossBucket).
						WithOSSObjectName(*ossObject))
				} else if !flags.Changed("code-bucket") && !flags.Changed("code-object") {
				} else {
					return fmt.Errorf("both code bucket and object should be provided")
				}
				_, err = client.UpdateFunction(input)
			}
			return err
		}

		mkf := &ishell.Cmd{
			Name: "mkf",
			Help: "create the function",
			Func: func(c *ishell.Context) {
				if len(c.Args) < 1 {
					c.Err(fmt.Errorf("invalid arguments: %s", c.Args))
					fmt.Println("Usage: mkf function [flags]")
					return
				}
				err := upsertFunction(c.Args, "CreateFunction")
				c.Err(err)
			},
		}

		upf := &ishell.Cmd{
			Name: "upf",
			Help: "update the function",
			Func: func(c *ishell.Context) {
				if len(c.Args) < 1 {
					c.Err(fmt.Errorf("invalid arguments: %s", c.Args))
					fmt.Println("Usage: upf function [flags]")
					return
				}
				err := upsertFunction(c.Args, "UpdateFunction")
				c.Err(err)
			},
		}

		upsertTrigger := func(args []string, op string) error {
			currPath := findFirstArg(args)
			if !isAbs(currPath) {
				currPath = path.Join(state.resrcAbsPath, currPath)
			}
			// create/update trigger
			flags := pflag.NewFlagSet("upsert-trigger", pflag.ContinueOnError)
			help := flags.Bool("help", false, "")
			etag := flags.String("etag", "", "trigger etag for update")
			triggerType := flags.StringP("type", "t", "oss", "trigger type, support oss, log, timer, http, cdn_events, mns_topic now")
			sourceARN := flags.StringP("source-arn", "s", "", "event source arn, timer type trigger optional")
			invocationRole := flags.StringP("invocation-role", "r", "", "invocation role, timer type trigger optional")
			triggerConfigFile := flags.StringP("trigger-config", "c", "", "trigger config file")
			err := flags.Parse(args)
			if err != nil {
				return err
			}
			if *help {
				fmt.Println(flags.FlagUsages())
				return nil
			}
			if op == "CreateTrigger" {
				if *triggerType == "" {
					return fmt.Errorf("please specify the type parameter")
				}
				if *triggerConfigFile == "" {
					return fmt.Errorf("please specify the trigger-config parameter")
				}
			}
			if !strings.HasPrefix(currPath, fcRootDir) {
				return fmt.Errorf(
					"invalid path: %s. The root directory of FC resources must be /fc", currPath)
			}
			resrcList := parseAbsPath(currPath)
			if len(resrcList) != 5 {
				return fmt.Errorf(
					"invalid arguments: %s. Trigger must be created under function", currPath)
			}
			serviceName := resrcList[2]
			functionName := resrcList[3]
			triggerName := resrcList[4]
			if op == "CreateTrigger" {
				triggerConfig, err := util.GetTriggerConfig(*triggerType, *triggerConfigFile)
				if err != nil {
					return err
				}
				input := fc.NewCreateTriggerInput(serviceName, functionName).
					WithTriggerName(triggerName).
					WithTriggerType(*triggerType).
					WithTriggerConfig(triggerConfig)

				if *sourceARN != "" {
					input.WithSourceARN(*sourceARN)
				}

				if *invocationRole != "" {
					input.WithInvocationRole(*invocationRole)
				}

				_, err = client.CreateTrigger(input)
				return err
			}
			input := fc.NewUpdateTriggerInput(serviceName, functionName, triggerName)
			if flags.Changed("etag") {
				input.WithIfMatch(*etag)
			}
			if flags.Changed("invocation-role") {
				input.WithInvocationRole(*invocationRole)
			}
			if flags.Changed("trigger-config") {
				triggerOutput, err := client.GetTrigger(fc.NewGetTriggerInput(serviceName, functionName, triggerName))
				if err != nil {
					return err
				}
				triggerType := *triggerOutput.TriggerType

				triggerConfig, err := util.GetTriggerConfig(triggerType, *triggerConfigFile)
				if err != nil {
					return err
				}
				input.WithTriggerConfig(triggerConfig)
			}
			_, err = client.UpdateTrigger(input)
			return err
		}

		mkt := &ishell.Cmd{
			Name: "mkt",
			Help: "create the trigger",
			Func: func(c *ishell.Context) {
				if len(c.Args) < 1 {
					c.Err(fmt.Errorf("invalid arguments: %s", c.Args))
					fmt.Println("Usage: mkt trigger [flags]")
					return
				}
				err := upsertTrigger(c.Args, "CreateTrigger")
				c.Err(err)
			},
		}

		upt := &ishell.Cmd{
			Name: "upt",
			Help: "update the trigger",
			Func: func(c *ishell.Context) {
				if len(c.Args) < 1 {
					c.Err(fmt.Errorf("invalid arguments: %s", c.Args))
					fmt.Println("Usage: upt trigger [flags]")
					return
				}
				err := upsertTrigger(c.Args, "UpdateTrigger")
				c.Err(err)
			},
		}

		createRole := func(args []string, roleType string) (string, error) {
			roleName := findFirstArg(args)
			flags := pflag.NewFlagSet("create-role", pflag.ContinueOnError)
			help := flags.Bool("help", false, "")
			err := flags.Parse(args)
			if err != nil {
				return "", err
			}
			if *help {
				fmt.Println(flags.FlagUsages())
				return "", nil
			}
			principal := ""
			if roleType == "ServiceRole" {
				principal = serviceRolePrincipal
			} else {
				principal = invocationRolePrincipal
			}
			return util.CreateRole(ramCli, roleName, principal)
		}

		mksr := &ishell.Cmd{
			Name: "mksr",
			Help: "create the service role",
			Func: func(c *ishell.Context) {
				_, err := createRole(c.Args, "ServiceRole")
				if err != nil {
					c.Err(err)
				}
			},
		}

		mkir := &ishell.Cmd{
			Name: "mkir",
			Help: "create the invocation role",
			Func: func(c *ishell.Context) {
				_, err := createRole(c.Args, "InvocationRole")
				if err != nil {
					c.Err(err)
				}
			},
		}

		mkrp := &ishell.Cmd{
			Name: "mkrp",
			Help: "create the ram policy",
			Func: func(c *ishell.Context) {
				policyName := findFirstArg(c.Args)
				flags := pflag.NewFlagSet("create-role", pflag.ContinueOnError)
				action := flags.StringP("action", "a", "", "action of the policy")
				resource := flags.StringP("resource", "r", "", "resource of the policy")
				help := flags.Bool("help", false, "")
				err := flags.Parse(c.Args)
				if err != nil {
					c.Err(err)
					return
				}
				if *help {
					c.Println("mkrp: create the ram policy")
					c.Println("flags:")
					c.Println(flags.FlagUsages())
					return
				}
				err = util.CreatePolicy(ramCli, policyName, *action, *resource)
				c.Err(err)
			},
		}

		mkl := &ishell.Cmd{
			Name: "mkl",
			Help: "create the log project and store",
			Func: func(c *ishell.Context) {
				flags := pflag.NewFlagSet("mkl", pflag.ContinueOnError)
				projName := flags.StringP("project", "p", "", "the log project")
				storeName := flags.StringP("store", "s", "", "the log store")
				shardCnt := flags.Int("shard", 1, "the shard count of the log store")
				ttl := flags.Int("ttl", 30, "the ttl of the log store, in days")
				help := flags.Bool("help", false, "")
				err := flags.Parse(c.Args)
				if err != nil {
					c.Err(err)
					return
				}
				if *help {
					c.Println("mkl: create the log project/store")
					c.Println("flags:")
					c.Println(flags.FlagUsages())
					return
				}

				if *projName == "" || *storeName == "" {
					c.Err(fmt.Errorf("log project/store name can not be empty"))
					return
				}

				proj, err := slsCli.GetProject(*projName)
				if proj == nil {
					proj, err = slsCli.CreateProject(
						*projName, "project created by function compute cli")
				}
				if proj == nil {
					c.Err(err)
					return
				}

				store, err := proj.GetLogStore(*storeName)
				if store == nil {
					c.Println(
						"Note: you have to pay at least 0.04 RMB/day " +
							"for the log store resource. For the detail billing " +
							"info please refer to:\n" +
							"https://www.aliyun.com/price/product#/sls/detail\n" +
							"Do you want to create the log store? [y/n]\n")
					reader := bufio.NewReader(os.Stdin)
					for {
						s, err := reader.ReadString('\n')
						if err != nil {
							c.Err(err)
							return
						}
						s = strings.TrimSpace(s)
						if s == "y" || s == "yes" || s == "Y" || s == "YES" {
							break
						} else if s == "n" || s == "no" || s == "N" || s == "NO" {
							return
						} else {
							fmt.Printf("Please input y/n:\n")
						}
					}
					proj.CreateLogStore(*storeName, *ttl, *shardCnt)
					store, err = proj.GetLogStore(*storeName)
				}
				if store == nil {
					c.Err(err)
					return
				}

				index := sls.Index{
					Keys: map[string]sls.IndexKey{
						"functionName": {
							Token: []string{
								"\n", "\t", ";", ",", "=", ":"},
							CaseSensitive: false,
							Type:          "text",
						},
					},
				}
				err = store.CreateIndex(index)
				if err != nil {
					err = store.UpdateIndex(index)
				}
				if err != nil {
					c.Err(err)
				}
			},
		}

		rm := &ishell.Cmd{
			Name: "rm",
			Help: "delete the resource",
			Func: func(c *ishell.Context) {
				flags := pflag.NewFlagSet("rm", pflag.ContinueOnError)
				if len(c.Args) == 0 {
					c.Err(fmt.Errorf("rm resource [flags]"))
					return
				}
				currPath := ""
				for _, arg := range c.Args {
					if !strings.HasPrefix(arg, "-") {
						currPath = arg
					}
				}
				help := flags.Bool("help", false, "")
				forced := flags.BoolP("forced", "f", false, "Remove the resource without confirmation.")
				err := flags.Parse(c.Args)
				if err != nil {
					c.Err(err)
					return
				}

				if *help {
					fmt.Println(flags.FlagUsages())
					return
				}

				if !isAbs(currPath) {
					currPath = path.Join(state.resrcAbsPath, currPath)
				}

				if !*forced {
					fmt.Printf("Do you want to remove the resource %s [y/n]:\n", currPath)
					reader := bufio.NewReader(os.Stdin)
					for {
						s, err := reader.ReadString('\n')
						if err != nil {
							c.Err(err)
							return
						}
						s = strings.TrimSpace(s)
						if s == "y" || s == "yes" || s == "Y" || s == "YES" {
							break
						} else if s == "n" || s == "no" || s == "N" || s == "NO" {
							return
						} else {
							fmt.Printf("Please input y/n:\n")
						}
					}
				}

				resrcList := parseAbsPath(currPath)
				if strings.HasPrefix(currPath, fcRootDir) {
					if len(resrcList) == 3 {
						// /fc/{service}
						serviceName := resrcList[2]
						input := fc.NewDeleteServiceInput(serviceName)
						_, err = client.DeleteService(input)
					} else if len(resrcList) == 4 {
						// /fc/{service}/{function}
						serviceName := resrcList[2]
						functionName := resrcList[3]
						input := fc.NewDeleteFunctionInput(serviceName, functionName)
						_, err = client.DeleteFunction(input)
					} else if len(resrcList) == 5 {
						// /fc/{service}/{function}/{trigger}
						serviceName := resrcList[2]
						functionName := resrcList[3]
						triggerName := resrcList[4]
						input := fc.NewDeleteTriggerInput(serviceName, functionName, triggerName)
						_, err = client.DeleteTrigger(input)
					} else {
						err = fmt.Errorf("invalid arguments: %s", currPath)
					}
				} else if strings.HasPrefix(currPath, ramRootDir) {
					if len(resrcList) != 4 {
						c.Err(fmt.Errorf("invalid arguments: %s", currPath))
						return
					}

					if strings.HasPrefix(currPath, path.Join(ramRootDir, "roles")) {
						roleName := resrcList[3]
						_, err = ramCli.DeleteRole(roleName)
					} else if strings.HasPrefix(currPath, path.Join(ramRootDir, "policies")) {
						policyName := resrcList[3]
						_, err = ramCli.DeletePolicy(policyName)
					} else {
						err = fmt.Errorf("invalid arguments: %s", currPath)
					}
				}
				c.Err(err)
			},
		}

		fcls := func(c *ishell.Context) {
			flags := pflag.NewFlagSet("ls", pflag.ContinueOnError)
			limit := flags.Int32P("limit", "l", 100, "limit number")
			nextToken := flags.StringP("next-token", "t", "", "list service with specified next token")
			prefix := flags.StringP("prefix", "p", "", "list service with specified prefix")
			startKey := flags.StringP("start-key", "k", "", "list service with specified start key")
			help := flags.Bool("help", false, "")
			err := flags.Parse(c.Args)
			if err != nil {
				c.Err(err)
				return
			}
			if *help {
				fmt.Println(flags.FlagUsages())
				return
			}
			currPath := findFirstArg(c.Args)
			if !isAbs(currPath) {
				currPath = path.Join(state.resrcAbsPath, currPath)
			}
			resrcList := parseAbsPath(currPath)
			if len(resrcList) == 2 {
				// "/fc"
				input := fc.NewListServicesInput().
					WithLimit(*limit).
					WithNextToken(*nextToken).
					WithPrefix(*prefix).
					WithStartKey(*startKey)
				resp, err := client.ListServices(input)
				if err != nil {
					c.Err(err)
					return
				}
				for _, svr := range resp.Services {
					fmt.Printf("%s\n", *svr.ServiceName)
				}
				if resp.NextToken != nil {
					fmt.Printf("NextToken: %s\n", *resp.NextToken)
				}
			} else if len(resrcList) == 3 {
				// "/fc/{service}"
				serviceName := resrcList[2]
				input := fc.NewListFunctionsInput(serviceName).
					WithLimit(*limit).
					WithNextToken(*nextToken).
					WithPrefix(*prefix).
					WithStartKey(*startKey)
				resp, err := client.ListFunctions(input)
				if err != nil {
					c.Err(err)
					return
				}
				for _, svr := range resp.Functions {
					fmt.Printf("%s\n", *svr.FunctionName)
				}
				if resp.NextToken != nil {
					fmt.Printf("NextToken: %s\n", *resp.NextToken)
				}
			} else if len(resrcList) == 4 {
				// "/fc/{service}/{function}"
				serviceName := resrcList[2]
				functionName := resrcList[3]
				input := fc.NewListTriggersInput(serviceName, functionName).
					WithLimit(*limit).
					WithNextToken(*nextToken).
					WithPrefix(*prefix).
					WithStartKey(*startKey)
				resp, err := client.ListTriggers(input)
				if err != nil {
					c.Err(err)
					return
				}
				for _, tri := range resp.Triggers {
					fmt.Printf("%s\n", *tri.TriggerName)
				}
				if resp.NextToken != nil {
					fmt.Printf("NextToken: %s\n", *resp.NextToken)
				}
			} else {
				c.Err(fmt.Errorf("resource does not exist: %s", path.Join(resrcList...)))
			}
		}

		listRole := func(c *ishell.Context) {
			resp, err := ramCli.ListRoles()
			if err != nil {
				c.Err(err)
				return
			}
			for _, r := range resp.Roles.Role {
				c.Println(r.RoleName)
			}
		}

		listPolicy := func(c *ishell.Context) {
			resp, err := ramCli.ListPolicies()
			if err != nil {
				c.Err(err)
				return
			}
			for _, r := range resp.Policies.Policy {
				c.Println(r.PolicyName)
			}
		}

		ls := &ishell.Cmd{
			Name: "ls",
			Help: "list the child resources of the current resource",
			Func: func(c *ishell.Context) {
				currPath := findFirstArg(c.Args)
				if !isAbs(currPath) {
					currPath = path.Join(state.resrcAbsPath, currPath)
				}
				if currPath == rootDir {
					c.Println(fcRootDir)
					c.Println(ramRootDir)
					c.Println(slsRootDir)
				} else if currPath == ramRootDir ||
					currPath == path.Join(ramRootDir, string(filepathSeparator)) {
					c.Println("roles")
					c.Println("policies")
				} else if strings.HasPrefix(currPath, fcRootDir) {
					fcls(c)
				} else if strings.HasPrefix(currPath, path.Join(ramRootDir, "roles")) {
					listRole(c)
				} else if strings.HasPrefix(currPath, path.Join(ramRootDir, "policies")) {
					listPolicy(c)
				} else {
					c.Err(fmt.Errorf("invalid path: %s", currPath))
				}
			},
		}

		dlc := &ishell.Cmd{
			Name: "dlc",
			Help: "download the function code",
			Func: func(c *ishell.Context) {
				flags := pflag.NewFlagSet("dlc", pflag.ContinueOnError)
				currPath := findFirstArg(c.Args)
				help := flags.Bool("help", false, "")
				outputFilePath := flags.StringP("output", "o", "", "write the response to a file")
				err := flags.Parse(c.Args)
				if err != nil {
					c.Err(err)
					return
				}
				if *help {
					fmt.Println(flags.FlagUsages())
					return
				}
				if !isAbs(currPath) {
					currPath = path.Join(state.resrcAbsPath, currPath)
				}
				resrcList := parseAbsPath(currPath)
				if len(resrcList) != 4 {
					c.Err(fmt.Errorf("invalid function: %s", currPath))
					return
				}
				serviceName := resrcList[2]
				functionName := resrcList[3]

				input := fc.NewGetFunctionCodeInput(serviceName, functionName)
				resp, err := client.GetFunctionCode(input)
				if err != nil {
					c.Err(err)
					return
				}
				// If user do not specify the output file path, then save it in the current directory.
				if *outputFilePath == "" {
					*outputFilePath = fmt.Sprintf("%s.zip", functionName)
				}

				c.Println("Downloading code from ", resp.URL)
				c.Println("Save it to ", *outputFilePath)

				output, err := os.Create(*outputFilePath)
				if err != nil {
					c.Err(err)
					return
				}
				defer output.Close()

				err = util.DownloadFromURL(resp.URL, output)
				if err != nil {
					c.Err(err)
				}
			},
		}

		invk := &ishell.Cmd{
			Name: "invk",
			Help: "invoke the function",
			Func: func(c *ishell.Context) {
				flags := pflag.NewFlagSet("invk", pflag.ContinueOnError)
				currPath := findFirstArg(c.Args)
				help := flags.Bool("help", false, "")
				eventFile := flags.StringP("event-file", "f", "", "read event from file")
				eventStr := flags.StringP("event-str", "s", "", "read event from string")
				outputFileName := flags.StringP("output-file", "o", "", "write the response to a file")
				invkType := flags.StringP(
					"type", "t", "Sync",
					"Sync/Async, invoke function synchronously or asynchronously")
				encode := flags.BoolP("encode", "e", false, "encode the response payload by base64")
				err := flags.Parse(c.Args)
				if err != nil {
					c.Err(err)
					return
				}
				if *help {
					fmt.Println(flags.FlagUsages())
					return
				}

				var event []byte
				if *eventStr != "" {
					event = []byte(*eventStr)
				} else if *eventFile != "" {
					data, err := ioutil.ReadFile(*eventFile)
					if err != nil {
						c.Err(fmt.Errorf(
							"Failed to read event from file: %s. Error: %v",
							*eventFile, err))
					}
					event = data
				}

				if !isAbs(currPath) {
					currPath = path.Join(state.resrcAbsPath, currPath)
				}
				resrcList := parseAbsPath(currPath)
				if len(resrcList) != 4 {
					c.Err(fmt.Errorf("invalid function: %s", currPath))
					return
				}

				serviceName := resrcList[2]
				functionName := resrcList[3]
				input := fc.NewInvokeFunctionInput(serviceName, functionName).
					WithInvocationType(*invkType).WithPayload(event).
					WithHeader(HeaderInvocationCodeVersion, InvocationCodeVersionLatest)
				resp, err := client.InvokeFunction(input)
				if err != nil {
					c.Err(err)
					return
				}
				errType := resp.GetErrorType()
				requestID := resp.GetRequestID()
				if errType != "" {
					c.Err(fmt.Errorf("Request id: %s. Error type: %s", requestID, errType))
				} else if gConfig.Debug {
					fmt.Printf("Request id: %s\n", requestID)
				}
				if *outputFileName != "" {
					err = ioutil.WriteFile(*outputFileName, resp.Payload, 0666)
					c.Err(err)
				} else if *encode {
					fmt.Println(base64.StdEncoding.EncodeToString(resp.Payload))
				} else {
					fmt.Println(string(resp.Payload))
				}
			},
		}

		logs := &ishell.Cmd{
			Name:     "logs",
			Help:     "display the service/function logs",
			LongHelp: "",
			Func: func(c *ishell.Context) {
				flags := pflag.NewFlagSet("logs", pflag.ContinueOnError)
				currPath := findFirstArg(c.Args)
				help := flags.Bool("help", false, "")
				start := flags.StringP("start", "s", "", "the start time of the logs. "+
					"time format is UTC RFC3339, such as 2017-01-01T01:02:03Z")
				end := flags.StringP("end", "e", "", "the end time of the logs. "+
					"time format is UTC RFC3339, such as 2017-01-01T01:02:03Z")
				count := flags.Int64P("count", "c", 1000, "the max count of returned lines")
				tail := flags.BoolP("tail", "t", false,
					"prints the last 'count' lines to standard output")
				duration := flags.Int64P(
					"duration", "d", 24*3600,
					"get the logs in the duration, in seconds")
				err := flags.Parse(c.Args)
				if err != nil {
					c.Err(err)
					return
				}
				if *help {
					fmt.Println(flags.FlagUsages())
					return
				}

				now := time.Now()
				startTimestamp := now.Add(-1 * time.Duration(*duration) * time.Second)
				endTimestamp := now
				if *start != "" {
					tmp, err := time.Parse(util.TimeLayoutInLogs, *start)
					if err != nil {
						c.Err(fmt.Errorf("invalid start time format: %s", *start))
						return
					}
					startTimestamp = tmp
				}
				if *end != "" {
					tmp, err := time.Parse(util.TimeLayoutInLogs, *end)
					if err != nil {
						c.Err(fmt.Errorf("invalid start time format: %s", *end))
						return
					}
					endTimestamp = tmp
				}

				if !isAbs(currPath) {
					currPath = path.Join(state.resrcAbsPath, currPath)
				}
				resrcList := parseAbsPath(currPath)
				var serviceName string
				var functionName string
				if len(resrcList) == 3 {
					// /fc/{service}
					serviceName = resrcList[2]
				} else if len(resrcList) == 4 {
					// /fc/{service}/{function}
					serviceName = resrcList[2]
					functionName = resrcList[3]
				} else {
					c.Err(fmt.Errorf("invalid service or function: %s", currPath))
					return
				}

				smeta, err := client.GetService(fc.NewGetServiceInput(serviceName))
				if err != nil {
					c.Err(err)
					return
				}

				projName := *smeta.LogConfig.Project
				storeName := *smeta.LogConfig.Logstore
				if projName == "" {
					c.Err(fmt.Errorf("the service %s has no log project", serviceName))
					return
				}

				if storeName == "" {
					c.Err(fmt.Errorf("the service %s has no log store", serviceName))
					return
				}

				slsProject, err := slsCli.GetProject(projName)
				if err != nil {
					c.Err(fmt.Errorf("failed to get project %s: %v", projName, err))
					return
				}
				slsLogstore, err := slsProject.GetLogStore(storeName)
				if err != nil {
					c.Err(fmt.Errorf("failed to get store %s: %v", storeName, err))
					return
				}
				err = util.GetLogs(
					slsLogstore, serviceName, functionName,
					startTimestamp.Unix(), endTimestamp.Unix(), *count, *tail)
				if err != nil {
					c.Err(fmt.Errorf(`failed to get logs of store "%s": %v`, storeName, err))
				}
			},
		}

		fcinfo := func(c *ishell.Context, currPath string) {
			if !isAbs(currPath) {
				currPath = path.Join(state.resrcAbsPath, currPath)
			}
			resrcList := parseAbsPath(currPath)
			if len(resrcList) == 2 {
				// info "/fc"
				c.Println(userConfigString())
			} else if len(resrcList) == 3 {
				// info "/fc/{service}"
				input := fc.NewGetServiceInput(resrcList[2])
				resp, err := client.GetService(input)
				if err != nil {
					c.Err(err)
					return
				}
				c.Printf("%s\n", resp.String())
			} else if len(resrcList) == 4 {
				// info "/fc/{service}/{function}"
				input := fc.NewGetFunctionInput(resrcList[2], resrcList[3])
				resp, err := client.GetFunction(input)
				if err != nil {
					c.Err(err)
					return
				}
				c.Printf("%s\n", resp.String())
			} else if len(resrcList) == 5 {
				// info "/fc/{service}/{function}/{trigger}"
				input := fc.NewGetTriggerInput(resrcList[2], resrcList[3], resrcList[4])
				resp, err := client.GetTrigger(input)
				if err != nil {
					c.Err(err)
					return
				}
				c.Printf("%s\n", resp.String())
			} else {
				c.Err(fmt.Errorf("resource does not exist: %s", path.Join(resrcList...)))
			}
		}

		raminfo := func(c *ishell.Context, currPath string) {
			resrcList := parseAbsPath(currPath)
			if strings.HasPrefix(currPath, path.Join(ramRootDir, "roles")) {
				if len(resrcList) == 4 {
					// /ram/roles/{role}
					roleName := resrcList[3]
					roleResp, err := ramCli.GetRole(roleName)
					if err != nil {
						c.Err(err)
						return
					}
					policyResp, err := ramCli.ListPoliciesForRole(roleName)
					if err != nil {
						c.Err(err)
						return
					}
					data, _ := json.MarshalIndent(roleResp.Role, "", "  ")
					c.Println("Role:")
					c.Println(string(data))
					c.Println("Attached policies:")
					c.Printf("%v\n", policyResp.Policies.Policy)
				}
			} else if strings.HasPrefix(currPath, path.Join(ramRootDir, "policies")) {
			} else {
				c.Err(fmt.Errorf("invalid path: %s", currPath))
			}
		}

		info := &ishell.Cmd{
			Name: "info",
			Help: "display the resource detail info",
			Func: func(c *ishell.Context) {
				flags := pflag.NewFlagSet("info", pflag.ContinueOnError)
				currPath := ""
				for _, arg := range c.Args {
					if !strings.HasPrefix(arg, "-") {
						currPath = arg
					}
				}
				help := flags.Bool("help", false, "")
				err := flags.Parse(c.Args)
				if err != nil {
					c.Err(err)
					c.Println(flags.FlagUsages())
					return
				}
				if *help {
					c.Println("display the resource detail info")
					c.Println(flags.FlagUsages())
					return
				}
				if !isAbs(currPath) {
					currPath = path.Join(state.resrcAbsPath, currPath)
				}
				if currPath == rootDir ||
					currPath == path.Join(rootDir, string(filepathSeparator)) {
					c.Println(userConfigString())
				} else if strings.HasPrefix(currPath, fcRootDir) {
					fcinfo(c, currPath)
				} else if strings.HasPrefix(currPath, ramRootDir) {
					raminfo(c, currPath)
				} else if strings.HasPrefix(currPath, slsRootDir) {
				}
			},
		}

		cd := &ishell.Cmd{
			Name:     "cd",
			Help:     "change the current resource",
			LongHelp: "",
			Func: func(c *ishell.Context) {
				sz := len(c.Args)
				var tmp shellState
				if sz == 0 {
					tmp.resrcName = ""
					tmp.resrcAbsPath = fcRootDir
				} else if sz == 1 {
					currPath := c.Args[0]
					if !isAbs(currPath) {
						currPath = path.Join(state.resrcAbsPath, currPath)
					}
					resrcList := parseAbsPath(currPath)
					tmp.resrcName = resrcList[len(resrcList)-1]
					tmp.resrcAbsPath = path.Join(resrcList...)
				} else {
					c.Err(fmt.Errorf("invalid arguments: %s", c.Args))
					return
				}

				if strings.HasPrefix(state.resrcAbsPath, fcRootDir) {
					_, err := getFCResrc(client, tmp.resrcAbsPath)
					if err != nil {
						c.Err(err)
						return
					}
				}
				state = tmp
			},
		}

		pwd := &ishell.Cmd{
			Name:     "pwd",
			Help:     "display the current resource",
			LongHelp: "",
			Func: func(c *ishell.Context) {
				c.Printf("%s\n", state.resrcAbsPath)
			},
		}

		attach := &ishell.Cmd{
			Name:     "attach",
			Help:     "attach the policy to a role",
			LongHelp: "",
			Func: func(c *ishell.Context) {
				flags := pflag.NewFlagSet("attach-policy", pflag.ContinueOnError)
				policy := flags.StringP("policy", "p", "", "the source policy")
				policyType := flags.StringP("policy-type", "", "", "the source policy type. Defaults to Custom")
				role := flags.StringP("role", "r", "", "the target role")
				help := flags.Bool("help", false, "")
				err := flags.Parse(c.Args)
				if err != nil {
					c.Err(err)
					return
				}
				if *help {
					c.Println("attach: attach the policy to a role")
					c.Println("usage: attach -p your_policy -r your_role")
					c.Println("flags:")
					c.Println(flags.FlagUsages())
					return
				}
				if !isAbs(*policy) {
					*policy = path.Join(state.resrcAbsPath, *policy)
				}
				resrcList := parseAbsPath(*policy)
				if len(resrcList) != 4 {
					// The path must be: /ram/policies/{policy}
					c.Err(fmt.Errorf("invalid policy path: %s", *policy))
					return
				}
				*policy = resrcList[3]
				if !isAbs(*role) {
					*role = path.Join(state.resrcAbsPath, *role)
				}
				resrcList = parseAbsPath(*role)
				if len(resrcList) != 4 {
					// The path must be: /ram/roles/{role}
					c.Err(fmt.Errorf("invalid role path: %s", *role))
					return
				}
				*role = resrcList[3]
				if *policyType == "" {
					*policyType = "Custom"
				}
				_, err = ramCli.AttachPolicyToRole(*policyType, *policy, *role)
				if err != nil {
					c.Err(err)
				}
			},
		}

		detach := &ishell.Cmd{
			Name:     "detach",
			Help:     "detach the policy from a role",
			LongHelp: "",
			Func: func(c *ishell.Context) {
				flags := pflag.NewFlagSet("detach-policy", pflag.ContinueOnError)
				policy := flags.StringP("policy", "p", "", "the target policy")
				policyType := flags.StringP("policy-type", "", "", "the source policy type. Defaults to Custom")
				role := flags.StringP("role", "r", "", "the source role")
				help := flags.Bool("help", false, "")
				err := flags.Parse(c.Args)
				if err != nil {
					c.Err(err)
					return
				}
				if *help {
					c.Println("detach: detach the policy from a role")
					c.Println("usage: detach -p your_policy -r your_role")
					c.Println("flags:")
					c.Println(flags.FlagUsages())
					return
				}
				if !isAbs(*policy) {
					*policy = path.Join(state.resrcAbsPath, *policy)
				}
				resrcList := parseAbsPath(*policy)
				if len(resrcList) != 4 {
					// The path must be: /ram/policies/{policy}
					c.Err(fmt.Errorf("invalid policy path: %s", *policy))
					return
				}
				*policy = resrcList[3]
				if !isAbs(*role) {
					*role = path.Join(state.resrcAbsPath, *role)
				}
				resrcList = parseAbsPath(*role)
				if len(resrcList) != 4 {
					// The path must be: /ram/roles/{role}
					c.Err(fmt.Errorf("invalid role path: %s", *role))
					return
				}
				*role = resrcList[3]

				if *policyType == "" {
					*policyType = "Custom"
				}
				_, err = ramCli.DetachPolicyFromRole(*policyType, *policy, *role)
				if err != nil {
					c.Err(err)
				}
			},
		}
		bindRolePolicyIdempotent := func(roleName, policyName, action, resource string) (string, error) {
			var roleARN string
			gresp, err := ramCli.GetRole(roleName)
			if gresp != nil {
				roleARN = gresp.Role.Arn
			} else {
				roleARN, err = util.CreateRole(ramCli, roleName, serviceRolePrincipal)
				if err != nil {
					return roleARN, err
				}
			}

			presp, err := ramCli.GetPolicy(policyName, "Custom")
			if presp == nil {
				err = util.CreatePolicy(ramCli, policyName, action, resource)
				if err != nil {
					return roleARN, err
				}
			}

			// check the policy has the resource permission
			policyVersion, err := util.GetDefaultPolicyVersion(ramCli, policyName, "Custom")
			if err != nil {
				return roleARN, err
			}
			hasResourcePermission, _ := util.CheckPolicyResourcePermission(policyVersion, resource)
			if !hasResourcePermission {
				err = util.CreatePolicyVersion(ramCli, policyName, action, resource)
			}
			if err != nil {
				return roleARN, err
			}
			alreadyAttach := false
			policies, err := ramCli.ListPoliciesForRole(roleName)
			if policies != nil {
				for _, v := range policies.Policies.Policy {
					if v.PolicyName == policyName {
						alreadyAttach = true
						break
					}
				}
			}
			if !alreadyAttach {
				err = util.AttachPolicy(ramCli, policyName, roleName)
				if err != nil {
					return roleARN, err
				}
			}
			// has latency when attach policy to role , Sleep 3 second to avoid the inconsistent for policy and role.
			for i := 0; i < 3; i++ {
				fmt.Print(".")
				time.Sleep(time.Second)
			}
			return roleARN, err
		}
		grantWriteLogstore := func(reader *bufio.Reader, serviceName, roleName, policyName string) error {
			fmt.Printf("Please input the log project: ")
			logProj, err := reader.ReadString('\n')
			if err != nil {
				return err
			}
			logProj = strings.TrimSpace(logProj)

			fmt.Printf("Please input the log store: ")
			logStore, err := reader.ReadString('\n')
			if err != nil {
				return err
			}
			logStore = strings.TrimSpace(logStore)

			uid, err := util.GetUIDFromEndpoint(gConfig.Endpoint)
			if err != nil {
				return err
			}
			// acs:log:*:uid:project/projectName/logstore/*
			// acs:log:*:uid:project/projectName/logstore/logstoreName
			action := `["log:PostLogStoreLogs"]`
			resource := fmt.Sprintf(
				`["acs:log:*:%s:project/%s/logstore/%s"]`, uid, logProj, logStore)
			roleARN, err := bindRolePolicyIdempotent(roleName, policyName, action, resource)
			if err != nil {
				return err
			}
			input := fc.NewUpdateServiceInput(serviceName).
				WithLogConfig(&fc.LogConfig{
					Project:  &logProj,
					Logstore: &logStore,
				}).WithRole(roleARN)
			_, err = client.UpdateService(input)
			if err == nil {
				fmt.Println("grant success")
				return nil
			}
			return wrapResponseError(err)
		}

		grantCopyCodeFromOSS := func(reader *bufio.Reader, serviceName, roleName, policyName string) error {
			fmt.Printf("Please input the OSS path (example: your_bucket/your_directory: ")
			path, err := reader.ReadString('\n')
			if err != nil {
				return err
			}
			path = strings.TrimSpace(path)

			uid, err := util.GetUIDFromEndpoint(gConfig.Endpoint)
			if err != nil {
				return err
			}
			// acs:oss:*:uid:ossPath
			action := `["oss:GetObject"]`
			resource := fmt.Sprintf(`["acs:oss:*:%s:%s/*"]`, uid, path)
			roleARN, err := bindRolePolicyIdempotent(roleName, policyName, action, resource)
			input := fc.NewUpdateServiceInput(serviceName).WithRole(roleARN)
			_, err = client.UpdateService(input)
			if err == nil {
				fmt.Println("grant success")
			}
			return err
		}

		grant := &ishell.Cmd{
			Name:     "grant",
			Help:     "grant the permission",
			LongHelp: "",
			Func: func(c *ishell.Context) {
				flags := pflag.NewFlagSet("grant", pflag.ContinueOnError)
				help := flags.Bool("help", false, "")
				err := flags.Parse(c.Args)
				if err != nil {
					c.Err(err)
					return
				}
				if *help {
					c.Println("grant: grant the permission")
					c.Println("usage: grant service/trigger [flags]")
					c.Println("flags:")
					c.Println(flags.FlagUsages())
					return
				}
				currPath := findFirstArg(c.Args)
				if !isAbs(currPath) {
					currPath = path.Join(state.resrcAbsPath, currPath)
				}
				resrcList := parseAbsPath(currPath)
				resrcName := ""
				if len(resrcList) == 3 {
					// /fc/service
					resrcName = resrcList[2]
				} else if len(resrcList) == 5 {
					// /fc/{service}/{function}{trigger}
					resrcName = resrcList[4]
				} else {
					c.Err(fmt.Errorf(
						"invalid path: %s. Only support "+
							"permission grant for service and trigger.",
						currPath))
					return
				}
				reader := bufio.NewReader(os.Stdin)
				numScenarios := 2

				fmt.Println("Please input the role name: ")
				roleName, err := reader.ReadString('\n')
				if err != nil {
					c.Err(err)
					return
				}
				roleName = strings.TrimSpace(roleName)

				fmt.Println("Please input the policy name: ")
				policyName, err := reader.ReadString('\n')
				if err != nil {
					c.Err(err)
					return
				}
				policyName = strings.TrimSpace(policyName)

				fmt.Println("Permission grant scenarios:")
				fmt.Println("1. Allow FC write function logs to your log store.")
				fmt.Println("2. Allow FC copy code from your OSS location.")
				fmt.Printf("Please input your choice [1-%d]:\n", numScenarios)
				scenario := 0
				for {
					s, err := reader.ReadString('\n')
					if err != nil {
						c.Err(err)
						return
					}
					s = strings.TrimSpace(s)
					scenario, err = strconv.Atoi(s)
					if err == nil && scenario <= numScenarios {
						break
					} else {
						c.Printf("Invalid input: %s\n", s)
						c.Printf("Please input your choice [1-%d]:\n", numScenarios)
					}
				}
				switch scenario {
				case 1:
					err = grantWriteLogstore(reader, resrcName, roleName, policyName)
				case 2:
					err = grantCopyCodeFromOSS(reader, resrcName, roleName, policyName)
				}
				c.Err(err)
			},
		}

		supportedRuntimes := map[string]string{
			"python2.7": "aliyunfc/runtime-python2.7",
			"python3":   "aliyunfc/runtime-python3.6",
			"nodejs6":   "aliyunfc/runtime-nodejs6",
			"nodejs8":   "aliyunfc/runtime-nodejs8",
			"java8":     "aliyunfc/runtime-java8",
			"php7.2":    "aliyunfc/runtime-php7.2",
		}

		sbox := &ishell.Cmd{
			Name: "sbox",
			Help: "a sandbox environment for installing " +
				"the 3rd party libararies and trouble shooting",
			LongHelp: "",
			Func: func(c *ishell.Context) {
				flags := pflag.NewFlagSet("config", pflag.ContinueOnError)
				help := flags.Bool("help", false, "")
				codeDir := flags.StringP("code-dir", "d", "", "the code directory")

				supportedRuntimeKeys := make([]string, 0, len(supportedRuntimes))
				for k := range supportedRuntimes {
					supportedRuntimeKeys = append(supportedRuntimeKeys, k)
				}

				runtime := flags.StringP("runtime", "t", "", "supported runtimes :  "+strings.Join(supportedRuntimeKeys, ", "))
				err := flags.Parse(c.Args)
				if err != nil {
					c.Err(err)
					return
				}
				if *help {
					c.Println("sbox: A sandbox environment for installing " +
						"the 3rd party libararies and trouble shooting. " +
						"It's consistent with the function execution environment " +
						"in FunctionCompute service. You should install dependent " +
						"libraries or test your function in the sandbox " +
						"to prenvent from any environment issues.")
					c.Println("note: the sbox feature requires docker is availabe on your machine.")
					c.Println("usage: sbox [flags]")
					c.Println("       type \"exit\" to exit the sandbox.")
					c.Println("flags:")
					c.Println(flags.FlagUsages())
					return
				}

				if *codeDir == "" {
					c.Err(fmt.Errorf("you have to provide the code direcotry and runtime info"))
					return
				}

				*codeDir, err = filepath.Abs(*codeDir)
				if err != nil {
					c.Err(err)
					return
				}
				// only verify the dir exist
				if _, err := os.Stat(*codeDir); err != nil {
					c.Err(fmt.Errorf("can't found path:%s", *codeDir))
					return
				}

				runtimeName := supportedRuntimes[*runtime]
				runtimeQualifier := runtimeName + ":" + dockerRuntimeImageTag

				// check local image
				localImageExisted := util.CheckLocalImage(runtimeName, dockerRuntimeImageTag)
				if localImageExisted {
					// check runtime
					lastDigest, err := util.GetPublicImageDigest(runtimeName, dockerRuntimeImageTag)
					if err != nil {
						c.Err(err)
						return
					}
					current, _ := exec.Command("docker", "inspect", "--format='{{index .RepoDigests 0}}'", runtimeQualifier).Output()
					currentDisgest := strings.Replace(strings.TrimRight(string(current), "\n"), "'", "", -1)
					currentDisgest = currentDisgest[strings.Index(currentDisgest, "@")+1:]
					if lastDigest != currentDisgest {
						c.Println("Warning: Your " + runtimeQualifier + " image is not the latest version")
						c.Println("Warning: You can use 'docker pull " + runtimeQualifier + "' to update image")
					}
				}

				dockerRunArgs := strings.Split(fmt.Sprintf(dockerRunParameter, *codeDir, runtimeQualifier), " ")
				subCmd := util.NewExecutableDockerCmd(dockerRunArgs...)
				c.Println("Entering the container. Your code is in the /code direcotry.")
				err = subCmd.Run()
				c.Err(err)
			},
		}

		config := &ishell.Cmd{
			Name:     "config",
			Help:     "config the fcli",
			LongHelp: "",
			Func: func(c *ishell.Context) {
				flags := pflag.NewFlagSet("config", pflag.ContinueOnError)
				help := flags.Bool("help", false, "")
				debug := flags.Bool("debug", false, "enable/disable debug mode")
				timeout := flags.Uint("timeout", 60, "timeout of the operation")
				endpoint := flags.String("endpoint", "", "endpoint of the function compute service")
				akid := flags.String("access-key-id", "", "access key id")
				aksecret := flags.String("access-key-secret", "", "access key secret")
				securityToken := flags.String("security-token", "", "ram security token")
				err := flags.Parse(c.Args)
				if err != nil {
					c.Err(err)
					return
				}
				if *help {
					c.Println("config: config the fcli")
					c.Println("usage: config [flags]")
					c.Println("flags:")
					c.Println(flags.FlagUsages())
					return
				}

				config, err := getConfigAways()
				if err != nil {
					return
				}
				if flags.Changed("debug") {
					config.Debug = *debug
				}
				if flags.Changed("timeout") {
					config.Timeout = *timeout
				}
				if flags.Changed("endpoint") {
					config.Endpoint = *endpoint
					config.Endpoint = strings.TrimSpace(config.Endpoint)
					config.SLSEndpoint = fmt.Sprintf(
						util.LogEndpointFmt, util.GetRegionNoForSLSEndpoint(config.Endpoint))
					ramCli, err = util.NewRAMClient(config.AccessKeyID, config.AccessKeySecret)
					if err != nil {
						fmt.Printf("Can not create ram client: %s\n", err)
					}
					slsCli = util.NewSLSClient(config)
				}

				if flags.Changed("access-key-id") {
					config.AccessKeyID = *akid
				}
				if flags.Changed("access-key-secret") {
					config.AccessKeySecret = *aksecret
				}
				if flags.Changed("security-token") {
					config.SecurityToken = *securityToken
				}
				client, err = util.NewFClient(config)
				if err != nil {
					fmt.Printf("Can not create fc client: %s\n", err)
					return
				}
				data, err := yaml.Marshal(config)
				if err != nil {
					fmt.Printf("Failed to marshal config: %v. Error: %v\n", data, err)
					return
				}
				err = ioutil.WriteFile(gConfigPath, data, 0600)
				if err != nil {
					fmt.Printf("Failed to write file: %s. Error: %v\n", gConfigPath, err)
					return
				}
			},
		}

		ver := &ishell.Cmd{
			Name:     "version",
			Help:     "display the fcli version",
			LongHelp: "",
			Func: func(c *ishell.Context) {
				c.Printf("fcli version: %s\n", version.Version)
			},
		}

		shell := ishell.New()
		shell.Println("Welcome to the function compute world. Have fun!")
		shell.SetHistoryPath(filepath.Join(gConfigDir, "history"))
		shell.AddCmd(mks)
		shell.AddCmd(ups)
		shell.AddCmd(mkf)
		shell.AddCmd(upf)
		shell.AddCmd(mkt)
		shell.AddCmd(upt)
		shell.AddCmd(mksr)
		shell.AddCmd(mkir)
		shell.AddCmd(mkrp)
		shell.AddCmd(mkl)
		shell.AddCmd(rm)
		shell.AddCmd(ls)
		shell.AddCmd(info)
		shell.AddCmd(dlc)
		shell.AddCmd(invk)
		shell.AddCmd(logs)
		shell.AddCmd(sbox)
		shell.AddCmd(cd)
		shell.AddCmd(pwd)
		shell.AddCmd(attach)
		shell.AddCmd(detach)
		shell.AddCmd(grant)
		shell.AddCmd(config)
		shell.AddCmd(ver)

		shell.Start()
		shell.Wait()
	},
}

func init() {
	RootCmd.AddCommand(shellCmd)
}
