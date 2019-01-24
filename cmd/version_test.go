package cmd

import (
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/suite"
	"testing"
)

// TestVersion ..
func TestVersion(t *testing.T) {
	suite.Run(t, new(VersionTestSuite))
}

type VersionTestSuite struct {
	suite.Suite
}

// TestVersionRun ..
func (s *VersionTestSuite) TestVersionRun() {
	err := versionCmd.RunE(&cobra.Command{}, nil)
	s.Nil(err)
}