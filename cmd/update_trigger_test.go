package cmd

import (
	"fmt"
	"github.com/aliyun/fc-go-sdk"
	"github.com/aliyun/fcli/util"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
)

// TestUpdateTrigger ..
func TestUpdateTrigger(t *testing.T) {
	suite.Run(t, new(UpdateTriggerTestSuite))
}

type UpdateTriggerTestSuite struct {
	suite.Suite
}

// TestUpdateTriggerString ..
func (s *UpdateTriggerTestSuite) TestUpdateTriggerString() {
	output := &UpdateTriggerCliOutput{
		triggerCliOutputDecorate: triggerCliOutputDecorate{},
		UpdateTriggerOutput:      fc.UpdateTriggerOutput{},
	}
	str := output.String()
	s.NotEmpty(str)
}

// TestUpdateTriggerMarshalJSON ..
func (s *UpdateTriggerTestSuite) TestUpdateTriggerMarshalJSON() {
	output := &UpdateTriggerCliOutput{
		triggerCliOutputDecorate: triggerCliOutputDecorate{},
		UpdateTriggerOutput:      fc.UpdateTriggerOutput{},
	}
	b, err := output.MarshalJSON()
	s.NotEmpty(b)
	s.Nil(err)
}

// TestUpdateTriggerRun ..
func (s *GetTriggerTestSuite) TestUpdateTriggerRun() {
	// error
	mockObject := &MockedManager{}
	mockObject.
		On("UpdateTrigger", mock.Anything).
		Return("", fmt.Errorf("error"))
	util.Client = mockObject
	output, err := updateTriggerRun(&cobra.Command{})
	s.Nil(output)
	s.NotNil(err)

	// success
	mockObject = &MockedManager{}
	mockObject.
		On("UpdateTrigger", mock.Anything).
		Return(`{"triggerType":"`+fc.TRIGGER_TYPE_OSS+`","triggerConfig":{}}`, nil)
	mockObject.
		On("GetTrigger", mock.Anything).
		Return(`{"triggerType":"`+fc.TRIGGER_TYPE_OSS+`","triggerConfig":{}}`, nil)

	util.Client = mockObject

	inputCmd := &cobra.Command{}
	inputCmd.Flags().String("etag", "", "")
	inputCmd.Flags().String("invocation-role", "", "")
	inputCmd.Flags().String("trigger-config", "", "")
	inputCmd.Flags().String("qualifier", "", "")

	err = inputCmd.Flags().Set("etag", "etag")
	err = inputCmd.Flags().Set("invocation-role", "invocation-role")
	err = inputCmd.Flags().Set("trigger-config", "trigger-config")
	err = inputCmd.Flags().Set("qualifier", "qualifier")
	*updateParam.TriggerConfigFile = "../example/ossTriggerConfig.yaml"
	*updateParam.Qualifier = "qualifier"

	output, err = updateTriggerRun(inputCmd)
	s.NotNil(output)
	s.Nil(err)

	err = inputCmd.Flags().Set("etag", "")
	err = inputCmd.Flags().Set("invocation-role", "")
	err = inputCmd.Flags().Set("trigger-config", "")
	err = inputCmd.Flags().Set("qualifier", "")

	util.Client = nil
}