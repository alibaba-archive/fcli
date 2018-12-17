package cmd

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type FunctionStructsTestSuite struct {
	suite.Suite
}

func (s *FunctionStructsTestSuite) TestInvokeFuncRun() {
	assert := s.Require()
	defer func() {
		assert.Nil(recover())
	}()
	invokeFuncRun()
}

func TestFunctionStructs(t *testing.T) {
	suite.Run(t, new(FunctionStructsTestSuite))
}
