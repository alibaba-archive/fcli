package cmd

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

// TestFunction ..
func TestFunction(t *testing.T) {
	suite.Run(t, new(FunctionTestSuite))
}

type FunctionTestSuite struct {
	suite.Suite
}

