package cmd

import (
	"fmt"
	"github.com/aliyun/fcli/util"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
)

// TestUpdateFunction ..
func TestUpdateFunction(t *testing.T) {
	suite.Run(t, new(UpdateFunctionTestSuite))
}

type UpdateFunctionTestSuite struct {
	suite.Suite
}

// TestUpdateFunctionRun ..
func (s *UpdateFunctionTestSuite) TestUpdateFunctionRun() {
	// error
	mockObject := &MockedManager{}
	mockObject.
		On("UpdateFunction", mock.Anything).
		Return("", fmt.Errorf("error"))
	util.Client = mockObject
	err := updateFuncCmd.RunE(&cobra.Command{}, nil)
	s.NotNil(err)

	// success
	mockObject = &MockedManager{}
	mockObject.
		On("UpdateFunction", mock.Anything).
		Return("{}", nil)
	util.Client = mockObject

	inputCmd := &cobra.Command{}
	inputCmd.Flags().String("description", "", "")
	inputCmd.Flags().String("etag", "", "")
	inputCmd.Flags().String("memory", "", "")
	inputCmd.Flags().String("timeout", "", "")
	inputCmd.Flags().String("initializationTimeout", "", "")
	inputCmd.Flags().String("handler", "", "")
	inputCmd.Flags().String("initializer", "", "")
	inputCmd.Flags().String("runtime", "", "")
	inputCmd.Flags().String("bucket", "", "")
	inputCmd.Flags().String("object", "", "")


	err = inputCmd.Flags().Set("description", "description")
	err = inputCmd.Flags().Set("etag", "etag")
	err = inputCmd.Flags().Set("memory", "memory")
	err = inputCmd.Flags().Set("timeout", "timeout")
	err = inputCmd.Flags().Set("initializationTimeout", "initializationTimeout")
	err = inputCmd.Flags().Set("handler", "handler")
	err = inputCmd.Flags().Set("initializer", "initializer")
	err = inputCmd.Flags().Set("runtime", "runtime")
	err = inputCmd.Flags().Set("bucket", "bucket")
	err = inputCmd.Flags().Set("object", "object")
	err = updateFuncCmd.RunE(inputCmd, nil)

	s.Nil(err)

	util.Client = nil
}
