package cmd

import (
	"fmt"
	"github.com/aliyun/fc-go-sdk"
	"github.com/stretchr/testify/suite"
	"strings"
	"testing"
)

// TestHelper ..
func TestHelper(t *testing.T) {
	suite.Run(t, new(HelperTestSuite))
}

type HelperTestSuite struct {
	suite.Suite
}

// TestGetClient ..
func (s *HelperTestSuite) TestGetClient() {
	client, err := getClient()
	s.Nil(err)
	s.NotNil(client)
}

// TestPrintStruct ..
func (s *HelperTestSuite) TestPrintStruct() {
	content := struct {}{}
	str, err := printStruct(content)
	s.Nil(err)
	s.NotNil(str)
}

// TestWrapResponseError ..
func (s *HelperTestSuite) TestWrapResponseError() {
	input := fmt.Errorf("error")
	res := wrapResponseError(input)
	s.Equal(input, res)

	data := `{"HTTPStatus":400}`
	input = fmt.Errorf(data)
	res = wrapResponseError(input)
	s.NotNil(input)

	data = `{"HTTPStatus":403}`
	input = fmt.Errorf(data)
	res = wrapResponseError(input)
	s.NotNil(input)

	data = `{"HTTPStatus":499}`
	input = fmt.Errorf(data)
	res = wrapResponseError(input)
	s.NotNil(input)
}

// TestDecorateTriggerOutput ..
func (s *HelperTestSuite) TestDecorateTriggerOutput() {
	triggerType := fc.TRIGGER_TYPE_HTTP
	output := &triggerCliOutputDecorate{}
	decorateTriggerOutput(&triggerType, output)
	expect := strings.Join([]string{
		gConfig.Endpoint,
		gConfig.APIVersion,
		"proxy",
		serviceName,
		functionName,
	}, "/")
	s.Equal(expect, *output.HTTPTriggerURL)
}

// TestPrettyPrint ..
func (s *HelperTestSuite) TestPrettyPrint() {
	prettyPrint(nil, fmt.Errorf("error"))
	prettyPrint(struct {}{}, nil)
	prettyPrint("", nil)
}