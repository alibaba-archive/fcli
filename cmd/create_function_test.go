package cmd

import (
	"fmt"
	"github.com/aliyun/fcli/util"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
)

// TestCreateFunction ..
func TestCreateFunction(t *testing.T) {
	suite.Run(t, new(CreateFunctionTestSuite))
}

type CreateFunctionTestSuite struct {
	suite.Suite
}

// TestCreateFunctionRun ..
func (s *CreateFunctionTestSuite) TestCreateFunctionRun() {
	// error
	mockObject := &MockedManager{}
	mockObject.
		On("CreateFunction", mock.Anything).
		Return("", fmt.Errorf("error"))
	util.Client = mockObject
	err := createFuncCmd.RunE(&cobra.Command{}, nil)
	s.NotNil(err)

	// codeFile error
	createFuncInput.codeFile = "./error_file.err"
	err = createFuncCmd.RunE(&cobra.Command{}, nil)
	s.NotNil(err)
	createFuncInput.codeFile = ""

	// success
	mockObject = &MockedManager{}
	mockObject.
		On("CreateFunction", mock.Anything).
		Return("{}", nil)
	util.Client = mockObject
	err = createFuncCmd.RunE(&cobra.Command{}, nil)
	s.Nil(err)

	util.Client = nil
}