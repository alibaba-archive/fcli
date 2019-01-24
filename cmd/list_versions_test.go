package cmd

import (
	"fmt"
	"github.com/aliyun/fcli/util"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
)

// TestListVersion ..
func TestListVersion(t *testing.T) {
	suite.Run(t, new(ListVersionTestSuite))
}

type ListVersionTestSuite struct {
	suite.Suite
}

// TestListVersionRun ..
func (s *ListVersionTestSuite) TestListVersionRun() {
	// error
	mockObject := &MockedManager{}
	mockObject.
		On("ListServiceVersions", mock.Anything).
		Return("", fmt.Errorf("error"))
	util.Client = mockObject
	err := listVersionsCmd.RunE(&cobra.Command{}, nil)
	s.NotNil(err)

	// success
	mockObject = &MockedManager{}
	mockObject.
		On("ListServiceVersions", mock.Anything).
		Return("{}", nil)
	util.Client = mockObject
	err = listVersionsCmd.RunE(&cobra.Command{}, nil)
	s.Nil(err)

	util.Client = nil
}
