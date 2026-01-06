package iam

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
	credentials   credentials.Credentials
	httpClient    *client.Client
	tokenProvider *client.OAuthTokenProvider
	debug         bool
}

func NewClient(region string, creds credentials.Credentials, hc *http.Client, debug bool) *Client {
	c := &Client{
		region:      region,
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
	baseURL := endpoint.Resolve(endpoint.ServiceIAM, c.region)
	opts := []client.ClientOption{
		client.WithDebug(c.debug),
	}
	c.httpClient = client.NewClient(baseURL, c.tokenProvider, opts...)
}

func (c *Client) ListOrganizations(ctx context.Context) (*ListOrganizationsOutput, error) {
	var out ListOrganizationsOutput
	if err := c.httpClient.GET(ctx, "/v1/organizations", &out); err != nil {
		return nil, fmt.Errorf("list organizations: %w", err)
	}
	return &out, nil
}

func (c *Client) GetOrganization(ctx context.Context, orgID string) (*GetOrganizationOutput, error) {
	var out GetOrganizationOutput
	if err := c.httpClient.GET(ctx, "/v1/organizations/"+orgID, &out); err != nil {
		return nil, fmt.Errorf("get organization %s: %w", orgID, err)
	}
	return &out, nil
}

func (c *Client) ListProjects(ctx context.Context, orgID string) (*ListProjectsOutput, error) {
	var out ListProjectsOutput
	path := fmt.Sprintf("/v1/organizations/%s/projects", orgID)
	if err := c.httpClient.GET(ctx, path, &out); err != nil {
		return nil, fmt.Errorf("list projects: %w", err)
	}
	return &out, nil
}

func (c *Client) GetProject(ctx context.Context, orgID, projectID string) (*GetProjectOutput, error) {
	var out GetProjectOutput
	path := fmt.Sprintf("/v1/organizations/%s/projects/%s", orgID, projectID)
	if err := c.httpClient.GET(ctx, path, &out); err != nil {
		return nil, fmt.Errorf("get project %s: %w", projectID, err)
	}
	return &out, nil
}

func (c *Client) ListMembers(ctx context.Context, orgID string) (*ListMembersOutput, error) {
	var out ListMembersOutput
	path := fmt.Sprintf("/v1/organizations/%s/members", orgID)
	if err := c.httpClient.GET(ctx, path, &out); err != nil {
		return nil, fmt.Errorf("list members: %w", err)
	}
	return &out, nil
}

func (c *Client) GetMember(ctx context.Context, orgID, memberID string) (*GetMemberOutput, error) {
	var out GetMemberOutput
	path := fmt.Sprintf("/v1/organizations/%s/members/%s", orgID, memberID)
	if err := c.httpClient.GET(ctx, path, &out); err != nil {
		return nil, fmt.Errorf("get member %s: %w", memberID, err)
	}
	return &out, nil
}

func (c *Client) InviteMember(ctx context.Context, orgID string, input *InviteMemberInput) (*InviteMemberOutput, error) {
	var out InviteMemberOutput
	path := fmt.Sprintf("/v1/organizations/%s/members", orgID)
	if err := c.httpClient.POST(ctx, path, input, &out); err != nil {
		return nil, fmt.Errorf("invite member: %w", err)
	}
	return &out, nil
}

func (c *Client) RemoveMember(ctx context.Context, orgID, memberID string) error {
	path := fmt.Sprintf("/v1/organizations/%s/members/%s", orgID, memberID)
	if err := c.httpClient.DELETE(ctx, path, nil); err != nil {
		return fmt.Errorf("remove member %s: %w", memberID, err)
	}
	return nil
}
