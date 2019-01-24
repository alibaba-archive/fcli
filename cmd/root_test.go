package cmd

import (
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
)

func (s *FunctionStructsTestSuite) TestEnv() {
	os.Setenv("ALIYUN_ACCESS_KEY_ID", "TEST_ID")
	os.Setenv("ALIYUN_ACCESS_KEY_SECRET", "TEST_SECRET")
	assert := s.Require()
	pickupConfigFromOldEnv()
	assert.Equal("TEST_ID", gConfig.AccessKeyID)
	assert.Equal("TEST_SECRET", gConfig.AccessKeySecret)

}

func TestRoot(t *testing.T) {
	suite.Run(t, new(FunctionStructsTestSuite))
}
