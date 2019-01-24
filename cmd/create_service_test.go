package cmd

import (
	"fmt"
	"github.com/aliyun/fcli/util"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
)

// TestCreateService ..
func TestCreateService(t *testing.T) {
	suite.Run(t, new(CreateServiceTestSuite))
}

type CreateServiceTestSuite struct {
	suite.Suite
}

// TestCreateServiceRun ..
func (s *CreateServiceTestSuite) TestCreateServiceRun() {
	// error
	mockObject := &MockedManager{}
	mockObject.
		On("CreateService", mock.Anything).
		Return("", fmt.Errorf("error"))
	util.Client = mockObject
	err := createServiceCmd.RunE(&cobra.Command{}, nil)
	s.NotNil(err)

	// success
	mockObject = &MockedManager{}
	mockObject.
		On("CreateService", mock.Anything).
		Return("{}", nil)
	util.Client = mockObject
	err = createServiceCmd.RunE(&cobra.Command{}, nil)
	s.Nil(err)

	// error nas server nas mount config
	createServiceInput.nasServer = &[]string{"1","2"}
	err = createServiceCmd.RunE(&cobra.Command{}, nil)
	s.NotNil(err)

	// correct nas server nas mount config
	createServiceInput.nasServer = &[]string{"1", "2"}
	createServiceInput.nasMount = &[]string{"1", "2"}
	err = createServiceCmd.RunE(&cobra.Command{}, nil)
	s.Nil(err)

	util.Client = nil
}
