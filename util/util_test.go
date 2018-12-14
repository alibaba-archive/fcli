package util

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestUtil(t *testing.T) {
	suite.Run(t, new(UtilTestSuite))
}

type UtilTestSuite struct {
	suite.Suite
}

func (s *UtilTestSuite) TestGetRegionNo() {
	testEndpoint := "1234567.fc.cn-hangzhou.aliyuncs.com:8080"
	region := GetRegionNoForSLSEndpoint(testEndpoint)
	s.Equal("cn-hangzhou", region)
}
