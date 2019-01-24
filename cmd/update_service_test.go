package cmd

import (
	"fmt"
	"github.com/aliyun/fcli/util"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
)

// TestUpdateService ..
func TestUpdateService(t *testing.T) {
	suite.Run(t, new(UpdateServiceTestSuite))
}

type UpdateServiceTestSuite struct {
	suite.Suite
}

// TestUpdateServiceRun ..
func (s *UpdateServiceTestSuite) TestUpdateServiceRun() {
	// error
	mockObject := &MockedManager{}
	mockObject.
		On("UpdateService", mock.Anything).
		Return("", fmt.Errorf("error"))
	util.Client = mockObject
	err := updateServiceCmd.RunE(&cobra.Command{}, nil)
	s.NotNil(err)

	// success
	mockObject = &MockedManager{}
	mockObject.
		On("UpdateService", mock.Anything).
		Return("{}", nil)
	util.Client = mockObject

	inputCmd := &cobra.Command{}
	inputCmd.Flags().String("description", "", "")
	inputCmd.Flags().String("internet-access", "", "")
	inputCmd.Flags().String("role", "", "")
	inputCmd.Flags().String("log-project", "", "")
	inputCmd.Flags().String("log-store", "", "")
	inputCmd.Flags().String("etag", "", "")
	inputCmd.Flags().String("vpc-id", "", "")
	inputCmd.Flags().String("nas-userid", "", "")
	inputCmd.Flags().String("nas-groupid", "", "")
	inputCmd.Flags().String("nas-server-addr", "", "")

	err = inputCmd.Flags().Set("description", "description")
	err = inputCmd.Flags().Set("internet-access", "internet-access")
	err = inputCmd.Flags().Set("role", "role")
	err = inputCmd.Flags().Set("log-project", "log-project")
	err = inputCmd.Flags().Set("log-store", "log-store")
	err = inputCmd.Flags().Set("etag", "etag")
	err = inputCmd.Flags().Set("vpc-id", "vpc-id")
	err = inputCmd.Flags().Set("nas-userid", "nas-userid")
	err = inputCmd.Flags().Set("nas-groupid", "nas-groupid")
	err = inputCmd.Flags().Set("nas-server-addr", "nas-server-addr")
	updateServiceInput.nasServer = &[]string{"1", "2"}
	updateServiceInput.nasMount = &[]string{"1", "2"}
	err = updateServiceCmd.RunE(inputCmd, nil)

	s.Nil(err)

	util.Client = nil
}
