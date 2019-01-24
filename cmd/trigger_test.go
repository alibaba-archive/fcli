package cmd

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

// TestTrigger ..
func TestTrigger(t *testing.T) {
	suite.Run(t, new(TriggerTestSuite))
}

type TriggerTestSuite struct {
	suite.Suite
}

// TestGetTriggerHelp ..
func (s *TriggerTestSuite) TestGetTriggerHelp() {
	str := getTriggerHelp("","")
	s.NotEmpty(str)
}

// TestPrintTrigger ..
func (s *TriggerTestSuite) TestPrintTrigger() {
	printTrigger()
}