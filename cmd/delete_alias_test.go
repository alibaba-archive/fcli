package cmd

import (
	"fmt"
	"github.com/aliyun/fcli/util"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
)

// TestDeleteAlias ..
func TestDeleteAlias(t *testing.T) {
	suite.Run(t, new(DeleteAliasTestSuite))
}

type DeleteAliasTestSuite struct {
	suite.Suite
}

// TestDeleteAliasRun ..
func (s *DeleteAliasTestSuite) TestDeleteAliasRun() {
	// error
	mockObject := &MockedManager{}
	mockObject.
		On("DeleteAlias", mock.Anything).
		Return("", fmt.Errorf("error"))
	util.Client = mockObject
	err := deleteAliasCmd.RunE(&cobra.Command{}, nil)
	s.NotNil(err)

	// success
	mockObject = &MockedManager{}
	mockObject.
		On("DeleteAlias", mock.Anything).
		Return("{}", nil)
	util.Client = mockObject
	err = deleteAliasCmd.RunE(&cobra.Command{}, nil)
	s.Nil(err)

	util.Client = nil
}