package cmd

import (
	"fmt"
	"github.com/aliyun/fcli/util"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
)

// TestGetService ..
func TestGetService(t *testing.T) {
	suite.Run(t, new(GetServiceTestSuite))
}

type GetServiceTestSuite struct {
	suite.Suite
}

// TestGetServicenRun ..
func (s *GetServiceTestSuite) TestGetServicenRun() {
	// error
	mockObject := &MockedManager{}
	mockObject.
		On("GetService", mock.Anything).
		Return("", fmt.Errorf("error"))
	util.Client = mockObject
	err := getServiceCmd.RunE(&cobra.Command{}, nil)
	s.NotNil(err)

	// success
	mockObject = &MockedManager{}
	mockObject.
		On("GetService", mock.Anything).
		Return("{}", nil)
	util.Client = mockObject
	err = getServiceCmd.RunE(&cobra.Command{}, nil)
	s.Nil(err)

	util.Client = nil
}
