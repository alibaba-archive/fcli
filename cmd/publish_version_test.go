package cmd

import (
	"fmt"
	"github.com/aliyun/fcli/util"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
)

// TestPublishVersion ..
func TestPublishVersion(t *testing.T) {
	suite.Run(t, new(PublishVersionTestSuite))
}

type PublishVersionTestSuite struct {
	suite.Suite
}

// TestPublishVersionRun ..
func (s *PublishVersionTestSuite) TestPublishVersionRun() {
	// error
	mockObject := &MockedManager{}
	mockObject.
		On("PublishServiceVersion", mock.Anything).
		Return("", fmt.Errorf("error"))
	util.Client = mockObject
	err := publishVersionCmd.RunE(&cobra.Command{}, nil)
	s.NotNil(err)

	// success
	mockObject = &MockedManager{}
	mockObject.
		On("PublishServiceVersion", mock.Anything).
		Return("{}", nil)
	util.Client = mockObject
	err = publishVersionCmd.RunE(&cobra.Command{}, nil)
	s.Nil(err)

	util.Client = nil
}
