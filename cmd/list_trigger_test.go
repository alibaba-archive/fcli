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

// TestListTrigger ..
func TestListTrigger(t *testing.T) {
	suite.Run(t, new(ListTriggerTestSuite))
}

type ListTriggerTestSuite struct {
	suite.Suite
}

// TestListTriggerRun ..
func (s *ListTriggerTestSuite) TestListTrigger() {
	// error
	mockObject := &MockedManager{}
	mockObject.
		On("ListTriggers", mock.Anything).
		Return("", fmt.Errorf("error"))
	util.Client = mockObject
	output, err := listTriggerRun(&cobra.Command{})
	s.Nil(output)
	s.NotNil(err)

	// error isShowAll
	isShowAll = true
	output, err = listTriggerRun(&cobra.Command{})
	s.Nil(output)
	s.NotNil(err)
	isShowAll = false

	// success
	mockObject = &MockedManager{}
	mockObject.
		On("ListTriggers", mock.Anything).
		Return("{}", nil)
	util.Client = mockObject
	output, err = listTriggerRun(&cobra.Command{})
	s.NotNil(output)
	s.Nil(err)

	// success onlyNames
	mockObject = &MockedManager{}
	mockObject.
		On("ListTriggers", mock.Anything).
		Return(`{"triggers":[{"triggerName":"trigger1","triggerConfig":{},"triggerType":"`+fc.TRIGGER_TYPE_OSS+`"}]}`, nil)
	util.Client = mockObject
	onlyNames = true
	output, err = listTriggerRun(&cobra.Command{})
	s.NotNil(output)
	s.Nil(err)
	onlyNames = false

	// success isShowAll
	isShowAll = true
	output, err = listTriggerRun(&cobra.Command{})
	s.NotNil(output)
	s.Nil(err)
	isShowAll = false

	// success isShowAll onlyNames
	isShowAll = true
	onlyNames = true
	output, err = listTriggerRun(&cobra.Command{})
	s.NotNil(output)
	s.Nil(err)
	isShowAll = false
	onlyNames = false

	util.Client = nil
}
