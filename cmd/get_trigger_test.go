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

// TestGetTrigger ..
func TestGetTrigger(t *testing.T) {
	suite.Run(t, new(GetTriggerTestSuite))
}

type GetTriggerTestSuite struct {
	suite.Suite
}

// TestGetTriggerString ..
func (s *GetTriggerTestSuite) TestGetTriggerString() {
	output := &GetTriggerCliOutput{
		triggerCliOutputDecorate: triggerCliOutputDecorate{},
		GetTriggerOutput:      fc.GetTriggerOutput{},
	}
	str := output.String()
	s.NotEmpty(str)
}

// TestGetTriggerMarshalJSON ..
func (s *GetTriggerTestSuite) TestGetTriggerMarshalJSON() {
	output := &GetTriggerCliOutput{
		triggerCliOutputDecorate: triggerCliOutputDecorate{},
		GetTriggerOutput:      fc.GetTriggerOutput{},
	}
	b, err := output.MarshalJSON()
	s.NotEmpty(b)
	s.Nil(err)
}

// TestGetTriggerRun ..
func (s *GetTriggerTestSuite) TestGetTriggerRun() {
	// error
	mockObject := &MockedManager{}
	mockObject.
		On("GetTrigger", mock.Anything).
		Return("", fmt.Errorf("error"))
	util.Client = mockObject
	output, err := getTriggerRun(&cobra.Command{})
	s.Nil(output)
	s.NotNil(err)

	// success
	mockObject = &MockedManager{}
	mockObject.
		On("GetTrigger", mock.Anything).
		Return(`{"triggerType":"`+fc.TRIGGER_TYPE_OSS+`","triggerConfig":{}}`, nil)
	util.Client = mockObject
	output, err = getTriggerRun(&cobra.Command{})
	s.NotNil(output)
	s.Nil(err)

	util.Client = nil
}