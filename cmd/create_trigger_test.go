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

// TestCreateTrigger ..
func TestCreateTrigger(t *testing.T) {
	suite.Run(t, new(CreateTriggerTestSuite))
}

type CreateTriggerTestSuite struct {
	suite.Suite
}

// TestCreateTriggerString ..
func (s *CreateTriggerTestSuite) TestCreateTriggerString() {
	output := &CreateTriggerCliOutput{
		triggerCliOutputDecorate: triggerCliOutputDecorate{},
		CreateTriggerOutput:      fc.CreateTriggerOutput{},
	}
	str := output.String()
	s.NotEmpty(str)
}

// TestCreateTriggerMarshalJSON ..
func (s *CreateTriggerTestSuite) TestCreateTriggerMarshalJSON() {
	output := &CreateTriggerCliOutput{
		triggerCliOutputDecorate: triggerCliOutputDecorate{},
		CreateTriggerOutput:      fc.CreateTriggerOutput{},
	}
	b, err := output.MarshalJSON()
	s.NotEmpty(b)
	s.Nil(err)
}

// TestPrepareCreateTriggerInput ..
func (s *CreateTriggerTestSuite) TestPrepareCreateTriggerInput() {
	*createTriggerInput.triggerConfigFile = `../example/notFound.yaml`
	input, err := prepareCreateTriggerInput()
	s.Nil(input)
	s.NotNil(err)

	// set triggerType and triggerConfigFile
	*createTriggerInput.triggerType = fc.TRIGGER_TYPE_OSS
	*createTriggerInput.triggerConfigFile = `../example/ossTriggerConfig.yaml`
	*createTriggerInput.sourceARN = "sourceARN"
	*createTriggerInput.invocationRole = "invocationRole"
	*createTriggerInput.qualifier = "qualifier"
	input, err = prepareCreateTriggerInput()
	*createTriggerInput.triggerType = ""
	*createTriggerInput.triggerConfigFile = ""
	*createTriggerInput.sourceARN = ""
	*createTriggerInput.invocationRole = ""
	*createTriggerInput.qualifier = ""
	s.NotNil(input)
	s.Nil(err)
}

// TestCreateTriggerRun ..
func (s *CreateTriggerTestSuite) TestCreateTriggerRun() {
	// set triggerType and triggerConfigFile
	*createTriggerInput.triggerType = fc.TRIGGER_TYPE_OSS
	*createTriggerInput.triggerConfigFile = `../example/ossTriggerConfig.yaml`
	*createTriggerInput.sourceARN = "sourceARN"
	*createTriggerInput.invocationRole = "invocationRole"
	*createTriggerInput.qualifier = "qualifier"

	// error
	mockObject := &MockedManager{}
	mockObject.
		On("CreateTrigger", mock.Anything).
		Return("", fmt.Errorf("error"))
	util.Client = mockObject
	output, err := createTriggerRun(&cobra.Command{})
	s.Nil(output)
	s.NotNil(err)

	// success
	mockObject = &MockedManager{}
	mockObject.
		On("CreateTrigger", mock.Anything).
		Return(`{"triggerType":"`+fc.TRIGGER_TYPE_OSS+`","triggerConfig":{}}`, nil)
	util.Client = mockObject
	output, err = createTriggerRun(&cobra.Command{})
	s.NotNil(output)
	s.Nil(err)

	*createTriggerInput.triggerType = ""
	*createTriggerInput.triggerConfigFile = ""
	*createTriggerInput.sourceARN = ""
	*createTriggerInput.invocationRole = ""
	*createTriggerInput.qualifier = ""

	util.Client = nil
}
