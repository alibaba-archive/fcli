package cmd

import (
	"fmt"
	"github.com/aliyun/fcli/util"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
)

// TestGetFunction ..
func TestGetFunction(t *testing.T) {
	suite.Run(t, new(GetFunctionTestSuite))
}

type GetFunctionTestSuite struct {
	suite.Suite
}

// TestGetFunctionRun ..
func (s *GetFunctionTestSuite) TestGetFunctionRun() {
	// error
	mockObject := &MockedManager{}
	mockObject.
		On("GetFunction", mock.Anything).
		Return("", fmt.Errorf("error"))
	util.Client = mockObject
	err := getFuncCmd.RunE(&cobra.Command{}, nil)
	s.NotNil(err)

	// success
	mockObject = &MockedManager{}
	mockObject.
		On("GetFunction", mock.Anything).
		Return("{}", nil)
	util.Client = mockObject
	err = getFuncCmd.RunE(&cobra.Command{}, nil)
	s.Nil(err)

	util.Client = nil
}
