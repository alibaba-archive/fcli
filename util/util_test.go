package util

import (
	"encoding/json"
	"fmt"
	"github.com/aliyun/aliyun-log-go-sdk"
	"github.com/aliyun/fc-go-sdk"
	"github.com/aliyun/fcli/ram"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gopkg.in/h2non/gock.v1"
	"os"
	"os/exec"
	"testing"
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

func (s *UtilTestSuite) TestCheckImageExist() {
	output = func(cmd *exec.Cmd) (bytes []byte, e error) {
		return []byte{}, nil
	}

	isExist := CheckImageExist("aliyunfc/runtime-nodejs6", "build")
	s.False(isExist)

	output = func(cmd *exec.Cmd) (bytes []byte, e error) {
		return []byte{1}, nil
	}
	isExist = CheckImageExist("aliyunfc/runtime-nodejs6", "build")
	s.True(isExist)
}

func (s *UtilTestSuite) TestGetPublicImageDigest() {
	defer gock.Off()
	gock.New("https://auth.docker.io").
		MatchParam("service", "registry.docker.io").
		MatchParam("scope", "repository:aliyunfc/runtime-nodejs6:pull").
		Get("token").
		Reply(200).
		JSON(map[string]string{"token": "testToken"})

	mockDigest := "sha256:123456"

	gock.New("https://registry-1.docker.io").
		MatchHeader("Authorization", "Bearer testToken").
		Get("v2/aliyunfc/runtime-nodejs6/manifests/build").
		Reply(200).
		AddHeader("Docker-Content-Digest", mockDigest)

	digest, _ := GetPublicImageDigest("aliyunfc/runtime-nodejs6", "build")
	s.Equal(mockDigest, digest)
}

func (s *UtilTestSuite) TestGetLocalImageDigest() {
	mockDigest := "sha256:123456"

	output = func(cmd *exec.Cmd) (bytes []byte, e error) {
		return []byte("'aliyunfc/runtime-nodejs6@" + mockDigest + "'"), nil
	}

	digest, _ := GetLocalImageDigest("aliyunfc/runtime-nodejs6", "build")
	s.Equal(mockDigest, digest)
}

// TestNewExecutableSubCmd ..
func (s *UtilTestSuite) TestNewExecutableSubCmd() {
	cmd := NewExecutableSubCmd("")
	s.NotNil(cmd)
}

// TestNewGlobalConfig ..
func (s *UtilTestSuite) TestNewGlobalConfig() {
	cfg := NewGlobalConfig()
	s.NotNil(cfg)
}

// TestNewFClient ..
func (s *UtilTestSuite) TestNewFClient() {
	cfg := NewGlobalConfig()
	s.NotNil(cfg)

	client, err := NewFClient(cfg)
	s.Nil(err)
	s.NotNil(client)
}

// TestNewRAMClient ..
func (s *UtilTestSuite) TestNewRAMClient() {
	accessKeyID := "testAccessKeyID"
	accessKeySecret := "testAccessKeySecret"

	client, err := NewRAMClient(accessKeyID, accessKeySecret)

	s.Nil(err)
	s.NotNil(client)
}

// TestNewSLSClient ..
func (s *UtilTestSuite) TestNewSLSClient() {
	cfg := NewGlobalConfig()
	s.NotNil(cfg)

	client := NewSLSClient(cfg)
	s.NotNil(client)
}

// TestPrettyPrintLog ..
func (s *UtilTestSuite) TestPrettyPrintLog() {
	logMap := map[string]string{
		"__time__":     "1545142890",
		"serviceName":  "testService",
		"functionName": "testFunction",
		"message":      "testMessage",
	}
	PrettyPrintLog(logMap)
}

// TestPrettyPrintLogs ..
func (s *UtilTestSuite) TestPrettyPrintLogs() {
	logMap := map[string]string{
		"__time__":     "1545142890",
		"serviceName":  "testService",
		"functionName": "testFunction",
		"message":      "testMessage",
	}
	logMaps := &[]map[string]string{
		logMap,
	}
	PrettyPrintLogs(logMaps, true)
	PrettyPrintLogs(logMaps, false)
}

// TestGetRegions ..
func (s *UtilTestSuite) TestGetRegions() {
	regions := GetRegions()
	s.NotNil(regions)
}

// TestGetRegionNoForEndpoint ..
func (s *UtilTestSuite) TestGetRegionNoForEndpoint() {
	regions := GetRegions()
	s.NotNil(regions)

	for _, region := range regions {
		s.Equal(region, GetRegionNoForEndpoint(region))
	}
}

// TestGetRegionNoForSLSEndpoint ..
func (s *UtilTestSuite) TestGetRegionNoForSLSEndpoint() {
	regions := GetSLSRegions()
	s.NotNil(regions)

	for _, region := range regions {
		s.Equal(region, GetRegionNoForSLSEndpoint(region))
	}
}

// TestGetUIDFromEndpoint ..
func (s *UtilTestSuite) TestGetUIDFromEndpoint() {
	endpoint := "https://123456789.cn-shanghai.fc.aliyuncs.com"
	uid, err := GetUIDFromEndpoint(endpoint)
	s.Equal(uid, "123456789")
	s.Nil(err)
}

// MockedManager mocked ram.RoleManager ram.PolicyManager
type MockedManager struct {
	mock.Mock
}

// CreateRole mocked ram.RoleManager create role
func (m *MockedManager) CreateRole(roleName, assumeRolePolicyDocument, decription string) (*ram.CreateRoleResponse, error) {
	args := m.Called(roleName, assumeRolePolicyDocument, decription)
	if !args.Bool(0) {
		return nil, fmt.Errorf("error")
	}
	data := `{"RequestID":"test","Role":{"Arn":"` + roleName + `"}}`
	resp := &ram.CreateRoleResponse{}
	err := json.Unmarshal([]byte(data), resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// CreatePolicy mocked ram.PolicyManager crate policy
func (m *MockedManager) CreatePolicy(policyName, policyDocument, description string) (*ram.CreatePolicyResponse, error) {
	args := m.Called(policyName, policyDocument, description)
	if !args.Bool(0) {
		return nil, fmt.Errorf("error")
	}
	data := `{"RequestID":"test"}`
	resp := &ram.CreatePolicyResponse{}
	err := json.Unmarshal([]byte(data), resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// CreatePolicyVersion mocked ram.PolicyManager crate policy version
func (m *MockedManager) CreatePolicyVersion(policyName, policyDocument, setAsDefault string) (*ram.PolicyVersionResponse, error) {
	args := m.Called(policyName, policyDocument, setAsDefault)
	if !args.Bool(0) {
		return nil, fmt.Errorf("error")
	}
	data := `{"RequestID":"test"}`
	resp := &ram.PolicyVersionResponse{}
	err := json.Unmarshal([]byte(data), resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// GetPolicyVersion mocked ram.PolicyManager get policy version
func (m *MockedManager) GetPolicyVersion(policyName, policyType, VersionID string) (*ram.PolicyVersionResponse, error) {
	args := m.Called(policyName, policyType, VersionID)
	if !args.Bool(0) {
		return nil, fmt.Errorf("error")
	}
	data := `{"RequestID":"test","PolicyVersion":{"VersionId":"1"}}`
	resp := &ram.PolicyVersionResponse{}
	err := json.Unmarshal([]byte(data), resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// ListPolicyVersions mocked ram.PolicyManager list policy versions
func (m *MockedManager) ListPolicyVersions(policyName string, policyType string) (*ram.ListPolicyVersionsResponse, error) {
	args := m.Called(policyName, policyType)
	if !args.Bool(0) {
		return nil, fmt.Errorf("error")
	}
	data := `{"RequestID":"test","PolicyVersions":{"PolicyVersion":[{"VersionID":"1","IsDefaultVersion":true}]}}`
	resp := &ram.ListPolicyVersionsResponse{}
	err := json.Unmarshal([]byte(data), resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// AttachPolicyToRole mocked ram.AttachPolicyToRole attach policy to role
func (m *MockedManager) AttachPolicyToRole(policyType, policyName, roleName string) (*ram.AttachPolicyToRoleResponse, error) {
	data := `{"RequestID":"1"}`
	resp := &ram.AttachPolicyToRoleResponse{}
	err := json.Unmarshal([]byte(data), resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// GetLogs mocked LogManager get logs
func (m *MockedManager) GetLogs(topic string, from int64, to int64, queryExp string,
	maxLineNum int64, offset int64, reverse bool) (*sls.GetLogsResponse, error) {
	data := ``
	if offset < 100 {
		data = `{"count":101,"Logs":[]}`
	} else {
		data = `{"count":0,"Logs":[]}`
	}
	resp := &sls.GetLogsResponse{}
	err := json.Unmarshal([]byte(data), resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// TestCreateRole ..
func (s *UtilTestSuite) TestCreateRole() {
	mockObject := &MockedManager{}
	mockObject.
		On("CreateRole", "correctRoleName", mock.Anything, mock.Anything).
		Return(true)
	roleArn, err := CreateRole(mockObject, "correctRoleName", "")
	s.Nil(err)
	s.Equal("correctRoleName", roleArn)

	mockObject.
		On("CreateRole", "errorRoleName", mock.Anything, mock.Anything).
		Return(false)
	roleArn, err = CreateRole(mockObject, "errorRoleName", "")
	s.NotNil(err)
	s.Equal("", roleArn)
}

// TestCreatePolicy ..
func (s *UtilTestSuite) TestCreatePolicy() {
	mockObject := &MockedManager{}
	mockObject.
		On("CreatePolicy", "correctPolicyName", mock.Anything, mock.Anything).
		Return(true)
	err := CreatePolicy(mockObject, "correctPolicyName", "", "")
	s.Nil(err)

	mockObject.
		On("CreatePolicy", "errorPolicyName", mock.Anything, mock.Anything).
		Return(false)
	err = CreatePolicy(mockObject, "errorPolicyName", "", "")
	s.NotNil(err)
}

// TestCreatePolicyVersion ..
func (s *UtilTestSuite) TestCreatePolicyVersion() {
	mockObject := &MockedManager{}
	mockObject.
		On("CreatePolicyVersion", "correctPolicyName", mock.Anything, mock.Anything).
		Return(true)
	err := CreatePolicyVersion(mockObject, "correctPolicyName", "", "")
	s.Nil(err)

	mockObject.
		On("CreatePolicyVersion", "errorPolicyName", mock.Anything, mock.Anything).
		Return(false)
	err = CreatePolicyVersion(mockObject, "errorPolicyName", "", "")
	s.NotNil(err)
}

// TestGetDefaultPolicyVersion ..
func (s *UtilTestSuite) TestGetDefaultPolicyVersion() {
	mockObject := &MockedManager{}
	mockObject.
		On("ListPolicyVersions", "correctPolicyName", mock.Anything).
		Return(true)
	mockObject.
		On("GetPolicyVersion", "correctPolicyName", mock.Anything, mock.Anything).
		Return(true)
	policyVersion, err := GetDefaultPolicyVersion(mockObject, "correctPolicyName", "")
	s.Nil(err)
	s.NotNil(*policyVersion)

	mockObject.
		On("ListPolicyVersions", "errorPolicyName", mock.Anything).
		Return(true)
	mockObject.
		On("GetPolicyVersion", "errorPolicyName", mock.Anything, mock.Anything).
		Return(false)
	policyVersion, err = GetDefaultPolicyVersion(mockObject, "errorPolicyName", "")
	s.NotNil(err)
	s.Nil(policyVersion)
}

// TestCheckPolicyResourcePermission ..
func (s *UtilTestSuite) TestCheckPolicyResourcePermission() {
	data := `{"VersionID":"1","PolicyDocument":"{"}`
	policyVersion := &ram.PolicyVersion{}
	err := json.Unmarshal([]byte(data), policyVersion)
	s.Nil(err)

	resp, err := CheckPolicyResourcePermission(policyVersion, "")
	s.NotNil(err)
	s.False(resp)

	data = `{"VersionID":"1","PolicyDocument":"{\"Statement\":[{\"Resource\":[\"testResource\"]}]}"}`
	policyVersion = &ram.PolicyVersion{}
	err = json.Unmarshal([]byte(data), policyVersion)
	s.Nil(err)

	resp, err = CheckPolicyResourcePermission(policyVersion, "[\"testResource\"]")
	s.Nil(err)
	s.True(resp)
}

// TestAttachPolicy ..
func (s *UtilTestSuite) TestAttachPolicy() {
	mockObject := &MockedManager{}
	err := AttachPolicy(mockObject, "testPolicy", "")
	s.Nil(err)
}

// TestGetBindingCmd ..
func (s *UtilTestSuite) TestGetBindingCmd() {
	cmd := GetBindingCmd("test", "", "")
	s.NotNil(cmd)
}

// TestParseAdditionalVersionWeight ..
func (s *UtilTestSuite) TestParseAdditionalVersionWeight() {
	routes := []string{"1=0.01", "2=0.02"}
	additionalVersionWeight := ParseAdditionalVersionWeight(routes)
	s.NotNil(additionalVersionWeight)
	s.Equal(0.01, additionalVersionWeight["1"])
	s.Equal(0.02, additionalVersionWeight["2"])

	routes = []string{"1=."}
	additionalVersionWeight = ParseAdditionalVersionWeight(routes)
	s.Nil(additionalVersionWeight)
}

// TestGetTriggerConfig
func (s *UtilTestSuite) TestGetTriggerConfig() {
	errTriggerConfig, err := GetTriggerConfig(fc.TRIGGER_TYPE_OSS, "")
	s.NotNil(err)
	s.Nil(errTriggerConfig)

	// oss
	ossTriggerConfig, err := GetTriggerConfig(fc.TRIGGER_TYPE_OSS, "../example/ossTriggerConfig.yaml")
	s.Nil(err)
	s.NotNil(ossTriggerConfig)

	// timer
	timerTriggerConfig, err := GetTriggerConfig(fc.TRIGGER_TYPE_TIMER, "../example/timerTriggerConfig.yaml")
	s.Nil(err)
	s.NotNil(timerTriggerConfig)

	// log
	logTriggerConfig, err := GetTriggerConfig(fc.TRIGGER_TYPE_LOG, "../example/logTriggerConfig.yaml")
	s.Nil(err)
	s.NotNil(logTriggerConfig)

	// cdn_events
	cdnEventTriggerConfig, err := GetTriggerConfig(fc.TRIGGER_TYPE_CDN_EVENTS, "../example/cdnEventTriggerConfig.yaml")
	s.Nil(err)
	s.NotNil(cdnEventTriggerConfig)

	// http
	httpTriggerConfig, err := GetTriggerConfig(fc.TRIGGER_TYPE_HTTP, "../example/httpTriggerConfig.yaml")
	s.Nil(err)
	s.NotNil(httpTriggerConfig)

	// mns_topic
	mnsTopicTriggerConfig, err := GetTriggerConfig(fc.TRIGGER_TYPE_MNS_TOPIC, "../example/mnsTopicTriggerConfig.yaml")
	s.Nil(err)
	s.NotNil(mnsTopicTriggerConfig)

	//default
	defaultConfig, err := GetTriggerConfig("default", "../example/mnsTopicTriggerConfig.yaml")
	s.NotNil(err)
	s.Nil(defaultConfig)
}

// TestGetAllLogsWithinTimeRange ..
func (s *UtilTestSuite) TestGetAllLogsWithinTimeRange() {
	mockObject := &MockedManager{}
	err := GetAllLogsWithinTimeRange(mockObject, "demo", "testFunc", 1548065021, 1548065081)
	s.Nil(err)
}

// TestDownloadFromURL ..
func (s *UtilTestSuite) TestDownloadFromURL() {
	defer gock.Off()
	gock.New("https://www.download.com").
		Reply(200).
		JSON(map[string]string{"token": "testToken"})

	err := DownloadFromURL("https://www.download.com", os.Stdout)
	s.Nil(err)
}
