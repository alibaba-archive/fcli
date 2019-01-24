package cmd

import (
	"github.com/stretchr/testify/suite"
	"os"
	"strings"
	"testing"
)

// TestConfig ..
func TestConfig(t *testing.T) {
	suite.Run(t, new(ConfigTestSuite))
}

type ConfigTestSuite struct {
	suite.Suite
}

// TestDisplayConfig ..
func (s *ConfigTestSuite) TestDisplayConfig() {
	config := displayConfig()
	s.NotNil(config)
}

// TestDisplayEnv ..
func (s *ConfigTestSuite) TestDisplayEnv() {
	err := os.Setenv("testKey", "testValue")
	if err == nil {
		str := displayEnv("testKey")
		s.NotNil(str)
		s.True(strings.Contains(str, "testValue"))
	}
}