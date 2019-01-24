package cmd

import (
	"github.com/aliyun/fc-go-sdk"
	"github.com/aliyun/fcli/util"

	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	serviceCmd.AddCommand(updateServiceCmd)

	updateServiceCmd.Flags().Bool("help", false, "Print Usage")

	updateServiceInput.serviceName = updateServiceCmd.Flags().StringP(
		"service-name", "s", "", "the service name")
	updateServiceInput.description = updateServiceCmd.Flags().String(
		"description", "", "service brief description")
	updateServiceInput.internetAccess = updateServiceCmd.Flags().Bool(
		"internet-access", true, "service internet access")
	updateServiceInput.role = updateServiceCmd.Flags().StringP(
		"role", "r", "", "the arn of the service RAM role for code copy and logging")
	updateServiceInput.logProject = updateServiceCmd.Flags().StringP(
		"log-project", "p", "", "loghub project name for logging")
	updateServiceInput.logStore = updateServiceCmd.Flags().StringP(
		"log-store", "l", "", "loghub logstore name for logging")
	updateServiceInput.etag = updateServiceCmd.Flags().String(
		"etag", "", "provide etag to do the conditional update. "+
			"If the specified etag does not match the service's, the update will fail.")
	updateServiceInput.vpcID = updateServiceCmd.Flags().StringP(
		"vpc-id", "", "", "vpc id is required to enable the vpc access")
	updateServiceInput.vSwitchIDs = updateServiceCmd.Flags().StringArrayP(
		"v-switch-ids", "", []string{}, "at least one vswitch id is required to enable the vpc access")
	updateServiceInput.securityGroupID = updateServiceCmd.Flags().StringP(
		"security-group-id", "", "", "security group id is required to enable the vpc access")
	updateServiceInput.nasUserID = updateServiceCmd.Flags().Int32P("nas-userid", "u", -1, "user id to access NAS volume")
	updateServiceInput.nasGroupID = updateServiceCmd.Flags().Int32P("nas-groupid", "g", -1, "group id to access NAS volume")
	updateServiceInput.nasServer = updateServiceCmd.Flags().StringArrayP("nas-server-addr", "", []string{},
		"at least one nas server is required to enable the NAS access")
	updateServiceInput.nasMount = updateServiceCmd.Flags().StringArrayP("nas-mount-dir", "", []string{},
		"at least one nas dir is required to enable the NAS access")
}

type updateServiceInputType struct {
	serviceName     *string
	description     *string
	internetAccess  *bool
	logProject      *string
	logStore        *string
	role            *string
	etag            *string
	vpcID           *string
	vSwitchIDs      *[]string
	securityGroupID *string
	nasUserID       *int32
	nasGroupID      *int32
	nasServer       *[]string
	nasMount        *[]string
}

var updateServiceInput updateServiceInputType

var updateServiceCmd = &cobra.Command{
	Use:     "update [option]",
	Aliases: []string{"u"},
	Short:   "update service",
	Long:    ``,

	RunE: func(cmd *cobra.Command, args []string) error {
		input := fc.NewUpdateServiceInput(*updateServiceInput.serviceName)
		if cmd.Flags().Changed("description") {
			input.WithDescription(*updateServiceInput.description)
		}
		if cmd.Flags().Changed("internet-access") {
			input.WithInternetAccess(*updateServiceInput.internetAccess)
		}
		if cmd.Flags().Changed("role") {
			input.WithRole(*updateServiceInput.role)
		}
		if cmd.Flags().Changed("log-project") && cmd.Flags().Changed("log-store") {
			input.WithLogConfig(
				fc.NewLogConfig().WithProject(*updateServiceInput.logProject).
					WithLogstore(*updateServiceInput.logStore))
		}
		if cmd.Flags().Changed("etag") {
			input.WithIfMatch(*updateServiceInput.etag)
		}
		if cmd.Flags().Changed("vpc-id") {
			input.WithVPCConfig(fc.NewVPCConfig().
				WithVPCID(*updateServiceInput.vpcID).
				WithVSwitchIDs(*updateServiceInput.vSwitchIDs).
				WithSecurityGroupID(*updateServiceInput.securityGroupID))
		}
		nasConfig := fc.NewNASConfig()
		if cmd.Flags().Changed("nas-userid") {
			nasConfig.WithUserID(*updateServiceInput.nasUserID)
		}
		if cmd.Flags().Changed("nas-groupid") {
			nasConfig.WithGroupID(*updateServiceInput.nasGroupID)
		}
		if cmd.Flags().Changed("nas-server-addr") || cmd.Flags().Changed("nas-mount-dir") {
			if len(*updateServiceInput.nasServer) != len(*updateServiceInput.nasMount) {
				return fmt.Errorf("nas server array length must match nas dir array length")
			}
			mountPoints := []fc.NASMountConfig{}
			for i, addr := range *updateServiceInput.nasServer {
				mountPoints = append(mountPoints, fc.NASMountConfig{
					ServerAddr: addr,
					MountDir:   (*updateServiceInput.nasMount)[i],
				})
			}
			nasConfig.WithMountPoints(mountPoints)
		}
		input.WithNASConfig(nasConfig)
		client, err := util.NewFClient(gConfig)
		if err != nil {
			return fmt.Errorf("can not create fc client: %s\n", err)
		}
		_, err = client.UpdateService(input)
		if err != nil {
			return fmt.Errorf("%s\n", err)
		}
		return nil
	},
}
