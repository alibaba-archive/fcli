package cmd

import (
	"fmt"
	"github.com/aliyun/fcli/util"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
)

// TestDeleteVersion ..
func TestDeleteVersion(t *testing.T) {
	suite.Run(t, new(DeleteVersionTestSuite))
}

type DeleteVersionTestSuite struct {
	suite.Suite
}

// TestDeleteVersionRun ..
func (s *DeleteServiceTestSuite) TestDeleteVersionRun() {
	// error
	mockObject := &MockedManager{}
	mockObject.
		On("DeleteServiceVersion", mock.Anything).
		Return("", fmt.Errorf("error"))
	util.Client = mockObject
	err := deleteVersionCmd.RunE(&cobra.Command{}, nil)
	s.NotNil(err)

	// success
	mockObject = &MockedManager{}
	mockObject.
		On("DeleteServiceVersion", mock.Anything).
		Return("{}", nil)
	util.Client = mockObject
	err = deleteVersionCmd.RunE(&cobra.Command{}, nil)
	s.Nil(err)

	util.Client = nil
}
