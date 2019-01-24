package cmd

import (
	"fmt"
	"github.com/aliyun/fcli/util"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
)

// TestListService ..
func TestListService(t *testing.T) {
	suite.Run(t, new(ListServiceTestSuite))
}

type ListServiceTestSuite struct {
	suite.Suite
}

// TestListServiceRun ..
func (s *ListServiceTestSuite) TestListServiceRun() {
	// error
	mockObject := &MockedManager{}
	mockObject.
		On("ListServices", mock.Anything).
		Return("", fmt.Errorf("error"))
	util.Client = mockObject
	err := listServiceCmd.RunE(&cobra.Command{}, nil)
	s.NotNil(err)

	// success
	mockObject = &MockedManager{}
	mockObject.
		On("ListServices", mock.Anything).
		Return(`{"services":[{"serviceName":"fc"}]}`, nil)
	util.Client = mockObject
	err = listServiceCmd.RunE(&cobra.Command{}, nil)
	s.Nil(err)

	*listServiceInput.nameOnly = false
	err = listServiceCmd.RunE(&cobra.Command{}, nil)
	*listServiceInput.nameOnly = true
	s.Nil(err)


	util.Client = nil
}
