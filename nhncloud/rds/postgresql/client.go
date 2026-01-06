package postgresql

import (
	"context"
	"fmt"
	"net/http"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/credentials"
	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/internal/client"
	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/internal/endpoint"
)

type Client struct {
	region        string
	appKey        string
	credentials   credentials.Credentials
	httpClient    *client.Client
	tokenProvider *client.OAuthTokenProvider
	debug         bool
}

func NewClient(region, appKey string, creds credentials.Credentials, hc *http.Client, debug bool) *Client {
	c := &Client{
		region:      region,
		appKey:      appKey,
		credentials: creds,
		debug:       debug,
	}

	if creds != nil {
		c.tokenProvider = client.NewOAuthTokenProvider(
			creds.GetAccessKeyID(),
			creds.GetSecretAccessKey(),
		)
		c.initHTTPClient()
	}

	return c
}

func (c *Client) initHTTPClient() {
	baseURL := endpoint.ResolveWithAppKey(endpoint.ServiceRDSPG, c.region, c.appKey)
	opts := []client.ClientOption{
		client.WithDebug(c.debug),
	}
	c.httpClient = client.NewClient(baseURL, c.tokenProvider, opts...)
}

func (c *Client) ListInstances(ctx context.Context) (*ListInstancesOutput, error) {
	var out ListInstancesOutput
	if err := c.httpClient.GET(ctx, "/db-instances", &out); err != nil {
		return nil, fmt.Errorf("list instances: %w", err)
	}
	return &out, nil
}

func (c *Client) GetInstance(ctx context.Context, instanceID string) (*GetInstanceOutput, error) {
	var out GetInstanceOutput
	if err := c.httpClient.GET(ctx, "/db-instances/"+instanceID, &out); err != nil {
		return nil, fmt.Errorf("get instance %s: %w", instanceID, err)
	}
	return &out, nil
}

func (c *Client) CreateInstance(ctx context.Context, input *CreateInstanceInput) (*CreateInstanceOutput, error) {
	var out CreateInstanceOutput
	if err := c.httpClient.POST(ctx, "/db-instances", input, &out); err != nil {
		return nil, fmt.Errorf("create instance: %w", err)
	}
	return &out, nil
}

func (c *Client) ModifyInstance(ctx context.Context, instanceID string, input *ModifyInstanceInput) (*GetInstanceOutput, error) {
	var out GetInstanceOutput
	if err := c.httpClient.PUT(ctx, "/db-instances/"+instanceID, input, &out); err != nil {
		return nil, fmt.Errorf("modify instance %s: %w", instanceID, err)
	}
	return &out, nil
}

func (c *Client) DeleteInstance(ctx context.Context, instanceID string) (*JobOutput, error) {
	var out JobOutput
	if err := c.httpClient.DELETE(ctx, "/db-instances/"+instanceID, &out); err != nil {
		return nil, fmt.Errorf("delete instance %s: %w", instanceID, err)
	}
	return &out, nil
}

func (c *Client) StartInstance(ctx context.Context, instanceID string) (*JobOutput, error) {
	var out JobOutput
	if err := c.httpClient.POST(ctx, "/db-instances/"+instanceID+"/start", nil, &out); err != nil {
		return nil, fmt.Errorf("start instance %s: %w", instanceID, err)
	}
	return &out, nil
}

func (c *Client) StopInstance(ctx context.Context, instanceID string) (*JobOutput, error) {
	var out JobOutput
	if err := c.httpClient.POST(ctx, "/db-instances/"+instanceID+"/stop", nil, &out); err != nil {
		return nil, fmt.Errorf("stop instance %s: %w", instanceID, err)
	}
	return &out, nil
}

func (c *Client) RestartInstance(ctx context.Context, instanceID string, useOnlineFailover bool) (*JobOutput, error) {
	req := map[string]interface{}{
		"useOnlineFailover": useOnlineFailover,
	}
	var out JobOutput
	if err := c.httpClient.POST(ctx, "/db-instances/"+instanceID+"/restart", req, &out); err != nil {
		return nil, fmt.Errorf("restart instance %s: %w", instanceID, err)
	}
	return &out, nil
}

func (c *Client) ListInstanceGroups(ctx context.Context) (*ListInstanceGroupsOutput, error) {
	var out ListInstanceGroupsOutput
	if err := c.httpClient.GET(ctx, "/db-instance-groups", &out); err != nil {
		return nil, fmt.Errorf("list instance groups: %w", err)
	}
	return &out, nil
}

func (c *Client) ListFlavors(ctx context.Context) (*ListFlavorsOutput, error) {
	var out ListFlavorsOutput
	if err := c.httpClient.GET(ctx, "/db-flavors", &out); err != nil {
		return nil, fmt.Errorf("list flavors: %w", err)
	}
	return &out, nil
}

func (c *Client) ListVersions(ctx context.Context) (*ListVersionsOutput, error) {
	var out ListVersionsOutput
	if err := c.httpClient.GET(ctx, "/db-versions", &out); err != nil {
		return nil, fmt.Errorf("list versions: %w", err)
	}
	return &out, nil
}

func (c *Client) ListSecurityGroups(ctx context.Context) (*ListSecurityGroupsOutput, error) {
	var out ListSecurityGroupsOutput
	if err := c.httpClient.GET(ctx, "/db-security-groups", &out); err != nil {
		return nil, fmt.Errorf("list security groups: %w", err)
	}
	return &out, nil
}

func (c *Client) ListParameterGroups(ctx context.Context, dbVersion string) (*ListParameterGroupsOutput, error) {
	path := "/parameter-groups"
	if dbVersion != "" {
		path += "?dbVersion=" + dbVersion
	}
	var out ListParameterGroupsOutput
	if err := c.httpClient.GET(ctx, path, &out); err != nil {
		return nil, fmt.Errorf("list parameter groups: %w", err)
	}
	return &out, nil
}

func (c *Client) ListBackups(ctx context.Context, instanceID string, page, size int) (*ListBackupsOutput, error) {
	path := fmt.Sprintf("/backups?page=%d&size=%d", page, size)
	if instanceID != "" {
		path += "&dbInstanceId=" + instanceID
	}
	var out ListBackupsOutput
	if err := c.httpClient.GET(ctx, path, &out); err != nil {
		return nil, fmt.Errorf("list backups: %w", err)
	}
	return &out, nil
}

func (c *Client) CreateBackup(ctx context.Context, instanceID, name string) (*JobOutput, error) {
	req := map[string]string{"backupName": name}
	var out JobOutput
	if err := c.httpClient.POST(ctx, "/db-instances/"+instanceID+"/backup", req, &out); err != nil {
		return nil, fmt.Errorf("create backup: %w", err)
	}
	return &out, nil
}

func (c *Client) DeleteBackup(ctx context.Context, backupID string) (*JobOutput, error) {
	var out JobOutput
	if err := c.httpClient.DELETE(ctx, "/backups/"+backupID, &out); err != nil {
		return nil, fmt.Errorf("delete backup %s: %w", backupID, err)
	}
	return &out, nil
}

func (c *Client) ListSubnets(ctx context.Context) (*ListSubnetsOutput, error) {
	var out ListSubnetsOutput
	if err := c.httpClient.GET(ctx, "/network/subnets", &out); err != nil {
		return nil, fmt.Errorf("list subnets: %w", err)
	}
	return &out, nil
}

func (c *Client) GetNetworkInfo(ctx context.Context, instanceID string) (*NetworkInfo, error) {
	var out NetworkInfo
	if err := c.httpClient.GET(ctx, "/db-instances/"+instanceID+"/network-info", &out); err != nil {
		return nil, fmt.Errorf("get network info: %w", err)
	}
	return &out, nil
}

func (c *Client) EnablePublicAccess(ctx context.Context, instanceID string) (*JobOutput, error) {
	req := map[string]bool{"usePublicAccess": true}
	var out JobOutput
	if err := c.httpClient.PUT(ctx, "/db-instances/"+instanceID+"/network-info", req, &out); err != nil {
		return nil, fmt.Errorf("enable public access: %w", err)
	}
	return &out, nil
}

func (c *Client) DisablePublicAccess(ctx context.Context, instanceID string) (*JobOutput, error) {
	req := map[string]bool{"usePublicAccess": false}
	var out JobOutput
	if err := c.httpClient.PUT(ctx, "/db-instances/"+instanceID+"/network-info", req, &out); err != nil {
		return nil, fmt.Errorf("disable public access: %w", err)
	}
	return &out, nil
}

func (c *Client) EnableHA(ctx context.Context, instanceID string) (*JobOutput, error) {
	req := map[string]bool{"useHighAvailability": true}
	var out JobOutput
	if err := c.httpClient.PUT(ctx, "/db-instances/"+instanceID+"/high-availability", req, &out); err != nil {
		return nil, fmt.Errorf("enable HA: %w", err)
	}
	return &out, nil
}

func (c *Client) DisableHA(ctx context.Context, instanceID string) (*JobOutput, error) {
	req := map[string]bool{"useHighAvailability": false}
	var out JobOutput
	if err := c.httpClient.PUT(ctx, "/db-instances/"+instanceID+"/high-availability", req, &out); err != nil {
		return nil, fmt.Errorf("disable HA: %w", err)
	}
	return &out, nil
}

func (c *Client) ExpandStorage(ctx context.Context, instanceID string, newSize int, useOnlineFailover bool) (*JobOutput, error) {
	req := map[string]interface{}{
		"storageSize":       newSize,
		"useOnlineFailover": useOnlineFailover,
	}
	var out JobOutput
	if err := c.httpClient.PUT(ctx, "/db-instances/"+instanceID+"/storage-info", req, &out); err != nil {
		return nil, fmt.Errorf("expand storage: %w", err)
	}
	return &out, nil
}

func (c *Client) SetDeletionProtection(ctx context.Context, instanceID string, enabled bool) (*JobOutput, error) {
	req := map[string]bool{"useDeletionProtection": enabled}
	var out JobOutput
	if err := c.httpClient.PUT(ctx, "/db-instances/"+instanceID+"/deletion-protection", req, &out); err != nil {
		return nil, fmt.Errorf("set deletion protection: %w", err)
	}
	return &out, nil
}

func (c *Client) ListExtensions(ctx context.Context) (*ListExtensionsOutput, error) {
	var out ListExtensionsOutput
	if err := c.httpClient.GET(ctx, "/extensions", &out); err != nil {
		return nil, fmt.Errorf("list extensions: %w", err)
	}
	return &out, nil
}

func (c *Client) ListHBARules(ctx context.Context, instanceID string) (*ListHBARulesOutput, error) {
	var out ListHBARulesOutput
	if err := c.httpClient.GET(ctx, "/db-instances/"+instanceID+"/hba-rules", &out); err != nil {
		return nil, fmt.Errorf("list HBA rules: %w", err)
	}
	return &out, nil
}
