package ram

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"gopkg.in/matryer/try.v1"
	"io/ioutil"
	"math"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/satori/go.uuid"
)

// RoleManager define role operation
type RoleManager interface {
	CreateRole(roleName, assumeRolePolicyDocument, decription string) (*CreateRoleResponse, error)
}

// PolicyManager define policy operation
type PolicyManager interface {
	CreatePolicy(policyName, policyDocument, description string) (*CreatePolicyResponse, error)
	CreatePolicyVersion(policyName, policyDocument, setAsDefault string) (*PolicyVersionResponse, error)
	GetPolicyVersion(policyName, policyType, VersionID string) (*PolicyVersionResponse, error)
	ListPolicyVersions(policyName string, policyType string) (*ListPolicyVersionsResponse, error)
	AttachPolicyToRole(policyType, policyName, roleName string) (*AttachPolicyToRoleResponse, error)
}

type RequestManager interface {
	sendRequest(rawParams map[string]string) ([]byte, error)
}

// ClientOption define sdk option config
type ClientOption struct {
	retryTimes int32
}

// Client define RAM Client
type Client struct {
	accessKeyID     string
	accessKeySecret string
	endpoint        *url.URL
	option          *ClientOption
	hclient         *http.Client
	RequestManager
}

func (option *ClientOption) setRetryTimes(retry int32) {
	option.retryTimes = retry
}

func getDefaultClientOption() *ClientOption {
	return &ClientOption{
		retryTimes: DefaultRetryTimes,
	}
}

func getWaitIntervalInMS(retry int32) int32 {
	if retry == int32(0) {
		return int32(0)
	}
	interval := math.Exp2(float64(retry)) * 100
	return int32(math.Min(interval, float64(DefaultMaxWaitIntervalInMS)))
}

// NewClient new RAM client
func NewClient(endpoint, accessKeyID, accessKeySecret string) (client *Client, err error) {
	client = &Client{}
	client.endpoint, err = url.Parse(endpoint)
	if err != nil {
		return nil, err
	}

	client.accessKeyID = accessKeyID
	client.accessKeySecret = accessKeySecret
	client.hclient = &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			Dial: (&net.Dialer{
				Timeout:   3 * time.Second, //设置DialTimeOut,默认3秒
				KeepAlive: 30 * time.Second,
			}).Dial,
			// 保持5个空闲连接
			MaxIdleConnsPerHost: 5,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
				CipherSuites: []uint16{
					tls.TLS_RSA_WITH_RC4_128_SHA,
					tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
					tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
					tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
					tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
					tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
					tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
					tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
					tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
					tls.TLS_RSA_WITH_AES_128_CBC_SHA,
					tls.TLS_RSA_WITH_AES_256_CBC_SHA,
					tls.TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA,
					tls.TLS_RSA_WITH_3DES_EDE_CBC_SHA,
				},
			},
		},
		Timeout: 3 * time.Second,
	}
	client.option = getDefaultClientOption()

	return
}

// WithRetryTimes : 进行可重入错误重试的次数, 目前对url.Error和500以及503错误进行重试
func (c *Client) WithRetryTimes(retry int32) *Client {
	c.option.setRetryTimes(retry)
	return c
}

func (c *Client) getCommontParam() map[string]string {
	u, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}
	params := map[string]string{ // common parameters, not include Signature
		"Format":           "JSON",
		"Version":          "2015-05-01",
		"SignatureMethod":  "HMAC-SHA1",
		"AccessKeyId":      c.accessKeyID,
		"SignatureVersion": "1.0",
		"SignatureNonce":   u.String(),
		"Timestamp":        time.Now().UTC().Format(time.RFC3339),
	}
	return params
}

// CreateRole :
func (c *Client) CreateRole(roleName, assumeRolePolicyDocument, decription string) (*CreateRoleResponse, error) {
	action := "CreateRole"
	params := map[string]string{}
	params["RoleName"] = roleName
	params["AssumeRolePolicyDocument"] = assumeRolePolicyDocument
	params["Description"] = decription
	params["Action"] = action

	var body []byte
	var err error

	if c.RequestManager != nil {
		body, err = c.RequestManager.sendRequest(params)
	} else {
		body, err = c.sendRequest(params)
	}

	if err != nil {
		return nil, err
	}
	resp := &CreateRoleResponse{}
	err = json.Unmarshal(body, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// GetRole :
func (c *Client) GetRole(roleName string) (*GetRoleResponse, error) {
	action := "GetRole"
	params := map[string]string{}
	params["RoleName"] = roleName
	params["Action"] = action

	var body []byte
	var err error

	if c.RequestManager != nil {
		body, err = c.RequestManager.sendRequest(params)
	} else {
		body, err = c.sendRequest(params)
	}

	if err != nil {
		return nil, err
	}
	resp := &GetRoleResponse{}
	err = json.Unmarshal(body, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// DeleteRole :
func (c *Client) DeleteRole(roleName string) (*DeleteRoleResponse, error) {
	action := "DeleteRole"
	params := map[string]string{}
	params["RoleName"] = roleName
	params["Action"] = action

	var body []byte
	var err error

	if c.RequestManager != nil {
		body, err = c.RequestManager.sendRequest(params)
	} else {
		body, err = c.sendRequest(params)
	}

	if err != nil {
		return nil, err
	}
	resp := &DeleteRoleResponse{}
	err = json.Unmarshal(body, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// ListRoles :
func (c *Client) ListRoles() (*ListRolesResponse, error) {
	action := "ListRoles"
	params := map[string]string{}
	params["Action"] = action


	var body []byte
	var err error

	if c.RequestManager != nil {
		body, err = c.RequestManager.sendRequest(params)
	} else {
		body, err = c.sendRequest(params)
	}

	if err != nil {
		return nil, err
	}
	resp := &ListRolesResponse{}
	err = json.Unmarshal(body, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// CreatePolicy :
func (c *Client) CreatePolicy(policyName, policyDocument, description string) (*CreatePolicyResponse, error) {
	action := "CreatePolicy"
	params := map[string]string{}
	params["PolicyName"] = policyName
	params["Description"] = description
	params["PolicyDocument"] = policyDocument
	params["Action"] = action


	var body []byte
	var err error

	if c.RequestManager != nil {
		body, err = c.RequestManager.sendRequest(params)
	} else {
		body, err = c.sendRequest(params)
	}

	if err != nil {
		return nil, err
	}
	resp := &CreatePolicyResponse{}
	err = json.Unmarshal(body, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// GetPolicy :
func (c *Client) GetPolicy(policyName string, policyType string) (*GetPolicyResponse, error) {
	action := "GetPolicy"
	params := map[string]string{}
	params["PolicyName"] = policyName
	params["Action"] = action
	params["PolicyType"] = policyType

	var body []byte
	var err error

	if c.RequestManager != nil {
		body, err = c.RequestManager.sendRequest(params)
	} else {
		body, err = c.sendRequest(params)
	}


	if err != nil {
		return nil, err
	}
	resp := &GetPolicyResponse{}
	err = json.Unmarshal(body, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// ListPolicyVersions :
func (c *Client) ListPolicyVersions(policyName string, policyType string) (*ListPolicyVersionsResponse, error) {
	action := "ListPolicyVersions"
	params := map[string]string{}
	params["PolicyName"] = policyName
	params["Action"] = action
	params["PolicyType"] = policyType

	var body []byte
	var err error

	if c.RequestManager != nil {
		body, err = c.RequestManager.sendRequest(params)
	} else {
		body, err = c.sendRequest(params)
	}

	if err != nil {
		return nil, err
	}
	resp := &ListPolicyVersionsResponse{}
	err = json.Unmarshal(body, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// GetPolicyVersion :
func (c *Client) GetPolicyVersion(policyName, policyType, VersionID string) (*PolicyVersionResponse, error) {
	action := "GetPolicyVersion"
	params := map[string]string{}
	params["PolicyName"] = policyName
	params["Action"] = action
	params["PolicyType"] = policyType
	params["VersionId"] = VersionID

	var body []byte
	var err error

	if c.RequestManager != nil {
		body, err = c.RequestManager.sendRequest(params)
	} else {
		body, err = c.sendRequest(params)
	}


	if err != nil {
		return nil, err
	}
	resp := &PolicyVersionResponse{}
	err = json.Unmarshal(body, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// CreatePolicyVersion :
func (c *Client) CreatePolicyVersion(policyName, policyDocument, setAsDefault string) (*PolicyVersionResponse, error) {
	action := "CreatePolicyVersion"
	params := map[string]string{}
	params["PolicyName"] = policyName
	params["Action"] = action
	params["PolicyDocument"] = policyDocument
	params["SetAsDefault"] = setAsDefault

	var body []byte
	var err error

	if c.RequestManager != nil {
		body, err = c.RequestManager.sendRequest(params)
	} else {
		body, err = c.sendRequest(params)
	}

	if err != nil {
		return nil, err
	}
	resp := &PolicyVersionResponse{}
	err = json.Unmarshal(body, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// DeletePolicy :
func (c *Client) DeletePolicy(policyName string) (*DeletePolicyResponse, error) {
	action := "DeletePolicy"
	params := map[string]string{}
	params["Action"] = action
	params["PolicyName"] = policyName


	var body []byte
	var err error

	if c.RequestManager != nil {
		body, err = c.RequestManager.sendRequest(params)
	} else {
		body, err = c.sendRequest(params)
	}

	if err != nil {
		return nil, err
	}
	resp := &DeletePolicyResponse{}
	err = json.Unmarshal(body, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// ListPolicies :
func (c *Client) ListPolicies() (*ListPoliciesResponse, error) {
	action := "ListPolicies"
	params := map[string]string{}
	params["Action"] = action


	var body []byte
	var err error

	if c.RequestManager != nil {
		body, err = c.RequestManager.sendRequest(params)
	} else {
		body, err = c.sendRequest(params)
	}

	if err != nil {
		return nil, err
	}
	resp := &ListPoliciesResponse{}
	err = json.Unmarshal(body, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// ListPoliciesForRole :
func (c *Client) ListPoliciesForRole(roleName string) (*ListPoliciesForRoleResponse, error) {
	action := "ListPoliciesForRole"
	params := map[string]string{}
	params["Action"] = action
	params["RoleName"] = roleName


	var body []byte
	var err error

	if c.RequestManager != nil {
		body, err = c.RequestManager.sendRequest(params)
	} else {
		body, err = c.sendRequest(params)
	}

	if err != nil {
		return nil, err
	}
	resp := &ListPoliciesForRoleResponse{}
	err = json.Unmarshal(body, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// AttachPolicyToRole :
func (c *Client) AttachPolicyToRole(policyType, policyName, roleName string) (*AttachPolicyToRoleResponse, error) {
	action := "AttachPolicyToRole"
	params := map[string]string{}
	params["Action"] = action
	params["PolicyType"] = policyType
	params["PolicyName"] = policyName
	params["RoleName"] = roleName


	var body []byte
	var err error

	if c.RequestManager != nil {
		body, err = c.RequestManager.sendRequest(params)
	} else {
		body, err = c.sendRequest(params)
	}

	if err != nil {
		return nil, err
	}
	resp := &AttachPolicyToRoleResponse{}
	err = json.Unmarshal(body, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// DetachPolicyFromRole :
func (c *Client) DetachPolicyFromRole(policyType, policyName, roleName string) (*DetachPolicyFromRoleResponse, error) {
	action := "DetachPolicyFromRole"
	params := map[string]string{}
	params["Action"] = action

	params["PolicyType"] = policyType
	params["PolicyName"] = policyName
	params["RoleName"] = roleName


	var body []byte
	var err error

	if c.RequestManager != nil {
		body, err = c.RequestManager.sendRequest(params)
	} else {
		body, err = c.sendRequest(params)
	}

	if err != nil {
		return nil, err
	}
	resp := &DetachPolicyFromRoleResponse{}
	err = json.Unmarshal(body, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *Client) sendRequest(rawParams map[string]string) ([]byte, error) {
	var body []byte
	err := try.Do(func(attempt int) (bool, error) {
		params := make(map[string]string)
		for k, v := range rawParams {
			params[k] = v
		}
		commonParams := c.getCommontParam()
		for k, v := range commonParams {
			params[k] = v
		}
		signature := CreateSignature(http.MethodPost, CreateQueryStr(params), c.accessKeySecret)
		params["Signature"] = signature

		targetURL := fmt.Sprintf("%s://%s/", c.endpoint.Scheme, c.endpoint.Host)
		queryParams := strings.NewReader(CreateQueryStr(params))

		req, err := http.NewRequest(http.MethodPost, targetURL, queryParams)
		if err != nil {
			return false, err
		}
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		resp, err := c.hclient.Do(req)
		if err != nil {
			//  if err is url.Error, just retry
			if _, ok := err.(*url.Error); ok {
				time.Sleep(time.Duration(getWaitIntervalInMS(int32(attempt))) * time.Millisecond)
				return int32(attempt) <= c.option.retryTimes, err
			}
			// else stop retry
			return false, err
		}
		defer resp.Body.Close()

		body, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return false, err
		}

		if resp.StatusCode != http.StatusOK {
			// if StatusCode == 500 or 503, just retry
			ramServiceErr := ServiceError{}
			ramServiceErr.HTTPStatus = resp.StatusCode
			err = json.Unmarshal(body, &ramServiceErr)
			if err != nil {
				return false, err
			}
			if resp.StatusCode == http.StatusInternalServerError ||
				resp.StatusCode == http.StatusServiceUnavailable {
				time.Sleep(time.Duration(getWaitIntervalInMS(int32(attempt))) * time.Millisecond)
				return int32(attempt) <= c.option.retryTimes, ramServiceErr
			}

			return false, ramServiceErr
		}
		return false, nil
	})
	if err != nil {
		return nil, err
	}
	return body, nil
}
