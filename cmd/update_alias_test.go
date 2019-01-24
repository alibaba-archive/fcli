package cmd

import (
	"fmt"
	"github.com/aliyun/fcli/util"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
)

// TestUpdateAlias ..
func TestUpdateAlias(t *testing.T) {
	suite.Run(t, new(UpdateAliasTestSuite))
}

type UpdateAliasTestSuite struct {
	suite.Suite
}

// TestUpdateAliasRun ..
func (s *UpdateAliasTestSuite) TestUpdateAliasRun() {
	// error
	mockObject := &MockedManager{}
	mockObject.
		On("UpdateAlias", mock.Anything).
		Return("", fmt.Errorf("error"))
	util.Client = mockObject
	err := updateAliasCmd.RunE(&cobra.Command{}, nil)
	s.NotNil(err)

	// success
	mockObject = &MockedManager{}
	mockObject.
		On("UpdateAlias", mock.Anything).
		Return("{}", nil)
	util.Client = mockObject
	err = updateAliasCmd.RunE(&cobra.Command{}, nil)
	s.Nil(err)

	util.Client = nil
}

