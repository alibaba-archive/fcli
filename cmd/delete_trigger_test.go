package cmd

import (
	"fmt"
	"github.com/aliyun/fcli/util"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
)

// TestDeleteTrigger ..
func TestDeleteTrigger(t *testing.T) {
	suite.Run(t, new(DeleteTriggerTestSuite))
}

type DeleteTriggerTestSuite struct {
	suite.Suite
}

// TestDeleteTriggerRun ..
func (s *DeleteTriggerTestSuite) TestDeleteTriggerRun() {
	// error
	mockObject := &MockedManager{}
	mockObject.
		On("DeleteTrigger", mock.Anything).
		Return("", fmt.Errorf("error"))
	util.Client = mockObject
	output, err := deleteTriggerRun(&cobra.Command{})
	s.NotNil(err)
	s.Nil(output)

	// success
	inputCmd := &cobra.Command{}
	inputCmd.Flags().String("etag", "", "")
	err = inputCmd.Flags().Set("etag", "etag")
	mockObject = &MockedManager{}
	mockObject.
		On("DeleteTrigger", mock.Anything).
		Return("{}", nil)
	util.Client = mockObject
	output, err = deleteTriggerRun(inputCmd)
	s.Nil(err)
	s.NotNil(output)

	util.Client = nil
}

