package cmd

import (
	"fmt"
	"github.com/aliyun/fcli/util"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
)

// TestDeleteFunction ..
func TestDeleteFunction(t *testing.T) {
	suite.Run(t, new(DeleteFunctionTestSuite))
}

type DeleteFunctionTestSuite struct {
	suite.Suite
}

// TestDeleteFunctionRun ..
func (s *DeleteFunctionTestSuite) TestDeleteFunctionRun() {
	// error
	mockObject := &MockedManager{}
	mockObject.
		On("DeleteFunction", mock.Anything).
		Return("", fmt.Errorf("error"))
	util.Client = mockObject
	err := deleteFuncCmd.RunE(&cobra.Command{}, nil)
	s.NotNil(err)

	// success
	inputCmd := &cobra.Command{}
	inputCmd.Flags().String("etag", "", "")
	err = inputCmd.Flags().Set("etag", "etag")
	mockObject = &MockedManager{}
	mockObject.
		On("DeleteFunction", mock.Anything).
		Return("{}", nil)
	util.Client = mockObject
	err = deleteFuncCmd.RunE(inputCmd, nil)
	s.Nil(err)

	util.Client = nil
}
