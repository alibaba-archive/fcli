package cmd

import (
	"fmt"
	"github.com/aliyun/fcli/util"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
)

// TestGetAlias ..
func TestGetAlias(t *testing.T) {
	suite.Run(t, new(GetAliasTestSuite))
}

type GetAliasTestSuite struct {
	suite.Suite
}

// TestGetAliasRun ..
func (s *GetAliasTestSuite) TestGetAliasRun() {
	// error
	mockObject := &MockedManager{}
	mockObject.
		On("GetAlias", mock.Anything).
		Return("", fmt.Errorf("error"))
	util.Client = mockObject
	err := getAliasCmd.RunE(&cobra.Command{}, nil)
	s.NotNil(err)

	// success
	mockObject = &MockedManager{}
	mockObject.
		On("GetAlias", mock.Anything).
		Return("{}", nil)
	util.Client = mockObject
	err = getAliasCmd.RunE(&cobra.Command{}, nil)
	s.Nil(err)

	util.Client = nil
}

