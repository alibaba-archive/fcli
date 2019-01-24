package cmd

import (
	"fmt"
	"github.com/aliyun/fcli/util"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
)

// TestDeleteService ..
func TestDeleteService(t *testing.T) {
	suite.Run(t, new(DeleteServiceTestSuite))
}

type DeleteServiceTestSuite struct {
	suite.Suite
}

// TestDeleteServiceRun ..
func (s *DeleteServiceTestSuite) TestDeleteServiceRun() {
	// error
	mockObject := &MockedManager{}
	mockObject.
		On("DeleteService", mock.Anything).
		Return("", fmt.Errorf("error"))
	util.Client = mockObject
	err := deleteServiceCmd.RunE(&cobra.Command{}, nil)
	s.NotNil(err)

	// success
	inputCmd := &cobra.Command{}
	inputCmd.Flags().String("etag", "", "")
	err = inputCmd.Flags().Set("etag", "etag")
	mockObject = &MockedManager{}
	mockObject.
		On("DeleteService", mock.Anything).
		Return("{}", nil)
	util.Client = mockObject
	err = deleteServiceCmd.RunE(inputCmd, nil)
	s.Nil(err)

	util.Client = nil
}
