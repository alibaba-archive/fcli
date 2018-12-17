package ram

// default params for client options
const (
	DefaultRetryTimes          = 5
	DefaultMaxWaitIntervalInMS = int32(120000)
)

// Role :
type Role struct {
	RoleID                   string `json:"RoleId"`
	RoleName                 string `json:"RoleName"`
	Arn                      string `json:"Arn"`
	Description              string `json:"Description"`
	AssumeRolePolicyDocument string `json:"AssumeRolePolicyDocument"`
	CreateDate               string `json:"CreateDate"`
}

// Roles :
type Roles struct {
	Role []Role `json:"Role"`
}

// Policy :
type Policy struct {
	PolicyName      string `json:"PolicyName"`
	PolicyType      string `json:"PolicyType"`
	Description     string `json:"Description"`
	DefaultVersion  string `json:"DefaultVersion"`
	CreateDate      string `json:"CreateDate"`
	UpdateDate      string `json:"UpdateDate"`
	AttachmentCount int64  `json:"AttachmentCount"`
}

// PolicyStatement :
type PolicyStatement struct {
	Effect   string   `json:"Effect"`
	Action   []string `json:"Action"`
	Resource []string `json:"Resource"`
}

// PolicyDocument :
type PolicyDocument struct {
	Version   string            `json:"Version"`
	Statement []PolicyStatement `json:"Statement"`
}

// Policies :
type Policies struct {
	Policy []Policy `json:"Policy"`
}

// PolicyVersion :
type PolicyVersion struct {
	VersionID        string `json:"VersionId"`
	IsDefaultVersion bool   `json:"IsDefaultVersion"`
	CreateDate       string `json:"CreateDate"`
	PolicyDocument   string `json:"PolicyDocument"`
}

// PolicyVersionResponse :
type PolicyVersionResponse struct {
	RequestID     string        `json:"RequestId"`
	PolicyVersion PolicyVersion `json:"PolicyVersion"`
}

// PolicyVersions :
type PolicyVersions struct {
	PolicyVersion []PolicyVersion `json:"PolicyVersion"`
}

// ListPolicyVersionsResponse :
type ListPolicyVersionsResponse struct {
	RequestID      string         `json:"RequestId"`
	PolicyVersions PolicyVersions `json:"PolicyVersions"`
}

// CreateRoleResponse :
type CreateRoleResponse struct {
	RequestID string `json:"RequestId"`
	Role      Role   `json:"Role"`
}

// GetRoleResponse :
type GetRoleResponse struct {
	RequestID string `json:"RequestId"`
	Role      Role   `json:"Role"`
}

// DeleteRoleResponse :
type DeleteRoleResponse struct {
	RequestID string `json:"RequestId"`
}

// ListRolesResponse :
type ListRolesResponse struct {
	RequestID string `json:"RequestId"`
	Roles     Roles  `json:"Roles"`
}

// CreatePolicyResponse :
type CreatePolicyResponse struct {
	RequestID string `json:"RequestId"`
	Policy    Policy `json:"Policy"`
}

// GetPolicyResponse :
type GetPolicyResponse struct {
	RequestID string `json:"RequestId"`
	Policy    Policy `json:"Policy"`
}

// DeletePolicyResponse :
type DeletePolicyResponse struct {
	RequestID string `json:"RequestId"`
}

// ListPoliciesResponse :
type ListPoliciesResponse struct {
	RequestID   string   `json:"RequestId"`
	IsTruncated bool     `json:"IsTruncated"`
	Marker      string   `json:"Marker"`
	Policies    Policies `json:"Policies"`
}

// AttachPolicyToRoleResponse :
type AttachPolicyToRoleResponse struct {
	RequestID string `json:"RequestId"`
}

// DetachPolicyFromRoleResponse :
type DetachPolicyFromRoleResponse struct {
	RequestID string `json:"RequestId"`
}

// ListPoliciesForRoleResponse :
type ListPoliciesForRoleResponse struct {
	RequestID string   `json:"RequestId"`
	Policies  Policies `json:"Policies"`
}
