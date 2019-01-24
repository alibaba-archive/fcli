package cmd

import (
	"fmt"
	"github.com/aliyun/fcli/util"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
)

// TestCreateAlias ..
func TestCreateAlias(t *testing.T) {
	suite.Run(t, new(CreateAliasTestSuite))
}

type CreateAliasTestSuite struct {
	suite.Suite
}

// TestCreateAliasRun ..
func (s *CreateAliasTestSuite) TestCreateAliasRun() {
	// error
	mockObject := &MockedManager{}
	mockObject.
		On("CreateAlias", mock.Anything).
		Return("", fmt.Errorf("error"))
	util.Client = mockObject
	err := createAliasCmd.RunE(&cobra.Command{}, nil)
	s.NotNil(err)

	// success
	mockObject = &MockedManager{}
	mockObject.
		On("CreateAlias", mock.Anything).
		Return("{}", nil)
	util.Client = mockObject
	err = createAliasCmd.RunE(&cobra.Command{}, nil)
	s.Nil(err)

	util.Client = nil
}