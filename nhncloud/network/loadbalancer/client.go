package loadbalancer

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

	baseURL, err := c.tokenProvider.GetServiceEndpoint("network", c.region)
	if err != nil {
		return err
	}

	opts := []client.ClientOption{
		client.WithDebug(c.debug),
	}

	c.httpClient = client.NewClient(baseURL, c.tokenProvider, opts...)
	return nil
}

func (c *Client) ListLoadBalancers(ctx context.Context) (*ListLoadBalancersOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	var out ListLoadBalancersOutput
	if err := c.httpClient.GET(ctx, "/v2.0/lbaas/loadbalancers", &out); err != nil {
		return nil, fmt.Errorf("list load balancers: %w", err)
	}
	return &out, nil
}

func (c *Client) GetLoadBalancer(ctx context.Context, lbID string) (*GetLoadBalancerOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	var out GetLoadBalancerOutput
	if err := c.httpClient.GET(ctx, "/v2.0/lbaas/loadbalancers/"+lbID, &out); err != nil {
		return nil, fmt.Errorf("get load balancer %s: %w", lbID, err)
	}
	return &out, nil
}

func (c *Client) CreateLoadBalancer(ctx context.Context, input *CreateLoadBalancerInput) (*CreateLoadBalancerOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	req := &CreateLoadBalancerRequest{LoadBalancer: *input}
	var out CreateLoadBalancerOutput
	if err := c.httpClient.POST(ctx, "/v2.0/lbaas/loadbalancers", req, &out); err != nil {
		return nil, fmt.Errorf("create load balancer: %w", err)
	}
	return &out, nil
}

func (c *Client) UpdateLoadBalancer(ctx context.Context, lbID string, input *UpdateLoadBalancerInput) (*GetLoadBalancerOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	req := &UpdateLoadBalancerRequest{LoadBalancer: *input}
	var out GetLoadBalancerOutput
	if err := c.httpClient.PUT(ctx, "/v2.0/lbaas/loadbalancers/"+lbID, req, &out); err != nil {
		return nil, fmt.Errorf("update load balancer %s: %w", lbID, err)
	}
	return &out, nil
}

func (c *Client) DeleteLoadBalancer(ctx context.Context, lbID string) error {
	if err := c.ensureClient(ctx); err != nil {
		return err
	}

	if err := c.httpClient.DELETE(ctx, "/v2.0/lbaas/loadbalancers/"+lbID, nil); err != nil {
		return fmt.Errorf("delete load balancer %s: %w", lbID, err)
	}
	return nil
}

func (c *Client) ListListeners(ctx context.Context) (*ListListenersOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	var out ListListenersOutput
	if err := c.httpClient.GET(ctx, "/v2.0/lbaas/listeners", &out); err != nil {
		return nil, fmt.Errorf("list listeners: %w", err)
	}
	return &out, nil
}

func (c *Client) GetListener(ctx context.Context, listenerID string) (*GetListenerOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	var out GetListenerOutput
	if err := c.httpClient.GET(ctx, "/v2.0/lbaas/listeners/"+listenerID, &out); err != nil {
		return nil, fmt.Errorf("get listener %s: %w", listenerID, err)
	}
	return &out, nil
}

func (c *Client) CreateListener(ctx context.Context, input *CreateListenerInput) (*GetListenerOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	req := &CreateListenerRequest{Listener: *input}
	var out GetListenerOutput
	if err := c.httpClient.POST(ctx, "/v2.0/lbaas/listeners", req, &out); err != nil {
		return nil, fmt.Errorf("create listener: %w", err)
	}
	return &out, nil
}

func (c *Client) DeleteListener(ctx context.Context, listenerID string) error {
	if err := c.ensureClient(ctx); err != nil {
		return err
	}

	if err := c.httpClient.DELETE(ctx, "/v2.0/lbaas/listeners/"+listenerID, nil); err != nil {
		return fmt.Errorf("delete listener %s: %w", listenerID, err)
	}
	return nil
}

func (c *Client) ListPools(ctx context.Context) (*ListPoolsOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	var out ListPoolsOutput
	if err := c.httpClient.GET(ctx, "/v2.0/lbaas/pools", &out); err != nil {
		return nil, fmt.Errorf("list pools: %w", err)
	}
	return &out, nil
}

func (c *Client) GetPool(ctx context.Context, poolID string) (*GetPoolOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	var out GetPoolOutput
	if err := c.httpClient.GET(ctx, "/v2.0/lbaas/pools/"+poolID, &out); err != nil {
		return nil, fmt.Errorf("get pool %s: %w", poolID, err)
	}
	return &out, nil
}

func (c *Client) CreatePool(ctx context.Context, input *CreatePoolInput) (*GetPoolOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	req := &CreatePoolRequest{Pool: *input}
	var out GetPoolOutput
	if err := c.httpClient.POST(ctx, "/v2.0/lbaas/pools", req, &out); err != nil {
		return nil, fmt.Errorf("create pool: %w", err)
	}
	return &out, nil
}

func (c *Client) DeletePool(ctx context.Context, poolID string) error {
	if err := c.ensureClient(ctx); err != nil {
		return err
	}

	if err := c.httpClient.DELETE(ctx, "/v2.0/lbaas/pools/"+poolID, nil); err != nil {
		return fmt.Errorf("delete pool %s: %w", poolID, err)
	}
	return nil
}

func (c *Client) ListMembers(ctx context.Context, poolID string) (*ListMembersOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	var out ListMembersOutput
	if err := c.httpClient.GET(ctx, "/v2.0/lbaas/pools/"+poolID+"/members", &out); err != nil {
		return nil, fmt.Errorf("list members: %w", err)
	}
	return &out, nil
}

func (c *Client) GetMember(ctx context.Context, poolID, memberID string) (*GetMemberOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	var out GetMemberOutput
	if err := c.httpClient.GET(ctx, "/v2.0/lbaas/pools/"+poolID+"/members/"+memberID, &out); err != nil {
		return nil, fmt.Errorf("get member %s: %w", memberID, err)
	}
	return &out, nil
}

func (c *Client) CreateMember(ctx context.Context, poolID string, input *CreateMemberInput) (*GetMemberOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	req := &CreateMemberRequest{Member: *input}
	var out GetMemberOutput
	if err := c.httpClient.POST(ctx, "/v2.0/lbaas/pools/"+poolID+"/members", req, &out); err != nil {
		return nil, fmt.Errorf("create member: %w", err)
	}
	return &out, nil
}

func (c *Client) DeleteMember(ctx context.Context, poolID, memberID string) error {
	if err := c.ensureClient(ctx); err != nil {
		return err
	}

	if err := c.httpClient.DELETE(ctx, "/v2.0/lbaas/pools/"+poolID+"/members/"+memberID, nil); err != nil {
		return fmt.Errorf("delete member %s: %w", memberID, err)
	}
	return nil
}

func (c *Client) ListHealthMonitors(ctx context.Context) (*ListHealthMonitorsOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	var out ListHealthMonitorsOutput
	if err := c.httpClient.GET(ctx, "/v2.0/lbaas/healthmonitors", &out); err != nil {
		return nil, fmt.Errorf("list health monitors: %w", err)
	}
	return &out, nil
}

func (c *Client) GetHealthMonitor(ctx context.Context, monitorID string) (*GetHealthMonitorOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	var out GetHealthMonitorOutput
	if err := c.httpClient.GET(ctx, "/v2.0/lbaas/healthmonitors/"+monitorID, &out); err != nil {
		return nil, fmt.Errorf("get health monitor %s: %w", monitorID, err)
	}
	return &out, nil
}

func (c *Client) CreateHealthMonitor(ctx context.Context, input *CreateHealthMonitorInput) (*GetHealthMonitorOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	req := &CreateHealthMonitorRequest{HealthMonitor: *input}
	var out GetHealthMonitorOutput
	if err := c.httpClient.POST(ctx, "/v2.0/lbaas/healthmonitors", req, &out); err != nil {
		return nil, fmt.Errorf("create health monitor: %w", err)
	}
	return &out, nil
}

func (c *Client) DeleteHealthMonitor(ctx context.Context, monitorID string) error {
	if err := c.ensureClient(ctx); err != nil {
		return err
	}

	if err := c.httpClient.DELETE(ctx, "/v2.0/lbaas/healthmonitors/"+monitorID, nil); err != nil {
		return fmt.Errorf("delete health monitor %s: %w", monitorID, err)
	}
	return nil
}
