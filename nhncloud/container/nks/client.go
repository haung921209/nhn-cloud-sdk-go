package nks

import (
	"context"
	"fmt"
	"net/http"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/credentials"
	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/internal/client"
)

type Client struct {
	region        string
	credentials   credentials.IdentityCredentials
	httpClient    *client.Client
	tokenProvider *client.IdentityTokenProvider
	debug         bool
}

func NewClient(region string, creds credentials.IdentityCredentials, hc *http.Client, debug bool) *Client {
	c := &Client{
		region:      region,
		credentials: creds,
		debug:       debug,
	}

	if creds != nil {
		c.tokenProvider = client.NewIdentityTokenProvider(
			creds.GetTenantID(),
			creds.GetUsername(),
			creds.GetPassword(),
		)
	}

	return c
}

func (c *Client) ensureClient(ctx context.Context) error {
	if c.httpClient != nil {
		return nil
	}

	if c.tokenProvider == nil {
		return ErrNoCredentials
	}

	baseURL, err := c.tokenProvider.GetServiceEndpoint("nks", c.region)
	if err != nil {
		return err
	}

	opts := []client.ClientOption{
		client.WithDebug(c.debug),
	}

	c.httpClient = client.NewClient(baseURL, c.tokenProvider, opts...)
	return nil
}

func (c *Client) ListClusters(ctx context.Context) (*ListClustersOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	var out ListClustersOutput
	if err := c.httpClient.GET(ctx, "/v1/clusters", &out); err != nil {
		return nil, fmt.Errorf("list clusters: %w", err)
	}
	return &out, nil
}

func (c *Client) GetCluster(ctx context.Context, clusterID string) (*GetClusterOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	var out GetClusterOutput
	if err := c.httpClient.GET(ctx, "/v1/clusters/"+clusterID, &out); err != nil {
		return nil, fmt.Errorf("get cluster %s: %w", clusterID, err)
	}
	return &out, nil
}

func (c *Client) CreateCluster(ctx context.Context, input *CreateClusterInput) (*CreateClusterOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	var out CreateClusterOutput
	if err := c.httpClient.POST(ctx, "/v1/clusters", input, &out); err != nil {
		return nil, fmt.Errorf("create cluster: %w", err)
	}
	return &out, nil
}

func (c *Client) DeleteCluster(ctx context.Context, clusterID string) error {
	if err := c.ensureClient(ctx); err != nil {
		return err
	}

	if err := c.httpClient.DELETE(ctx, "/v1/clusters/"+clusterID, nil); err != nil {
		return fmt.Errorf("delete cluster %s: %w", clusterID, err)
	}
	return nil
}

func (c *Client) UpdateCluster(ctx context.Context, clusterID string, input *UpdateClusterInput) error {
	if err := c.ensureClient(ctx); err != nil {
		return err
	}

	if err := c.httpClient.PATCH(ctx, "/v1/clusters/"+clusterID, input, nil); err != nil {
		return fmt.Errorf("update cluster %s: %w", clusterID, err)
	}
	return nil
}

func (c *Client) GetKubeconfig(ctx context.Context, clusterID string) (*GetKubeconfigOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	var kubeconfig string
	if err := c.httpClient.GET(ctx, "/v1/clusters/"+clusterID+"/kubeconfig", &kubeconfig); err != nil {
		return nil, fmt.Errorf("get kubeconfig for cluster %s: %w", clusterID, err)
	}
	return &GetKubeconfigOutput{Kubeconfig: kubeconfig}, nil
}

func (c *Client) ListNodeGroups(ctx context.Context, clusterID string) (*ListNodeGroupsOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	var out ListNodeGroupsOutput
	if err := c.httpClient.GET(ctx, "/v1/clusters/"+clusterID+"/nodegroups", &out); err != nil {
		return nil, fmt.Errorf("list node groups for cluster %s: %w", clusterID, err)
	}
	return &out, nil
}

func (c *Client) GetNodeGroup(ctx context.Context, clusterID, nodeGroupID string) (*GetNodeGroupOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	var out GetNodeGroupOutput
	if err := c.httpClient.GET(ctx, "/v1/clusters/"+clusterID+"/nodegroups/"+nodeGroupID, &out); err != nil {
		return nil, fmt.Errorf("get node group %s: %w", nodeGroupID, err)
	}
	return &out, nil
}

func (c *Client) CreateNodeGroup(ctx context.Context, clusterID string, input *CreateNodeGroupInput) (*CreateNodeGroupOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	var out CreateNodeGroupOutput
	if err := c.httpClient.POST(ctx, "/v1/clusters/"+clusterID+"/nodegroups", input, &out); err != nil {
		return nil, fmt.Errorf("create node group: %w", err)
	}
	return &out, nil
}

func (c *Client) UpdateNodeGroup(ctx context.Context, clusterID, nodeGroupID string, input *UpdateNodeGroupInput) error {
	if err := c.ensureClient(ctx); err != nil {
		return err
	}

	if err := c.httpClient.PATCH(ctx, "/v1/clusters/"+clusterID+"/nodegroups/"+nodeGroupID, input, nil); err != nil {
		return fmt.Errorf("update node group %s: %w", nodeGroupID, err)
	}
	return nil
}

func (c *Client) DeleteNodeGroup(ctx context.Context, clusterID, nodeGroupID string) error {
	if err := c.ensureClient(ctx); err != nil {
		return err
	}

	if err := c.httpClient.DELETE(ctx, "/v1/clusters/"+clusterID+"/nodegroups/"+nodeGroupID, nil); err != nil {
		return fmt.Errorf("delete node group %s: %w", nodeGroupID, err)
	}
	return nil
}

func (c *Client) ListClusterTemplates(ctx context.Context) (*ListClusterTemplatesOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	var out ListClusterTemplatesOutput
	if err := c.httpClient.GET(ctx, "/v1/clustertemplates", &out); err != nil {
		return nil, fmt.Errorf("list cluster templates: %w", err)
	}
	return &out, nil
}
