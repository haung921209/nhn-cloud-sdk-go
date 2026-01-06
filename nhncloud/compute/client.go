package compute

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

	if _, err := c.tokenProvider.GetToken(ctx); err != nil {
		return fmt.Errorf("authenticate: %w", err)
	}

	baseURL, err := c.tokenProvider.GetServiceEndpoint("compute", c.region)
	if err != nil {
		return fmt.Errorf("resolve endpoint: %w", err)
	}

	opts := []client.ClientOption{
		client.WithDebug(c.debug),
	}

	c.httpClient = client.NewClient(baseURL, c.tokenProvider, opts...)
	return nil
}

func (c *Client) ListServers(ctx context.Context) (*ListServersOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	var out ListServersOutput
	if err := c.httpClient.GET(ctx, "/servers/detail", &out); err != nil {
		return nil, fmt.Errorf("list servers: %w", err)
	}
	return &out, nil
}

func (c *Client) GetServer(ctx context.Context, serverID string) (*GetServerOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	var out GetServerOutput
	if err := c.httpClient.GET(ctx, "/servers/"+serverID, &out); err != nil {
		return nil, fmt.Errorf("get server %s: %w", serverID, err)
	}
	return &out, nil
}

func (c *Client) CreateServer(ctx context.Context, input *CreateServerInput) (*CreateServerOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	req := map[string]interface{}{"server": input}
	var out CreateServerOutput
	if err := c.httpClient.POST(ctx, "/servers", req, &out); err != nil {
		return nil, fmt.Errorf("create server: %w", err)
	}
	return &out, nil
}

func (c *Client) DeleteServer(ctx context.Context, serverID string) error {
	if err := c.ensureClient(ctx); err != nil {
		return err
	}

	if err := c.httpClient.DELETE(ctx, "/servers/"+serverID, nil); err != nil {
		return fmt.Errorf("delete server %s: %w", serverID, err)
	}
	return nil
}

func (c *Client) StartServer(ctx context.Context, serverID string) error {
	if err := c.ensureClient(ctx); err != nil {
		return err
	}

	req := map[string]interface{}{"os-start": nil}
	if err := c.httpClient.POST(ctx, "/servers/"+serverID+"/action", req, nil); err != nil {
		return fmt.Errorf("start server %s: %w", serverID, err)
	}
	return nil
}

func (c *Client) StopServer(ctx context.Context, serverID string) error {
	if err := c.ensureClient(ctx); err != nil {
		return err
	}

	req := map[string]interface{}{"os-stop": nil}
	if err := c.httpClient.POST(ctx, "/servers/"+serverID+"/action", req, nil); err != nil {
		return fmt.Errorf("stop server %s: %w", serverID, err)
	}
	return nil
}

func (c *Client) RebootServer(ctx context.Context, serverID string, hard bool) error {
	if err := c.ensureClient(ctx); err != nil {
		return err
	}

	rebootType := "SOFT"
	if hard {
		rebootType = "HARD"
	}
	req := map[string]interface{}{
		"reboot": map[string]string{"type": rebootType},
	}
	if err := c.httpClient.POST(ctx, "/servers/"+serverID+"/action", req, nil); err != nil {
		return fmt.Errorf("reboot server %s: %w", serverID, err)
	}
	return nil
}

func (c *Client) ResizeServer(ctx context.Context, serverID, flavorRef string) error {
	if err := c.ensureClient(ctx); err != nil {
		return err
	}

	req := map[string]interface{}{
		"resize": map[string]string{"flavorRef": flavorRef},
	}
	if err := c.httpClient.POST(ctx, "/servers/"+serverID+"/action", req, nil); err != nil {
		return fmt.Errorf("resize server %s: %w", serverID, err)
	}
	return nil
}

func (c *Client) ConfirmResize(ctx context.Context, serverID string) error {
	if err := c.ensureClient(ctx); err != nil {
		return err
	}

	req := map[string]interface{}{"confirmResize": nil}
	if err := c.httpClient.POST(ctx, "/servers/"+serverID+"/action", req, nil); err != nil {
		return fmt.Errorf("confirm resize %s: %w", serverID, err)
	}
	return nil
}

func (c *Client) ListFlavors(ctx context.Context) (*ListFlavorsOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	var out ListFlavorsOutput
	if err := c.httpClient.GET(ctx, "/flavors/detail", &out); err != nil {
		return nil, fmt.Errorf("list flavors: %w", err)
	}
	return &out, nil
}

func (c *Client) ListImages(ctx context.Context) (*ListImagesOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	var out ListImagesOutput
	if err := c.httpClient.GET(ctx, "/images/detail", &out); err != nil {
		return nil, fmt.Errorf("list images: %w", err)
	}
	return &out, nil
}

func (c *Client) ListKeyPairs(ctx context.Context) (*ListKeyPairsOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	var out ListKeyPairsOutput
	if err := c.httpClient.GET(ctx, "/os-keypairs", &out); err != nil {
		return nil, fmt.Errorf("list keypairs: %w", err)
	}
	return &out, nil
}

func (c *Client) CreateKeyPair(ctx context.Context, input *CreateKeyPairInput) (*CreateKeyPairOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	req := map[string]interface{}{"keypair": input}
	var out CreateKeyPairOutput
	if err := c.httpClient.POST(ctx, "/os-keypairs", req, &out); err != nil {
		return nil, fmt.Errorf("create keypair: %w", err)
	}
	return &out, nil
}

func (c *Client) DeleteKeyPair(ctx context.Context, name string) error {
	if err := c.ensureClient(ctx); err != nil {
		return err
	}

	if err := c.httpClient.DELETE(ctx, "/os-keypairs/"+name, nil); err != nil {
		return fmt.Errorf("delete keypair %s: %w", name, err)
	}
	return nil
}
