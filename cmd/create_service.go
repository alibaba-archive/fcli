package cmd

import (
	"github.com/aliyun/fc-go-sdk"
	"github.com/aliyun/fcli/util"

	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	serviceCmd.AddCommand(createServiceCmd)

	createServiceCmd.Flags().Bool("help", false, "create service")

	createServiceInput.serviceName = createServiceCmd.Flags().StringP(
		"service-name", "s", "", "the service name")
	createServiceInput.description = createServiceCmd.Flags().String(
		"description", "", "service description")
	createServiceInput.internetAccess = createServiceCmd.Flags().Bool(
		"internet-access", true, "service internet access")
	createServiceInput.logProject = createServiceCmd.Flags().StringP(
		"log-project", "p", "", "loghub project for logging")
	createServiceInput.logStore = createServiceCmd.Flags().StringP(
		"log-store", "l", "", "loghub logstore for logging")
	createServiceInput.role = createServiceCmd.Flags().StringP(
		"role-arn", "r", "", "role arn for oss code copy, function execution and logging")
	createServiceInput.vpcID = createServiceCmd.Flags().StringP(
		"vpc-id", "", "", "vpc id is required to enable the vpc access")
	createServiceInput.vSwitchIDs = createServiceCmd.Flags().StringArrayP(
		"v-switch-ids", "", []string{}, "at least one vswitch id is required to enable the vpc access")
	createServiceInput.securityGroupID = createServiceCmd.Flags().StringP(
		"security-group-id", "", "", "security group id is required to enable the vpc access")
	createServiceInput.nasUserID = createServiceCmd.Flags().Int32P("nas-userid", "u", -1, "user id to access NAS volume")
	createServiceInput.nasGroupID = createServiceCmd.Flags().Int32P("nas-groupid", "g", -1, "group id to access NAS volume")
	createServiceInput.nasServer = createServiceCmd.Flags().StringArrayP("nas-server-addr", "", []string{},
		"at least one nas server is required to enable the NAS access")
	createServiceInput.nasMount = createServiceCmd.Flags().StringArrayP("nas-mount-dir", "", []string{},
		"at least one nas dir is required to enable the NAS access")
}

// ServiceInput defines service input
type createServiceInputType struct {
	serviceName     *string
	description     *string
	internetAccess  *bool
	role            *string
	logProject      *string
	logStore        *string
	vpcID           *string
	vSwitchIDs      *[]string
	securityGroupID *string
	nasUserID       *int32
	nasGroupID      *int32
	nasServer       *[]string
	nasMount        *[]string
}

var createServiceInput createServiceInputType

var createServiceCmd = &cobra.Command{
	Use:     "create [option]",
	Aliases: []string{"c"},
	Short:   "create service",
	Long:    ``,
	Run: func(cmd *cobra.Command, args []string) {
		mountPoints := []fc.NASMountConfig{}
		if len(*createServiceInput.nasServer) != len(*createServiceInput.nasMount) {
			fmt.Println("nas server array length must match nas dir array length")
			return
		}
		for i, addr := range *createServiceInput.nasServer {
			mountPoints = append(mountPoints, fc.NASMountConfig{
				ServerAddr: addr,
				MountDir:   (*createServiceInput.nasMount)[i],
			})
		}
		input := fc.NewCreateServiceInput().
			WithServiceName(*createServiceInput.serviceName).
			WithDescription(*createServiceInput.description).
			WithInternetAccess(*createServiceInput.internetAccess).
			WithLogConfig(fc.NewLogConfig().
				WithProject(*createServiceInput.logProject).
				WithLogstore(*createServiceInput.logStore)).
			WithRole(*createServiceInput.role).
			WithVPCConfig(fc.NewVPCConfig().
				WithVPCID(*createServiceInput.vpcID).
				WithVSwitchIDs(*createServiceInput.vSwitchIDs).
				WithSecurityGroupID(*createServiceInput.securityGroupID)).
			WithNASConfig(fc.NewNASConfig().
				WithUserID(*createServiceInput.nasUserID).
				WithGroupID(*createServiceInput.nasGroupID).
				WithMountPoints(mountPoints))

		client, err := util.NewFClient(gConfig)
		if err != nil {
			fmt.Printf("Error: can not create fc client: %s\n", err)
			return
		}
		_, err = client.CreateService(input)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
		}
	},
}
