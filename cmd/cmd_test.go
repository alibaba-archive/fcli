package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/aliyun/fc-go-sdk"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestCmd(t *testing.T) {
	suite.Run(t, new(CmdTestSuite))
}

type CmdTestSuite struct {
	suite.Suite
}

// MockedManager mocked util.IClient
type MockedManager struct {
	mock.Mock
}

// GetAccountSettings ..
func (m *MockedManager) GetAccountSettings(input *fc.GetAccountSettingsInput) (*fc.GetAccountSettingsOutput, error) {
	args := m.Called(input)
	if args.Error(1) != nil {
		return nil, fmt.Errorf("error")
	}
	data := args.String(0)
	output := &fc.GetAccountSettingsOutput{}
	err := json.Unmarshal([]byte(data), output)
	if err != nil {
		return nil, err
	}
	return output, nil
}

// GetService ..
func (m *MockedManager) GetService(input *fc.GetServiceInput) (*fc.GetServiceOutput, error) {
	args := m.Called(input)
	if args.Error(1) != nil {
		return nil, fmt.Errorf("error")
	}
	data := args.String(0)
	output := &fc.GetServiceOutput{}
	err := json.Unmarshal([]byte(data), output)
	if err != nil {
		return nil, err
	}
	return output, nil
}

// ListServices ..
func (m *MockedManager) ListServices(input *fc.ListServicesInput) (*fc.ListServicesOutput, error) {
	args := m.Called(input)
	if args.Error(1) != nil {
		return nil, fmt.Errorf("error")
	}
	data := args.String(0)
	output := &fc.ListServicesOutput{}
	err := json.Unmarshal([]byte(data), output)
	if err != nil {
		return nil, err
	}
	return output, nil
}

// UpdateService ..
func (m *MockedManager) UpdateService(input *fc.UpdateServiceInput) (*fc.UpdateServiceOutput, error) {
	args := m.Called(input)
	if args.Error(1) != nil {
		return nil, fmt.Errorf("error")
	}
	data := args.String(0)
	output := &fc.UpdateServiceOutput{}
	err := json.Unmarshal([]byte(data), output)
	if err != nil {
		return nil, err
	}
	return output, nil
}

// CreateService ..
func (m *MockedManager) CreateService(input *fc.CreateServiceInput) (*fc.CreateServiceOutput, error) {
	args := m.Called(input)
	if args.Error(1) != nil {
		return nil, fmt.Errorf("error")
	}
	data := args.String(0)
	output := &fc.CreateServiceOutput{}
	err := json.Unmarshal([]byte(data), output)
	if err != nil {
		return nil, err
	}
	return output, nil
}

// DeleteService ..
func (m *MockedManager) DeleteService(input *fc.DeleteServiceInput) (*fc.DeleteServiceOutput, error) {
	args := m.Called(input)
	if args.Error(1) != nil {
		return nil, fmt.Errorf("error")
	}
	data := args.String(0)
	output := &fc.DeleteServiceOutput{}
	err := json.Unmarshal([]byte(data), output)
	if err != nil {
		return nil, err
	}
	return output, nil
}

// PublishServiceVersion ..
func (m *MockedManager) PublishServiceVersion(input *fc.PublishServiceVersionInput) (*fc.PublishServiceVersionOutput, error) {
	args := m.Called(input)
	if args.Error(1) != nil {
		return nil, fmt.Errorf("error")
	}
	data := args.String(0)
	output := &fc.PublishServiceVersionOutput{}
	err := json.Unmarshal([]byte(data), output)
	if err != nil {
		return nil, err
	}
	return output, nil
}

// ListServiceVersions ..
func (m *MockedManager) ListServiceVersions(input *fc.ListServiceVersionsInput) (*fc.ListServiceVersionsOutput, error) {
	args := m.Called(input)
	if args.Error(1) != nil {
		return nil, fmt.Errorf("error")
	}
	data := args.String(0)
	output := &fc.ListServiceVersionsOutput{}
	err := json.Unmarshal([]byte(data), output)
	if err != nil {
		return nil, err
	}
	return output, nil
}

// DeleteServiceVersion ..
func (m *MockedManager) DeleteServiceVersion(input *fc.DeleteServiceVersionInput) (*fc.DeleteServiceVersionOutput, error) {
	args := m.Called(input)
	if args.Error(1) != nil {
		return nil, fmt.Errorf("error")
	}
	data := args.String(0)
	output := &fc.DeleteServiceVersionOutput{}
	err := json.Unmarshal([]byte(data), output)
	if err != nil {
		return nil, err
	}
	return output, nil
}

// CreateAlias ..
func (m *MockedManager) CreateAlias(input *fc.CreateAliasInput) (*fc.CreateAliasOutput, error) {
	args := m.Called(input)
	if args.Error(1) != nil {
		return nil, fmt.Errorf("error")
	}
	data := args.String(0)
	output := &fc.CreateAliasOutput{}
	err := json.Unmarshal([]byte(data), output)
	if err != nil {
		return nil, err
	}
	return output, nil
}

// UpdateAlias ..
func (m *MockedManager) UpdateAlias(input *fc.UpdateAliasInput) (*fc.UpdateAliasOutput, error) {
	args := m.Called(input)
	if args.Error(1) != nil {
		return nil, fmt.Errorf("error")
	}
	data := args.String(0)
	output := &fc.UpdateAliasOutput{}
	err := json.Unmarshal([]byte(data), output)
	if err != nil {
		return nil, err
	}
	return output, nil
}

// GetAlias ..
func (m *MockedManager) GetAlias(input *fc.GetAliasInput) (*fc.GetAliasOutput, error) {
	args := m.Called(input)
	if args.Error(1) != nil {
		return nil, fmt.Errorf("error")
	}
	data := args.String(0)
	output := &fc.GetAliasOutput{}
	err := json.Unmarshal([]byte(data), output)
	if err != nil {
		return nil, err
	}
	return output, nil
}

// ListAliases ..
func (m *MockedManager) ListAliases(input *fc.ListAliasesInput) (*fc.ListAliasesOutput, error) {
	args := m.Called(input)
	if args.Error(1) != nil {
		return nil, fmt.Errorf("error")
	}
	data := args.String(0)
	output := &fc.ListAliasesOutput{}
	err := json.Unmarshal([]byte(data), output)
	if err != nil {
		return nil, err
	}
	return output, nil
}

// DeleteAlias ..
func (m *MockedManager) DeleteAlias(input *fc.DeleteAliasInput) (*fc.DeleteAliasOutput, error) {
	args := m.Called(input)
	if args.Error(1) != nil {
		return nil, fmt.Errorf("error")
	}
	data := args.String(0)
	output := &fc.DeleteAliasOutput{}
	err := json.Unmarshal([]byte(data), output)
	if err != nil {
		return nil, err
	}
	return output, nil
}

// CreateFunction ..
func (m *MockedManager) CreateFunction(input *fc.CreateFunctionInput) (*fc.CreateFunctionOutput, error) {
	args := m.Called(input)
	if args.Error(1) != nil {
		return nil, fmt.Errorf("error")
	}
	data := args.String(0)
	output := &fc.CreateFunctionOutput{}
	err := json.Unmarshal([]byte(data), output)
	if err != nil {
		return nil, err
	}
	return output, nil
}

// DeleteFunction ..
func (m *MockedManager) DeleteFunction(input *fc.DeleteFunctionInput) (*fc.DeleteFunctionOutput, error) {
	args := m.Called(input)
	if args.Error(1) != nil {
		return nil, fmt.Errorf("error")
	}
	data := args.String(0)
	output := &fc.DeleteFunctionOutput{}
	err := json.Unmarshal([]byte(data), output)
	if err != nil {
		return nil, err
	}
	return output, nil
}

// GetFunction ..
func (m *MockedManager) GetFunction(input *fc.GetFunctionInput) (*fc.GetFunctionOutput, error) {
	args := m.Called(input)
	if args.Error(1) != nil {
		return nil, fmt.Errorf("error")
	}
	data := args.String(0)
	output := &fc.GetFunctionOutput{}
	err := json.Unmarshal([]byte(data), output)
	if err != nil {
		return nil, err
	}
	return output, nil
}

// GetFunctionCode ..
func (m *MockedManager) GetFunctionCode(input *fc.GetFunctionCodeInput) (*fc.GetFunctionCodeOutput, error) {
	args := m.Called(input)
	if args.Error(1) != nil {
		return nil, fmt.Errorf("error")
	}
	data := args.String(0)
	output := &fc.GetFunctionCodeOutput{}
	err := json.Unmarshal([]byte(data), output)
	if err != nil {
		return nil, err
	}
	return output, nil
}

// ListFunctions ..
func (m *MockedManager) ListFunctions(input *fc.ListFunctionsInput) (*fc.ListFunctionsOutput, error) {
	args := m.Called(input)
	if args.Error(1) != nil {
		return nil, fmt.Errorf("error")
	}
	data := args.String(0)
	output := &fc.ListFunctionsOutput{}
	err := json.Unmarshal([]byte(data), output)
	if err != nil {
		return nil, err
	}
	return output, nil
}

// UpdateFunction ..
func (m *MockedManager) UpdateFunction(input *fc.UpdateFunctionInput) (*fc.UpdateFunctionOutput, error) {
	args := m.Called(input)
	if args.Error(1) != nil {
		return nil, fmt.Errorf("error")
	}
	data := args.String(0)
	output := &fc.UpdateFunctionOutput{}
	err := json.Unmarshal([]byte(data), output)
	if err != nil {
		return nil, err
	}
	return output, nil
}

// CreateTrigger ..
func (m *MockedManager) CreateTrigger(input *fc.CreateTriggerInput) (*fc.CreateTriggerOutput, error) {
	args := m.Called(input)
	if args.Error(1) != nil {
		return nil, fmt.Errorf("error")
	}
	data := args.String(0)
	output := &fc.CreateTriggerOutput{}
	err := json.Unmarshal([]byte(data), output)
	if err != nil {
		return nil, err
	}
	return output, nil
}

// GetTrigger ..
func (m *MockedManager) GetTrigger(input *fc.GetTriggerInput) (*fc.GetTriggerOutput, error) {
	args := m.Called(input)
	if args.Error(1) != nil {
		return nil, fmt.Errorf("error")
	}
	data := args.String(0)
	output := &fc.GetTriggerOutput{}
	err := json.Unmarshal([]byte(data), output)
	if err != nil {
		return nil, err
	}
	return output, nil
}

// UpdateTrigger ..
func (m *MockedManager) UpdateTrigger(input *fc.UpdateTriggerInput) (*fc.UpdateTriggerOutput, error) {
	args := m.Called(input)
	if args.Error(1) != nil {
		return nil, fmt.Errorf("error")
	}
	data := args.String(0)
	output := &fc.UpdateTriggerOutput{}
	err := json.Unmarshal([]byte(data), output)
	if err != nil {
		return nil, err
	}
	return output, nil
}

// DeleteTrigger ..
func (m *MockedManager) DeleteTrigger(input *fc.DeleteTriggerInput) (*fc.DeleteTriggerOutput, error) {
	args := m.Called(input)
	if args.Error(1) != nil {
		return nil, fmt.Errorf("error")
	}
	data := args.String(0)
	output := &fc.DeleteTriggerOutput{}
	err := json.Unmarshal([]byte(data), output)
	if err != nil {
		return nil, err
	}
	return output, nil
}

// ListTriggers ..
func (m *MockedManager) ListTriggers(input *fc.ListTriggersInput) (*fc.ListTriggersOutput, error) {
	args := m.Called(input)
	if args.Error(1) != nil {
		return nil, fmt.Errorf("error")
	}
	data := args.String(0)
	output := &fc.ListTriggersOutput{}
	err := json.Unmarshal([]byte(data), output)
	if err != nil {
		return nil, err
	}
	return output, nil
}

// InvokeFunction ..
func (m *MockedManager) InvokeFunction(input *fc.InvokeFunctionInput) (*fc.InvokeFunctionOutput, error) {
	args := m.Called(input)
	if args.Error(1) != nil {
		return nil, fmt.Errorf("error")
	}
	data := args.String(0)
	output := &fc.InvokeFunctionOutput{}
	err := json.Unmarshal([]byte(data), output)
	if err != nil {
		return nil, err
	}
	return output, nil
}