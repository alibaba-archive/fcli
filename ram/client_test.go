package ram

import (
	"fmt"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gopkg.in/h2non/gock.v1"
	"net/http"
	"net/url"
	"testing"
)

func TestClient(t *testing.T) {
	suite.Run(t, new(ClientTestSuite))
}

type ClientTestSuite struct {
	suite.Suite
}

// TestSetRetryTimes ..
func (s *ClientTestSuite) TestSetRetryTimes() {
	clientOption := &ClientOption{}
	clientOption.setRetryTimes(1)
	s.Equal(int32(1), clientOption.retryTimes)
}

// TestGetDefaultClientOption ..
func (s *ClientTestSuite) TestGetDefaultClientOption() {
	clientOption := getDefaultClientOption()
	s.NotNil(clientOption)
}

// TestGetWaitIntervalInMS ..
func (s *ClientTestSuite) TestGetWaitIntervalInMS() {
	retry := getWaitIntervalInMS(0)
	s.Equal(int32(0), retry)

	retry = getWaitIntervalInMS(1)
	s.Equal(int32(200), retry)
}

// TestNewClient ..
func (s *ClientTestSuite) TestNewClient() {
	client, err := NewClient("testEndpoint", "testId", "testSecret")
	s.Nil(err)
	s.NotNil(client)
}

// TestWithRetryTimes ..
func (s *ClientTestSuite) TestWithRetryTimes() {
	client := &Client{option: &ClientOption{}}
	client = client.WithRetryTimes(1)
	s.NotNil(client)
	s.NotNil(client.option)
	s.Equal(int32(1), client.option.retryTimes)
}

// TestGetCommontParam ..
func (s *ClientTestSuite) TestGetCommontParam() {
	client := &Client{}
	params := client.getCommontParam()
	s.NotNil(params)
}

// MockedRoleManager mocked ram.RequestManager
type MockedManager struct {
	mock.Mock
}

// CreateRole mocked ram.RoleManager create role
func (m *MockedManager) sendRequest(rawParams map[string]string) ([]byte, error) {
	args := m.Called(rawParams)
	if args.Error(1) != nil {
		return nil, fmt.Errorf("error")
	}
	return []byte(args.String(0)), nil
}

// TestCreateRole ..
func (s *ClientTestSuite) TestCreateRole() {
	// error
	mockObject := &MockedManager{}
	mockObject.
		On("sendRequest", mock.Anything).
		Return("", fmt.Errorf(""))
	client := &Client{
		RequestManager: mockObject,
	}

	resp, err := client.CreateRole("roleName", "", "")
	s.NotNil(err)
	s.Nil(resp)

	// error json string
	mockObject = &MockedManager{}
	mockObject.
		On("sendRequest", mock.Anything).
		Return(`{"RequestId:1}`, nil)
	client = &Client{
		RequestManager: mockObject,
	}

	resp, err = client.CreateRole("roleName", "", "")
	s.NotNil(err)
	s.Nil(resp)

	// success
	mockObject = &MockedManager{}
	mockObject.
		On("sendRequest", mock.Anything).
		Return(`{"RequestId":"1"}`, nil)
	client = &Client{
		RequestManager: mockObject,
	}

	resp, err = client.CreateRole("roleName", "", "")
	s.Nil(err)
	s.NotNil(resp)
}

// TestGetRole ..
func (s *ClientTestSuite) TestGetRole() {
	// error
	mockObject := &MockedManager{}
	mockObject.
		On("sendRequest", mock.Anything).
		Return("", fmt.Errorf(""))
	client := &Client{
		RequestManager: mockObject,
	}

	resp, err := client.GetRole("roleName")
	s.NotNil(err)
	s.Nil(resp)

	// error json string
	mockObject = &MockedManager{}
	mockObject.
		On("sendRequest", mock.Anything).
		Return(`{"RequestId":"1"`, nil)
	client = &Client{
		RequestManager: mockObject,
	}

	resp, err = client.GetRole("roleName")
	s.NotNil(err)
	s.Nil(resp)

	// success
	mockObject = &MockedManager{}
	mockObject.
		On("sendRequest", mock.Anything).
		Return(`{"RequestId":"1"}`, nil)
	client = &Client{
		RequestManager: mockObject,
	}

	resp, err = client.GetRole("roleName")
	s.Nil(err)
	s.NotNil(resp)
}

// TestDeleteRole ..
func (s *ClientTestSuite) TestDeleteRole() {
	// error
	mockObject := &MockedManager{}
	mockObject.
		On("sendRequest", mock.Anything).
		Return("", fmt.Errorf(""))
	client := &Client{
		RequestManager: mockObject,
	}

	resp, err := client.DeleteRole("roleName")
	s.NotNil(err)
	s.Nil(resp)

	// error json string
	mockObject = &MockedManager{}
	mockObject.
		On("sendRequest", mock.Anything).
		Return(`{"RequestId":"1"`, nil)
	client = &Client{
		RequestManager: mockObject,
	}

	resp, err = client.DeleteRole("roleName")
	s.NotNil(err)
	s.Nil(resp)

	// success
	mockObject = &MockedManager{}
	mockObject.
		On("sendRequest", mock.Anything).
		Return(`{"RequestId":"1"}`, nil)
	client = &Client{
		RequestManager: mockObject,
	}

	resp, err = client.DeleteRole("roleName")
	s.Nil(err)
	s.NotNil(resp)
}

// TestListRoles ..
func (s *ClientTestSuite) TestListRoles() {
	// error
	mockObject := &MockedManager{}
	mockObject.
		On("sendRequest", mock.Anything).
		Return("", fmt.Errorf(""))
	client := &Client{
		RequestManager: mockObject,
	}

	resp, err := client.ListRoles()
	s.NotNil(err)
	s.Nil(resp)

	// error json string
	mockObject = &MockedManager{}
	mockObject.
		On("sendRequest", mock.Anything).
		Return(`{"RequestId":"1"`, nil)
	client = &Client{
		RequestManager: mockObject,
	}

	resp, err = client.ListRoles()
	s.NotNil(err)
	s.Nil(resp)

	// success
	mockObject = &MockedManager{}
	mockObject.
		On("sendRequest", mock.Anything).
		Return(`{"RequestId":"1"}`, nil)
	client = &Client{
		RequestManager: mockObject,
	}

	resp, err = client.ListRoles()
	s.Nil(err)
	s.NotNil(resp)
}

// TestCreatePolicy ..
func (s *ClientTestSuite) TestCreatePolicy() {
	// error
	mockObject := &MockedManager{}
	mockObject.
		On("sendRequest", mock.Anything).
		Return("", fmt.Errorf(""))
	client := &Client{
		RequestManager: mockObject,
	}

	resp, err := client.CreatePolicy("policyName", "", "")
	s.NotNil(err)
	s.Nil(resp)

	// error json string
	mockObject = &MockedManager{}
	mockObject.
		On("sendRequest", mock.Anything).
		Return(`{"RequestId":"1"`, nil)
	client = &Client{
		RequestManager: mockObject,
	}

	resp, err = client.CreatePolicy("policyName", "", "")
	s.NotNil(err)
	s.Nil(resp)

	// success
	mockObject = &MockedManager{}
	mockObject.
		On("sendRequest", mock.Anything).
		Return(`{"RequestId":"1"}`, nil)
	client = &Client{
		RequestManager: mockObject,
	}

	resp, err = client.CreatePolicy("policyName", "", "")
	s.Nil(err)
	s.NotNil(resp)
}

// TestGetPolicy ..
func (s *ClientTestSuite) TestGetPolicy() {
	// error
	mockObject := &MockedManager{}
	mockObject.
		On("sendRequest", mock.Anything).
		Return("", fmt.Errorf(""))
	client := &Client{
		RequestManager: mockObject,
	}

	resp, err := client.GetPolicy("policyName", "")
	s.NotNil(err)
	s.Nil(resp)

	// error json string
	mockObject = &MockedManager{}
	mockObject.
		On("sendRequest", mock.Anything).
		Return(`{"RequestId":"1"`, nil)
	client = &Client{
		RequestManager: mockObject,
	}

	resp, err = client.GetPolicy("policyName", "")
	s.NotNil(err)
	s.Nil(resp)

	// success
	mockObject = &MockedManager{}
	mockObject.
		On("sendRequest", mock.Anything).
		Return(`{"RequestId":"1"}`, nil)
	client = &Client{
		RequestManager: mockObject,
	}

	resp, err = client.GetPolicy("policyName", "")
	s.Nil(err)
	s.NotNil(resp)
}

// TestListPolicyVersions ..
func (s *ClientTestSuite) TestListPolicyVersions() {
	// error
	mockObject := &MockedManager{}
	mockObject.
		On("sendRequest", mock.Anything).
		Return("", fmt.Errorf(""))
	client := &Client{
		RequestManager: mockObject,
	}

	resp, err := client.ListPolicyVersions("policyName", "")
	s.NotNil(err)
	s.Nil(resp)

	// error json string
	mockObject = &MockedManager{}
	mockObject.
		On("sendRequest", mock.Anything).
		Return(`{"RequestId":"1"`, nil)
	client = &Client{
		RequestManager: mockObject,
	}

	resp, err = client.ListPolicyVersions("policyName", "")
	s.NotNil(err)
	s.Nil(resp)

	// success
	mockObject = &MockedManager{}
	mockObject.
		On("sendRequest", mock.Anything).
		Return(`{"RequestId":"1"}`, nil)
	client = &Client{
		RequestManager: mockObject,
	}

	resp, err = client.ListPolicyVersions("policyName", "")
	s.Nil(err)
	s.NotNil(resp)
}

// TestGetPolicyVersion ..
func (s *ClientTestSuite) TestGetPolicyVersion() {
	// error
	mockObject := &MockedManager{}
	mockObject.
		On("sendRequest", mock.Anything).
		Return("", fmt.Errorf(""))
	client := &Client{
		RequestManager: mockObject,
	}

	resp, err := client.GetPolicyVersion("policyName", "", "1")
	s.NotNil(err)
	s.Nil(resp)

	// error json string
	mockObject = &MockedManager{}
	mockObject.
		On("sendRequest", mock.Anything).
		Return(`{"RequestId":"1"`, nil)
	client = &Client{
		RequestManager: mockObject,
	}

	resp, err = client.GetPolicyVersion("policyName", "", "1")
	s.NotNil(err)
	s.Nil(resp)

	// success
	mockObject = &MockedManager{}
	mockObject.
		On("sendRequest", mock.Anything).
		Return(`{"RequestId":"1"}`, nil)
	client = &Client{
		RequestManager: mockObject,
	}

	resp, err = client.GetPolicyVersion("policyName", "", "1")
	s.Nil(err)
	s.NotNil(resp)
}

// TestCreatePolicyVersion ..
func (s *ClientTestSuite) TestCreatePolicyVersion() {
	// error
	mockObject := &MockedManager{}
	mockObject.
		On("sendRequest", mock.Anything).
		Return("", fmt.Errorf(""))
	client := &Client{
		RequestManager: mockObject,
	}

	resp, err := client.CreatePolicyVersion("policyName", "", "")
	s.NotNil(err)
	s.Nil(resp)

	// error json string
	mockObject = &MockedManager{}
	mockObject.
		On("sendRequest", mock.Anything).
		Return(`{"RequestId":"1"`, nil)
	client = &Client{
		RequestManager: mockObject,
	}

	resp, err = client.CreatePolicyVersion("policyName", "", "")
	s.NotNil(err)
	s.Nil(resp)

	// success
	mockObject = &MockedManager{}
	mockObject.
		On("sendRequest", mock.Anything).
		Return(`{"RequestId":"1"}`, nil)
	client = &Client{
		RequestManager: mockObject,
	}

	resp, err = client.CreatePolicyVersion("policyName", "", "")
	s.Nil(err)
	s.NotNil(resp)
}

// TestDeletePolicy ..
func (s *ClientTestSuite) TestDeletePolicy() {
	// error
	mockObject := &MockedManager{}
	mockObject.
		On("sendRequest", mock.Anything).
		Return("", fmt.Errorf(""))
	client := &Client{
		RequestManager: mockObject,
	}

	resp, err := client.DeletePolicy("policyName")
	s.NotNil(err)
	s.Nil(resp)

	// error json string
	mockObject = &MockedManager{}
	mockObject.
		On("sendRequest", mock.Anything).
		Return(`{"RequestId":"1"`, nil)
	client = &Client{
		RequestManager: mockObject,
	}

	resp, err = client.DeletePolicy("policyName")
	s.NotNil(err)
	s.Nil(resp)

	// success
	mockObject = &MockedManager{}
	mockObject.
		On("sendRequest", mock.Anything).
		Return(`{"RequestId":"1"}`, nil)
	client = &Client{
		RequestManager: mockObject,
	}

	resp, err = client.DeletePolicy("policyName")
	s.Nil(err)
	s.NotNil(resp)
}

// TestListPolicies ..
func (s *ClientTestSuite) TestListPolicies() {
	// error
	mockObject := &MockedManager{}
	mockObject.
		On("sendRequest", mock.Anything).
		Return("", fmt.Errorf(""))
	client := &Client{
		RequestManager: mockObject,
	}

	resp, err := client.ListPolicies()
	s.NotNil(err)
	s.Nil(resp)

	// error json string
	mockObject = &MockedManager{}
	mockObject.
		On("sendRequest", mock.Anything).
		Return(`{"RequestId":"1"`, nil)
	client = &Client{
		RequestManager: mockObject,
	}

	resp, err = client.ListPolicies()
	s.NotNil(err)
	s.Nil(resp)

	// success
	mockObject = &MockedManager{}
	mockObject.
		On("sendRequest", mock.Anything).
		Return(`{"RequestId":"1"}`, nil)
	client = &Client{
		RequestManager: mockObject,
	}

	resp, err = client.ListPolicies()
	s.Nil(err)
	s.NotNil(resp)
}

// TestListPoliciesForRole ..
func (s *ClientTestSuite) TestListPoliciesForRole() {
	// error
	mockObject := &MockedManager{}
	mockObject.
		On("sendRequest", mock.Anything).
		Return("", fmt.Errorf(""))
	client := &Client{
		RequestManager: mockObject,
	}

	resp, err := client.ListPoliciesForRole("roleName")
	s.NotNil(err)
	s.Nil(resp)

	// error json string
	mockObject = &MockedManager{}
	mockObject.
		On("sendRequest", mock.Anything).
		Return(`{"RequestId":"1"`, nil)
	client = &Client{
		RequestManager: mockObject,
	}

	resp, err = client.ListPoliciesForRole("roleName")
	s.NotNil(err)
	s.Nil(resp)

	// success
	mockObject = &MockedManager{}
	mockObject.
		On("sendRequest", mock.Anything).
		Return(`{"RequestId":"1"}`, nil)
	client = &Client{
		RequestManager: mockObject,
	}

	resp, err = client.ListPoliciesForRole("roleName")
	s.Nil(err)
	s.NotNil(resp)
}

// TestAttachPolicyToRole ..
func (s *ClientTestSuite) TestAttachPolicyToRole() {
	// error
	mockObject := &MockedManager{}
	mockObject.
		On("sendRequest", mock.Anything).
		Return("", fmt.Errorf(""))
	client := &Client{
		RequestManager: mockObject,
	}

	resp, err := client.AttachPolicyToRole("", "policyName", "roleName")
	s.NotNil(err)
	s.Nil(resp)

	// error json string
	mockObject = &MockedManager{}
	mockObject.
		On("sendRequest", mock.Anything).
		Return(`{"RequestId":"1"`, nil)
	client = &Client{
		RequestManager: mockObject,
	}

	resp, err = client.AttachPolicyToRole("", "policyName", "roleName")
	s.NotNil(err)
	s.Nil(resp)

	// success
	mockObject = &MockedManager{}
	mockObject.
		On("sendRequest", mock.Anything).
		Return(`{"RequestId":"1"}`, nil)
	client = &Client{
		RequestManager: mockObject,
	}

	resp, err = client.AttachPolicyToRole("", "policyName", "roleName")
	s.Nil(err)
	s.NotNil(resp)
}

// TestDetachPolicyFromRole ..
func (s *ClientTestSuite) TestDetachPolicyFromRole() {
	// error
	mockObject := &MockedManager{}
	mockObject.
		On("sendRequest", mock.Anything).
		Return("", fmt.Errorf(""))
	client := &Client{
		RequestManager: mockObject,
	}

	resp, err := client.DetachPolicyFromRole("", "policyName", "roleName")
	s.NotNil(err)
	s.Nil(resp)

	// error json string
	mockObject = &MockedManager{}
	mockObject.
		On("sendRequest", mock.Anything).
		Return(`{"RequestId":"1"`, nil)
	client = &Client{
		RequestManager: mockObject,
	}

	resp, err = client.DetachPolicyFromRole("", "policyName", "roleName")
	s.NotNil(err)
	s.Nil(resp)

	// success
	mockObject = &MockedManager{}
	mockObject.
		On("sendRequest", mock.Anything).
		Return(`{"RequestId":"1"}`, nil)
	client = &Client{
		RequestManager: mockObject,
	}

	resp, err = client.DetachPolicyFromRole("", "policyName", "roleName")
	s.Nil(err)
	s.NotNil(resp)
}

// TestSendRequest ..
func (s *ClientTestSuite) TestSendRequest() {
	defer gock.Off()

	client := &Client{
		accessKeyID: "testKey",
		accessKeySecret: "testSecret",
		hclient: &http.Client{},
		endpoint: &url.URL{
			Scheme: "https",
			Host: "ram.aliyuncs.com",
		},
		option: &ClientOption{},
	}

	// success
	gock.New("https://ram.aliyuncs.com").
		Reply(200).
		JSON(map[string]string{"token": "testToken"})
	resp, err := client.sendRequest(make(map[string]string))
	s.Nil(err)
	s.NotNil(resp)

	// success 500 retry
	gock.New("https://ram.aliyuncs.com").
		Reply(500).
		JSON(map[string]string{})
	resp, err = client.sendRequest(make(map[string]string))
	s.NotNil(err)
	s.Nil(resp)

	// success 500 error JSON
	gock.New("https://ram.aliyuncs.com").
		Reply(500).
		BodyString(`{"key":"1}`)
	resp, err = client.sendRequest(make(map[string]string))
	s.NotNil(err)
	s.Nil(resp)

}