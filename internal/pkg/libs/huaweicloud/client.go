package huaweicloud

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth/global"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/config"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/httphandler"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iam/v3/model"

	obs "github.com/huaweicloud/huaweicloud-sdk-go-obs/obs"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth/basic"
	bss "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/bss/v2"
	cbr "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/cbr/v1"
	ces "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/ces/v1"
	ecs "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/ecs/v2"
	eip "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/eip/v2"
	evs "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/evs/v2"
	iam "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iam/v3"
	ims "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/ims/v2"
	vpc "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/vpc/v2"
)

type Client struct {
	AccessConfig
}

// AccessConfig is for common configuration related to HuaweiCloud access
type AccessConfig struct {
	// The access key of the HuaweiCloud to use.
	// If omitted, the HW_ACCESS_KEY environment variable is used.
	AccessKey string `mapstructure:"access_key" required:"true"`
	// The secret key of the HuaweiCloud to use.
	// If omitted, the HW_SECRET_KEY environment variable is used.
	SecretKey string `mapstructure:"secret_key" required:"true"`
	// The HuaweiCloud region in which to launch the server to create the image.
	// If omitted, the HW_REGION_NAME environment variable is used.
	Region string `mapstructure:"region" required:"true"`
	// The name of the project to login with.
	// If omitted, the HW_PROJECT_NAME environment variable or `region` is used.
	ProjectName string `mapstructure:"project_name" required:"false"`
	// The ID of the project to login with.
	// If omitted, the HW_PROJECT_ID environment variable is used.
	ProjectID string `mapstructure:"project_id" required:"false"`
	// The security token to authenticate with a
	// [temporary security credential](https://support.huaweicloud.com/intl/en-us/iam_faq/iam_01_0620.html).
	// If omitted, the HW_SECURITY_TOKEN environment variable is used.
	SecurityToken string `mapstructure:"security_token" required:"false"`
	// The Identity authentication URL.
	// If omitted, the HW_AUTH_URL environment variable is used.
	// This is not required if you use HuaweiCloud.
	IdentityEndpoint string `mapstructure:"auth_url" required:"false"`
	// Trust self-signed SSL certificates.
	// By default this is false.
	Insecure bool `mapstructure:"insecure" required:"false"`

	cloud string
}

func NewAccessConfig(ak, sk, region string) (*AccessConfig, error) {
	cfg := &AccessConfig{
		AccessKey: ak,
		SecretKey: sk,
		Region:    region,
	}

	if err := cfg.Prepare(); err != nil {
		return nil, err
	}
	return cfg, nil
}

func (c *AccessConfig) Prepare() error {
	if c.AccessKey == "" {
		c.AccessKey = os.Getenv("HW_ACCESS_KEY")
	}
	if c.SecretKey == "" {
		c.SecretKey = os.Getenv("HW_SECRET_KEY")
	}
	if c.Region == "" {
		c.Region = os.Getenv("HW_REGION_NAME")
	}
	// access parameters validation
	if c.AccessKey == "" || c.SecretKey == "" {
		paraErr := fmt.Errorf("access_key, secret_key  must be set")
		return paraErr
	}

	//// access parameters validation
	//if c.AccessKey == "" || c.SecretKey == "" || c.Region == "" {
	//	paraErr := fmt.Errorf("access_key, secret_key and region must be set")
	//	return paraErr
	//}

	if c.SecurityToken == "" {
		c.SecurityToken = os.Getenv("HW_SECURITY_TOKEN")
	}
	if c.ProjectID == "" {
		c.ProjectID = os.Getenv("HW_PROJECT_ID")
	}
	if c.ProjectName == "" {
		c.ProjectName = os.Getenv("HW_PROJECT_NAME")
	}
	// if neither "project_name" nor HW_PROJECT_NAME was specified, defaults to c.Region
	if c.ProjectName == "" {
		c.ProjectName = c.Region
	}

	if c.IdentityEndpoint == "" {
		c.IdentityEndpoint = os.Getenv("HW_AUTH_URL")
	}
	// if neither "auth_url" nor HW_AUTH_URL was specified, defaults to "iam.xxx.myhuaweicloud.com"
	// In Europe site(e.g. eu-west-101), the default endpoint is "iam.eu-west-10x.myhuaweicloud.eu"
	if c.IdentityEndpoint == "" {
		c.IdentityEndpoint = buildDefaultIamEndpoint(c.Region)
	}

	// 获取在此区域的项目列表
	if c.ProjectID == "" {
		projectID, err := c.getProjectID(c.ProjectName)
		if err != nil {
			log.Println(err.Error())
			return err
		}

		c.ProjectID = projectID
	}

	cloudDomain, err := GetCloudFromAuth(c.IdentityEndpoint)
	if err != nil {
		return err
	}
	c.cloud = cloudDomain

	return nil
}

type ServiceCatalog struct {
	Name  string
	Scope string
	Admin bool
}

// NewHcClient is the common client using huaweicloud-sdk-go-v3 package
func NewHcClient(c *AccessConfig, region, product string) (*core.HcHttpClient, error) {
	endpoint := GetServiceEndpoint(c.cloud, product, region)
	if endpoint == "" {
		return nil, fmt.Errorf("failed to get the endpoint of %q service in region %s", product, region)
	}

	builder := core.NewHcHttpClientBuilder().WithEndpoints([]string{endpoint}).WithHttpConfig(buildHTTPConfig(c))

	credentials := basic.Credentials{
		BaseCredentials: auth.BaseCredentials{
			AK:            c.AccessKey,
			SK:            c.SecretKey,
			SecurityToken: c.SecurityToken,
		},

		// 作为请求头传递
		ProjectId: c.ProjectID,
	}
	builder.WithCredential(&credentials)

	headers := map[string]string{
		//	"User-Agent": UserAgent,
	}
	return builder.Build().PreInvoke(headers), nil
}

// https://developer.huaweicloud.com/endpoint
var serviceEndpoints = map[string]ServiceCatalog{
	"ecs": {
		Name: "ecs",
	},
	"ims": {
		Name: "ims",
	},
	"vpc": {
		Name: "vpc",
	},
	"eip": {
		Name: "vpc",
	},
	"evs": {
		Name: "evs",
	},
	"cbr": {
		Name: "cbr",
	},
	"ces": {
		Name: "ces",
	},

	"bss": {
		Name:  "bss",
		Scope: "global",
	},
	"iam": {
		Name:  "iam",
		Scope: "global",
	},
	"oss": {
		Name: "obs",
	},
}

// GetServiceEndpoint try to get the endpoint from customizing map
func GetServiceEndpoint(cloud, srv, region string) string {
	// get the endpoint from build-in service catalog
	catalog, ok := serviceEndpoints[srv]
	if !ok {
		return ""
	}

	var ep string
	if catalog.Scope == "global" {
		ep = fmt.Sprintf("https://%s.%s/", catalog.Name, cloud)
	} else {
		ep = fmt.Sprintf("https://%s.%s.%s/", catalog.Name, region, cloud)
	}
	return ep
}

func buildDefaultIamEndpoint(region string) string {
	if region == "" {
		return fmt.Sprintf("https://iam.myhuaweicloud.com")
	}
	if strings.HasPrefix(region, "eu-west-10") {
		// In Europe site(e.g. eu-west-101), the default endpoint is "iam.eu-west-10x.myhuaweicloud.eu"
		return fmt.Sprintf("https://iam.%s.myhuaweicloud.eu", region)
	}
	return fmt.Sprintf("https://iam.%s.myhuaweicloud.com", region)
}

//func buildDefaultIamEndpoint() string {
//	return fmt.Sprintf("https://iam.myhuaweicloud.com")
//}

func GetCloudFromAuth(auth string) (string, error) {
	var cloud string

	u, err := url.Parse(auth)
	if err != nil {
		return "", err
	}

	// the parsed Host in host:port format, get rid of the port
	hosts := strings.SplitN(u.Host, ":", 2)

	subhosts := strings.Split(hosts[0], ".")
	total := len(subhosts)
	if total == 3 {
		// without region: iam.myhuaweicloud.com
		cloud = strings.Join(subhosts[1:], ".")
	} else if total > 3 {
		// with region: iam.cn-north-1.myhuaweicloud.com
		// iam.eu-west-0.prod-cloud-ocb.orange-business.com
		cloud = strings.Join(subhosts[2:], ".")
	} else {
		return "", fmt.Errorf("the auth_url is invalid")
	}

	return cloud, nil
}

func (c *AccessConfig) getProjectID(projectName string) (string, error) {
	builder := core.NewHcHttpClientBuilder().WithEndpoints([]string{c.IdentityEndpoint}).WithHttpConfig(buildHTTPConfig(c))

	credentials := global.Credentials{
		BaseCredentials: auth.BaseCredentials{
			AK:            c.AccessKey,
			SK:            c.SecretKey,
			SecurityToken: c.SecurityToken,
		},
	}
	builder.WithCredentialsType("global.Credentials").WithCredential(&credentials)

	headers := map[string]string{}
	client := iam.NewIamClient(builder.Build().PreInvoke(headers))
	request := &model.KeystoneListProjectsRequest{
		Name: &projectName,
	}

	response, err := client.KeystoneListProjects(request)
	if err != nil {
		return "", fmt.Errorf("can not get the project ID of %s: %s", c.ProjectName, err)
	}

	if response.Projects == nil || len(*response.Projects) == 0 {
		return "", fmt.Errorf("can not get the project ID of %s", c.ProjectName)
	}

	queriedProjects := *response.Projects
	return queriedProjects[0].Id, nil
}

func buildHTTPConfig(c *AccessConfig) *config.HttpConfig {
	httpConfig := config.DefaultHttpConfig()

	if c.Insecure {
		httpConfig = httpConfig.WithIgnoreSSLVerification(true)
	}

	if logEnabled() {
		httpHandler := httphandler.NewHttpHandler().
			AddRequestHandler(logRequestHandler).
			AddResponseHandler(logResponseHandler)
		httpConfig = httpConfig.WithHttpHandler(httpHandler)
	}

	if proxyURL := getProxyFromEnv(); proxyURL != "" {
		if parsed, err := url.Parse(proxyURL); err == nil {
			log.Printf("[DEBUG] using https proxy: %s://%s", parsed.Scheme, parsed.Host)

			httpProxy := config.Proxy{
				Schema:   parsed.Scheme,
				Host:     parsed.Host,
				Username: parsed.User.Username(),
			}
			if pwd, ok := parsed.User.Password(); ok {
				httpProxy.Password = pwd
			}

			httpConfig = httpConfig.WithProxy(&httpProxy)
		} else {
			log.Printf("[WARN] parsing https proxy failed: %s", err)
		}
	}

	return httpConfig
}

func logEnabled() bool {
	debugEnv := os.Getenv("HW_DEBUG")
	return debugEnv != "" && debugEnv != "0"
}

func getProxyFromEnv() string {
	var proxyURL string

	envNames := []string{"HTTPS_PROXY", "https_proxy"}
	for _, n := range envNames {
		if val := os.Getenv(n); val != "" {
			proxyURL = val
			break
		}
	}

	return proxyURL
}

// HcImsClient is the IMS service client using huaweicloud-sdk-go-v3 package
func (c *AccessConfig) HcImsClient(region string) (*ims.ImsClient, error) {
	hcClient, err := NewHcClient(c, region, "ims")
	if err != nil {
		return nil, err
	}

	return ims.NewImsClient(hcClient), nil
}

// HcEcsClient is the ECS service client using huaweicloud-sdk-go-v3 package
func (c *AccessConfig) HcEcsClient(region string) (*ecs.EcsClient, error) {
	hcClient, err := NewHcClient(c, region, "ecs")
	if err != nil {
		return nil, err
	}

	return ecs.NewEcsClient(hcClient), nil
}

// HcVpcClient is the VPC service client using huaweicloud-sdk-go-v3 package
func (c *AccessConfig) HcVpcClient(region string) (*vpc.VpcClient, error) {
	hcClient, err := NewHcClient(c, region, "vpc")
	if err != nil {
		return nil, err
	}

	return vpc.NewVpcClient(hcClient), nil
}

// HcEipClient is the EIP service client using huaweicloud-sdk-go-v3 package
func (c *AccessConfig) HcEipClient(region string) (*eip.EipClient, error) {
	hcClient, err := NewHcClient(c, region, "eip")
	if err != nil {
		return nil, err
	}

	return eip.NewEipClient(hcClient), nil
}

// HcEvsClient is the EVS service client using huaweicloud-sdk-go-v3 package
func (c *AccessConfig) HcEvsClient(region string) (*evs.EvsClient, error) {
	hcClient, err := NewHcClient(c, region, "evs")
	if err != nil {
		return nil, err
	}

	return evs.NewEvsClient(hcClient), nil
}

// HcBssClient is the Business service client using huaweicloud-sdk-go-v3 package
func (c *AccessConfig) HcBssClient() (*bss.BssClient, error) {
	hcClient, err := NewHcClient(c, "", "bss")
	if err != nil {
		return nil, err
	}

	return bss.NewBssClient(hcClient), nil
}

// HcIamClient is the IMS service client using huaweicloud-sdk-go-v3 package
func (c *AccessConfig) HcIamClient() (*iam.IamClient, error) {
	hcClient, err := NewHcClient(c, "", "iam")
	if err != nil {
		return nil, err
	}

	return iam.NewIamClient(hcClient), nil
}

// HcOssClient is the Obs service client
func (c *AccessConfig) HcOssClient(region string) (*obs.ObsClient, error) {
	endpoint := GetServiceEndpoint(c.cloud, "oss", region)
	if endpoint == "" {
		return nil, fmt.Errorf("failed to get the endpoint of %q service in region %s", "oss", region)
	}
	fmt.Println(endpoint)
	ossClient, err := obs.New(c.AccessKey, c.SecretKey, endpoint)
	if err != nil {
		return nil, err
	}

	return ossClient, nil
}

// HcCbrClient is the Obs service client
func (c *AccessConfig) HcCbrClient(region string) (*cbr.CbrClient, error) {
	hcClient, err := NewHcClient(c, "", "cbr")
	if err != nil {
		return nil, err
	}

	return cbr.NewCbrClient(hcClient), nil
}

// HcCesClient is the Obs service client
func (c *AccessConfig) HcCesClient(region string) (*ces.CesClient, error) {
	hcClient, err := NewHcClient(c, "", "ces")
	if err != nil {
		return nil, err
	}

	return ces.NewCesClient(hcClient), nil
}
