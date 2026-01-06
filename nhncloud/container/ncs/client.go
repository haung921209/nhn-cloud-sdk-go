package ncs

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
	baseURL := endpoint.ResolveWithAppKey(endpoint.ServiceNCS, c.region, c.appKey)
	opts := []client.ClientOption{
		client.WithDebug(c.debug),
	}
	c.httpClient = client.NewClient(baseURL, c.tokenProvider, opts...)
}

func (c *Client) ListWorkloads(ctx context.Context, namespace string) (*ListWorkloadsOutput, error) {
	path := "/workloads"
	if namespace != "" {
		path += "?namespace=" + namespace
	}
	var out ListWorkloadsOutput
	if err := c.httpClient.GET(ctx, path, &out); err != nil {
		return nil, fmt.Errorf("list workloads: %w", err)
	}
	return &out, nil
}

func (c *Client) GetWorkload(ctx context.Context, workloadID string) (*GetWorkloadOutput, error) {
	var out GetWorkloadOutput
	if err := c.httpClient.GET(ctx, "/workloads/"+workloadID, &out); err != nil {
		return nil, fmt.Errorf("get workload %s: %w", workloadID, err)
	}
	return &out, nil
}

func (c *Client) CreateWorkload(ctx context.Context, input *CreateWorkloadInput) (*CreateWorkloadOutput, error) {
	var out CreateWorkloadOutput
	if err := c.httpClient.POST(ctx, "/workloads", input, &out); err != nil {
		return nil, fmt.Errorf("create workload: %w", err)
	}
	return &out, nil
}

func (c *Client) UpdateWorkload(ctx context.Context, workloadID string, input *UpdateWorkloadInput) (*GetWorkloadOutput, error) {
	var out GetWorkloadOutput
	if err := c.httpClient.PUT(ctx, "/workloads/"+workloadID, input, &out); err != nil {
		return nil, fmt.Errorf("update workload %s: %w", workloadID, err)
	}
	return &out, nil
}

func (c *Client) DeleteWorkload(ctx context.Context, workloadID string) error {
	if err := c.httpClient.DELETE(ctx, "/workloads/"+workloadID, nil); err != nil {
		return fmt.Errorf("delete workload %s: %w", workloadID, err)
	}
	return nil
}

func (c *Client) RestartWorkload(ctx context.Context, workloadID string) error {
	if err := c.httpClient.POST(ctx, "/workloads/"+workloadID+"/restart", nil, nil); err != nil {
		return fmt.Errorf("restart workload %s: %w", workloadID, err)
	}
	return nil
}

func (c *Client) ScaleWorkload(ctx context.Context, workloadID string, replicas int) error {
	req := map[string]int{"replicas": replicas}
	if err := c.httpClient.POST(ctx, "/workloads/"+workloadID+"/scale", req, nil); err != nil {
		return fmt.Errorf("scale workload %s: %w", workloadID, err)
	}
	return nil
}

func (c *Client) ListTemplates(ctx context.Context) (*ListTemplatesOutput, error) {
	var out ListTemplatesOutput
	if err := c.httpClient.GET(ctx, "/templates", &out); err != nil {
		return nil, fmt.Errorf("list templates: %w", err)
	}
	return &out, nil
}

func (c *Client) GetTemplate(ctx context.Context, templateID string) (*GetTemplateOutput, error) {
	var out GetTemplateOutput
	if err := c.httpClient.GET(ctx, "/templates/"+templateID, &out); err != nil {
		return nil, fmt.Errorf("get template %s: %w", templateID, err)
	}
	return &out, nil
}

func (c *Client) ListServices(ctx context.Context, namespace string) (*ListServicesOutput, error) {
	path := "/services"
	if namespace != "" {
		path += "?namespace=" + namespace
	}
	var out ListServicesOutput
	if err := c.httpClient.GET(ctx, path, &out); err != nil {
		return nil, fmt.Errorf("list services: %w", err)
	}
	return &out, nil
}

func (c *Client) GetService(ctx context.Context, serviceID string) (*GetServiceOutput, error) {
	var out GetServiceOutput
	if err := c.httpClient.GET(ctx, "/services/"+serviceID, &out); err != nil {
		return nil, fmt.Errorf("get service %s: %w", serviceID, err)
	}
	return &out, nil
}

func (c *Client) CreateService(ctx context.Context, input *CreateServiceInput) (*CreateServiceOutput, error) {
	var out CreateServiceOutput
	if err := c.httpClient.POST(ctx, "/services", input, &out); err != nil {
		return nil, fmt.Errorf("create service: %w", err)
	}
	return &out, nil
}

func (c *Client) DeleteService(ctx context.Context, serviceID string) error {
	if err := c.httpClient.DELETE(ctx, "/services/"+serviceID, nil); err != nil {
		return fmt.Errorf("delete service %s: %w", serviceID, err)
	}
	return nil
}
