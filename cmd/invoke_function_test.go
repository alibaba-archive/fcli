package cmd

import (
	"fmt"
	"github.com/aliyun/fcli/util"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/mock"
	"testing"

	"github.com/stretchr/testify/suite"
)

type FunctionStructsTestSuite struct {
	suite.Suite
}

func (s *FunctionStructsTestSuite) TestInvokeFuncRun() {
	assert := s.Require()
	defer func() {
		assert.Nil(recover())
	}()
	invokeFuncRun()
}

func TestFunctionStructs(t *testing.T) {
	suite.Run(t, new(FunctionStructsTestSuite))
}

// TestInvokeFunctionRun ..
func (s *FunctionStructsTestSuite) TestInvokeFunctionRun() {
	// error
	mockObject := &MockedManager{}
	mockObject.
		On("InvokeFunction", mock.Anything).
		Return("", fmt.Errorf("error"))
	util.Client = mockObject
	err := invokeFuncCmd.RunE(&cobra.Command{}, nil)
	s.NotNil(err)

	// success
	mockObject = &MockedManager{}
	mockObject.
		On("InvokeFunction", mock.Anything).
		Return("{}", nil)
	util.Client = mockObject
	err = invokeFuncCmd.RunE(&cobra.Command{}, nil)
	s.Nil(err)

	// success set eventStr invkDebugEnabled
	mockObject = &MockedManager{}
	mockObject.
		On("InvokeFunction", mock.Anything).
		Return(`{"Payload":[123,125]}`, nil)
	util.Client = mockObject
	eventStr = "eventStr"
	invkDebugEnabled = true
	err = invokeFuncCmd.RunE(&cobra.Command{}, nil)
	s.Nil(err)
	eventStr = ""
	invkDebugEnabled = false

	util.Client = nil
}