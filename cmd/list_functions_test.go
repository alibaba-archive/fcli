package cmd

import (
	"fmt"
	"github.com/aliyun/fcli/util"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
)

// TestListFunction ..
func TestListFunction(t *testing.T) {
	suite.Run(t, new(ListFunctionTestSuite))
}

type ListFunctionTestSuite struct {
	suite.Suite
}

// TestListFunctionRun ..
func (s *ListFunctionTestSuite) TestListFunctionRun() {
	// error
	mockObject := &MockedManager{}
	mockObject.
		On("ListFunctions", mock.Anything).
		Return("", fmt.Errorf("error"))
	util.Client = mockObject
	err := listFuncCmd.RunE(&cobra.Command{}, nil)
	s.NotNil(err)

	// success
	mockObject = &MockedManager{}
	mockObject.
		On("ListFunctions", mock.Anything).
		Return(`{"functions":[{"functionName":"fc"}]}`, nil)
	util.Client = mockObject
	err = listFuncCmd.RunE(&cobra.Command{}, nil)
	s.Nil(err)

	*listFuncInput.nameOnly = false
	err = listFuncCmd.RunE(&cobra.Command{}, nil)
	*listFuncInput.nameOnly = true
	s.Nil(err)


	util.Client = nil
}
