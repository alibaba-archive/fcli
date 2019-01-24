package cmd

import (
	"fmt"
	"github.com/aliyun/fcli/util"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
)

// TestListAlias ..
func TestListAlias(t *testing.T) {
	suite.Run(t, new(ListAliasTestSuite))
}

type ListAliasTestSuite struct {
	suite.Suite
}

// TestListAliasRun ..
func (s *ListAliasTestSuite) TestListAliasRun() {
	// error
	mockObject := &MockedManager{}
	mockObject.
		On("ListAliases", mock.Anything).
		Return("", fmt.Errorf("error"))
	util.Client = mockObject
	err := listAliasesCmd.RunE(&cobra.Command{}, nil)
	s.NotNil(err)

	// success
	mockObject = &MockedManager{}
	mockObject.
		On("ListAliases", mock.Anything).
		Return("{}", nil)
	util.Client = mockObject
	err = listAliasesCmd.RunE(&cobra.Command{}, nil)
	s.Nil(err)

	util.Client = nil
}
