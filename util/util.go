package util

import (
	"encoding/json"
	"fmt"
	"github.com/aliyun/fcli/version"
	"io"
	"io/ioutil"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/aliyun-log-go-sdk"
	"github.com/denisbrodbeck/machineid"
	"github.com/spf13/viper"

	"math"

	"github.com/aliyun/fc-go-sdk"
	"github.com/aliyun/fcli/ram"
	"os"
	"os/exec"
)

const (
	// MaxLineNumPerGet defines GetLogs max lines per times
	MaxLineNumPerGet int64 = 100

	// EndpointFmt fc endpoint fmt
	EndpointFmt = `https://%s.%s.fc.aliyuncs.com`

	// LogEndpointFmt loghub endpoint fmt
	LogEndpointFmt = `%s.log.aliyuncs.com`

	// GetLogsIntervalTimeInSecs read logs for the function in streaming way
	// https://help.aliyun.com/document_detail/29029.html?spm=5176.doc29025.6.682.h3H3ac
	// 根据loghub api说明文档，实时数据写入至可查询的延时为1分钟
	// 因此查询间隔为1分钟，如果过短则会出现数据写入后查询时没有数据的情况，查询区间继续往下走则会漏掉这部分数据
	GetLogsIntervalTimeInSecs int64 = 1 * 60

	// IncompleteProgress defines progress status in GetLogs with specific query
	IncompleteProgress = "Incomplete"

	// TimeLayoutInLogs defines input time format for log cmds
	TimeLayoutInLogs = "2006-01-02T15:04:05Z"

	// KeyValueDelimiter defines key value delimiter of map params
	KeyValueDelimiter = "="
)

// GlobalConfig define the global configurations.
type GlobalConfig struct {
	Endpoint        string `yaml:"endpoint"`
	APIVersion      string `yaml:"api_version"`
	AccessKeyID     string `yaml:"access_key_id"`
	AccessKeySecret string `yaml:"access_key_secret"`
	SecurityToken   string `yaml:"security_token"`
	UserAgent       string
	Debug           bool   `yaml:"debug"`
	Timeout         uint   `yaml:"timeout"`
	SLSEndpoint     string `yaml:"sls_endpoint"`
}

// NewGlobalConfig create a global config.
func NewGlobalConfig() *GlobalConfig {
	cfg := &GlobalConfig{}
	cfg.APIVersion = "2016-08-15"
	cfg.Timeout = 60

	processVersion := runtime.Version() // eg: go1.11.2
	platform := runtime.GOOS            // eg: windows, darwin, freebsd, linux, and so on
	arch := runtime.GOARCH              // eg: 386, amd64, arm, s390x, and so on
	zone, _ := os.LookupEnv("LANG")
	mid, _ := machineid.ID()
	cfg.UserAgent = fmt.Sprintf("@alicloud/fcli/%s ( %s; OS %s %s; language %s; mid %s )",
		version.Version, processVersion, platform, arch, zone, mid)

	return cfg
}

// NewFClient create a fc client.
func NewFClient(cfg *GlobalConfig) (*fc.Client, error) {
	return fc.NewClient(
		cfg.Endpoint,
		cfg.APIVersion,
		cfg.AccessKeyID,
		cfg.AccessKeySecret,
		fc.WithSecurityToken(cfg.SecurityToken),
		fc.WithTimeout(cfg.Timeout),
		func(c *fc.Client) {
			c.Config.UserAgent = cfg.UserAgent
		},
	)
}

// NewRAMClient create ram client
func NewRAMClient(accessKeyID, accessKeySecret string) (*ram.Client, error) {
	const endpoint = "https://ram.aliyuncs.com"
	client, err := ram.NewClient(endpoint, accessKeyID, accessKeySecret)
	if err != nil {
		return nil, fmt.Errorf("get ram client err: %s", err)
	}
	return client, nil
}

// NewSLSClient create sls client
func NewSLSClient(cfg *GlobalConfig) *sls.Client {
	c := &sls.Client{
		Endpoint:        cfg.SLSEndpoint,
		AccessKeyID:     cfg.AccessKeyID,
		AccessKeySecret: cfg.AccessKeySecret,
		SecurityToken:   cfg.SecurityToken,
	}
	return c
}

// GetLogs read the log data from loghub with the count limit.
func GetLogs(store *sls.LogStore, topic, queryExp string, from, to, maxTotalLineNum int64, reverse bool) error {
	var offset int64
	lineNumPerGet := MaxLineNumPerGet
	if lineNumPerGet > maxTotalLineNum {
		lineNumPerGet = maxTotalLineNum
	}
	for offset < maxTotalLineNum {
		resp, err := store.GetLogs(topic, from, to, queryExp, lineNumPerGet, offset, reverse)
		if err != nil {
			return err
		}
		if resp.Progress == IncompleteProgress {
			// offset不变，继续读，不输出
			time.Sleep(1 * time.Second)
			continue
		} else {
			// Complete状态，输出
			PrettyPrintLogs(&resp.Logs, reverse)

			// 当返回的日志的count值小于指定的行数，则该段时间内的日志全部读取完毕
			if resp.Count < MaxLineNumPerGet {
				break
			}

			// 更改offset继续读取
			offset = offset + resp.Count
		}
	}
	return nil
}

// GetAllLogsWithinTimeRange ...
func GetAllLogsWithinTimeRange(store *sls.LogStore, topic, queryExp string, from, to int64) error {
	return GetLogs(store, topic, queryExp, from, to, math.MaxInt64, false)
}

// PrettyPrintLogs ..
func PrettyPrintLogs(logs *[]map[string]string, reverse bool) {
	if reverse {
		for i := len(*logs) - 1; i >= 0; i-- {
			PrettyPrintLog((*logs)[i])
		}
		return
	}
	for _, v := range *logs {
		PrettyPrintLog(v)
	}
}

// PrettyPrintLog ..
func PrettyPrintLog(v map[string]string) {
	timestamp, _ := strconv.ParseInt(v["__time__"], 10, 64)
	timeStr := time.Unix(timestamp, 0).Format(time.RFC3339)
	fmt.Println(fmt.Sprintf("%s\tserviceName:%s\tfunctionName:%s\tmessage:%s", timeStr,
		v["serviceName"], v["functionName"], v["message"]))
}

// GetRegions get region list of fc
func GetRegions() []string {
	return []string{
		"cn-beijing",
		"cn-qingdao",
		"cn-zhangjiakou",
		"cn-hangzhou",
		"cn-shanghai",
		"cn-shenzhen",
		"cn-hongkong",
		"ap-southeast-1",
		"ap-southeast-2",
		"ap-northeast-1",
		"us-west-1",
		"us-east-1",
		"eu-central-1",
		"ap-south-1",
	}
}

// GetRegionNoForEndpoint get region no from fc endpoint for endpoint
func GetRegionNoForEndpoint(endpoint string) string {
	regions := GetRegions()
	for _, region := range regions {
		if strings.Contains(endpoint, region) {
			return region
		}
	}
	return regions[0]
}

// GetRegionNoForSLSEndpoint get region no from fc endpoint for slsendpoint
func GetRegionNoForSLSEndpoint(endpoint string) string {
	// TODO Support the intranet endpoint.
	regions := []string{
		"cn-beijing",
		"cn-hangzhou",
		"cn-shanghai",
		"cn-qingdao",
		"cn-zhangjiakou",
		"cn-huhehaote",
		"cn-shenzhen",
		"cn-chengdu",
		"cn-hongkong",
		"ap-northeast-1",
		"ap-southeast-1",
		"ap-southeast-2",
		"ap-southeast-3",
		"ap-southeast-5",
		"me-east-1",
		"us-west-1",
		"eu-central-1",
		"us-east-1",
		"ap-south-1",
		"eu-west-1",
	}
	for _, region := range regions {
		if strings.Contains(endpoint, region) {
			return region
		}
	}
	return regions[0]
}

// GetUIDFromEndpoint extract the uid from the fc service endpoint.
func GetUIDFromEndpoint(endpoint string) (string, error) {
	tmp := strings.SplitN(endpoint, ".", 2)
	if len(tmp) < 1 {
		return "", fmt.Errorf("invalid fc endpoint: %s", endpoint)
	}
	uid := tmp[0]
	uid = strings.TrimPrefix(uid, "http://")
	uid = strings.TrimPrefix(uid, "https://")
	return uid, nil
}

// CreateRole create the RAM role.
func CreateRole(cli *ram.Client, roleName, principal string) (roleARN string, err error) {
	tmpl := `{"Statement": [{"Action": "sts:AssumeRole", "Effect": "Allow", "Principal": { "%s": ["%s"]}}], "Version": "1"}`
	const roleType = "Service"
	doc := fmt.Sprintf(tmpl, roleType, principal)
	desc := fmt.Sprintf("create the role %s", roleName)
	resp, err := cli.CreateRole(roleName, doc, desc)
	if err != nil {
		return "", fmt.Errorf("failed to create role %s due to:%s", roleName, err)
	}
	roleARN = resp.Role.Arn
	return roleARN, nil
}

// CreatePolicy create the RAM policy.
func CreatePolicy(cli *ram.Client, policyName, action, resource string) error {
	tmpl := `{"Version": "1", "Statement": [{"Effect": "Allow", "Action": %s, "Resource": %s}]}`
	doc := fmt.Sprintf(tmpl, action, resource)
	desc := fmt.Sprintf("create the policy %s", policyName)
	_, err := cli.CreatePolicy(policyName, doc, desc)
	if err != nil {
		return fmt.Errorf("failed to create policy %s due to %s", policyName, err)
	}
	return err
}

// CreatePolicyVersion update resource permission.
func CreatePolicyVersion(cli *ram.Client, policyName, action, resource string) error {
	tmpl := `{"Version": "1", "Statement": [{"Effect": "Allow", "Action": %s, "Resource": %s}]}`
	doc := fmt.Sprintf(tmpl, action, resource)
	_, err := cli.CreatePolicyVersion(policyName, doc, "true")
	if err != nil {
		return fmt.Errorf("failed to create policy %s due to %s", policyName, err)
	}
	return err
}

// GetDefaultPolicyVersion :
func GetDefaultPolicyVersion(cli *ram.Client, policyName, policyType string) (*ram.PolicyVersion, error) {
	versions, err := cli.ListPolicyVersions(policyName, policyType)
	if err != nil {
		return nil, err
	}
	var defaultVersion string
	for _, val := range versions.PolicyVersions.PolicyVersion {
		if val.IsDefaultVersion {
			defaultVersion = val.VersionID
			break
		}
	}
	versionResp, err := cli.GetPolicyVersion(policyName, policyType, defaultVersion)
	if err != nil {
		return nil, err
	}
	return &versionResp.PolicyVersion, nil
}

// CheckPolicyResourcePermission :
func CheckPolicyResourcePermission(p *ram.PolicyVersion, resource string) (bool, error) {
	hasResourcePermission := false
	pd := &ram.PolicyDocument{}
	err := json.Unmarshal([]byte(p.PolicyDocument), pd)
	if err != nil {
		return false, err
	}
	for _, val := range pd.Statement {
		for _, vr := range val.Resource {
			if resource == ("[\"" + vr + "\"]") {
				hasResourcePermission = true
				break
			}
		}
		if hasResourcePermission {
			break
		}
	}
	return hasResourcePermission, nil
}

// AttachPolicy attach the policy to the specified role.
func AttachPolicy(cli *ram.Client, policyName, roleName string) error {
	_, err := cli.AttachPolicyToRole("Custom", policyName, roleName)
	return err
}

// DownloadFromURL download file from the specified url.
func DownloadFromURL(url string, writer io.Writer) error {
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	_, err = io.Copy(writer, response.Body)
	return err
}

// GetBindingCmd return cmd and bind with out/err/in
func GetBindingCmd(name, arg, cmdString string) *exec.Cmd {
	cmd := exec.Command(name, arg, cmdString)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd
}

// CheckLocalImage Check local image exist
func CheckLocalImage(name, tag string) bool {
	_, err := exec.Command("docker", "inspect", name+":"+tag).Output()
	if err != nil {
		return false
	}
	return true
}

// ParseAdditionalVersionWeight parse route string to map and return
func ParseAdditionalVersionWeight(routes []string) map[string]float64 {
	additionalVersionWeight := make(map[string]float64)
	for _, route := range routes {
		items := strings.Split(strings.TrimSpace(route), KeyValueDelimiter)
		key := items[0]
		value, err := strconv.ParseFloat(route[len(items[0])+1:], 64)
		if err != nil {
			fmt.Printf("Error: can not parse version weight from route: %s\n", err)
			return (map[string]float64)(nil)
		}
		additionalVersionWeight[key] = value
	}
	return additionalVersionWeight
}

// GetTriggerConfig ...
func GetTriggerConfig(triggerType string, triggerConfigFile string) (interface{}, error) {
	viper.SetConfigFile(triggerConfigFile)
	err := viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s due to %v", triggerConfigFile, err)
	}

	switch triggerType {
	case fc.TRIGGER_TYPE_OSS:
		ossTriggerConfig := fc.OSSTriggerConfig{}
		viper.UnmarshalKey("triggerConfig", &ossTriggerConfig)
		return ossTriggerConfig, nil
	case fc.TRIGGER_TYPE_TIMER:
		timeTriggerConfig := fc.TimeTriggerConfig{}
		viper.UnmarshalKey("triggerConfig", &timeTriggerConfig)
		return timeTriggerConfig, nil
	case fc.TRIGGER_TYPE_LOG:
		logTriggerConfig := fc.LogTriggerConfig{}
		viper.UnmarshalKey("triggerConfig", &logTriggerConfig)
		return logTriggerConfig, nil
	case fc.TRIGGER_TYPE_CDN_EVENTS:
		cdnEvnetsTriggerConfig := fc.CDNEventsTriggerConfig{}
		viper.UnmarshalKey("triggerConfig", &cdnEvnetsTriggerConfig)
		return cdnEvnetsTriggerConfig, nil
	case fc.TRIGGER_TYPE_HTTP:
		httpTriggerConfig := fc.HTTPTriggerConfig{}
		viper.UnmarshalKey("triggerConfig", &httpTriggerConfig)
		return httpTriggerConfig, nil
	// case fc.TRIGGER_TYPE_RDS:
	// 	rdsTriggerConfig := fc.RdsTriggerConfig{}
	// 	viper.UnmarshalKey("triggerConfig", &rdsTriggerConfig)
	// 	return rdsTriggerConfig, nil
	// case fc.TRIGGER_TYPE_TABLESTORE:
	// 	tableStoreTriggerConfig := fc.TableStoreTriggerConfig{}
	// 	viper.UnmarshalKey("triggerConfig", &tableStoreTriggerConfig)
	// 	return tableStoreTriggerConfig, nil
	case fc.TRIGGER_TYPE_MNS_TOPIC:
		mnsTopicTriggerConfig := fc.MnsTopicTriggerConfig{}
		viper.UnmarshalKey("triggerConfig", &mnsTopicTriggerConfig)
		return mnsTopicTriggerConfig, nil
	default:
		return nil, fmt.Errorf("unsupported trigger type, expect oss, log, timer, http, cdn_events, mns_topic, actual %s", triggerType)
	}

}

// GetPublicImageDigest Get docker hub public image digest: sha256:xxxxxxx
func GetPublicImageDigest(name, tag string) (string, error) {
	// get token
	resp, err := http.Get("https://auth.docker.io/token?service=registry.docker.io&scope=repository:" + name + ":pull")
	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	err = resp.Body.Close()
	if err != nil {
		return "", err
	}

	var f interface{}
	err = json.Unmarshal(body, &f)
	if err != nil {
		return "", err
	}

	m := f.(map[string]interface{})
	token := m["token"].(string)

	// get digest
	req, _ := http.NewRequest("GET", "https://registry-1.docker.io/v2/"+name+"/manifests/"+tag, nil)
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Accept", "application/vnd.docker.distribution.manifest.v2+json")

	client := &http.Client{}
	resp, _ = client.Do(req)
	digest := resp.Header.Get("Docker-Content-Digest")

	return digest, nil
}
