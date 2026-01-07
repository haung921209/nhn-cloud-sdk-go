package object

import (
	"context"
	"fmt"
	"io"
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

	baseURL, err := c.tokenProvider.GetServiceEndpoint("object-store", c.region)
	if err != nil {
		return fmt.Errorf("resolve endpoint: %w", err)
	}

	opts := []client.ClientOption{
		client.WithDebug(c.debug),
	}

	c.httpClient = client.NewClient(baseURL, c.tokenProvider, opts...)
	return nil
}

func (c *Client) ListContainers(ctx context.Context) (*ListContainersOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	var containers []Container
	if err := c.httpClient.GET(ctx, "/?format=json", &containers); err != nil {
		return nil, fmt.Errorf("list containers: %w", err)
	}
	return &ListContainersOutput{Containers: containers}, nil
}

func (c *Client) CreateContainer(ctx context.Context, input *CreateContainerInput) error {
	if err := c.ensureClient(ctx); err != nil {
		return err
	}

	if err := c.httpClient.PUT(ctx, "/"+input.Name, nil, nil); err != nil {
		return fmt.Errorf("create container %s: %w", input.Name, err)
	}
	return nil
}

func (c *Client) DeleteContainer(ctx context.Context, containerName string) error {
	if err := c.ensureClient(ctx); err != nil {
		return err
	}

	if err := c.httpClient.DELETE(ctx, "/"+containerName, nil); err != nil {
		return fmt.Errorf("delete container %s: %w", containerName, err)
	}
	return nil
}

func (c *Client) ListObjects(ctx context.Context, containerName string, input *ListObjectsInput) (*ListObjectsOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	path := "/" + containerName + "?format=json"
	if input != nil {
		if input.Prefix != "" {
			path += "&prefix=" + input.Prefix
		}
		if input.Delimiter != "" {
			path += "&delimiter=" + input.Delimiter
		}
		if input.Marker != "" {
			path += "&marker=" + input.Marker
		}
		if input.Limit > 0 {
			path += fmt.Sprintf("&limit=%d", input.Limit)
		}
	}

	var objects []Object
	if err := c.httpClient.GET(ctx, path, &objects); err != nil {
		return nil, fmt.Errorf("list objects in %s: %w", containerName, err)
	}
	return &ListObjectsOutput{Objects: objects}, nil
}

func (c *Client) PutObject(ctx context.Context, containerName, objectName string, body io.Reader, metadata *ObjectMetadata) error {
	if err := c.ensureClient(ctx); err != nil {
		return err
	}

	path := "/" + containerName + "/" + objectName
	if err := c.httpClient.PUT(ctx, path, body, nil); err != nil {
		return fmt.Errorf("put object %s/%s: %w", containerName, objectName, err)
	}
	return nil
}

func (c *Client) GetObject(ctx context.Context, containerName, objectName string) (*GetObjectOutput, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	path := "/" + containerName + "/" + objectName
	var out GetObjectOutput
	if err := c.httpClient.GET(ctx, path, &out); err != nil {
		return nil, fmt.Errorf("get object %s/%s: %w", containerName, objectName, err)
	}
	return &out, nil
}

func (c *Client) DeleteObject(ctx context.Context, containerName, objectName string) error {
	if err := c.ensureClient(ctx); err != nil {
		return err
	}

	path := "/" + containerName + "/" + objectName
	if err := c.httpClient.DELETE(ctx, path, nil); err != nil {
		return fmt.Errorf("delete object %s/%s: %w", containerName, objectName, err)
	}
	return nil
}

func (c *Client) CopyObject(ctx context.Context, input *CopyObjectInput) error {
	if err := c.ensureClient(ctx); err != nil {
		return err
	}

	destPath := "/" + input.DestinationContainer + "/" + input.DestinationObjectName
	if err := c.httpClient.PUT(ctx, destPath, nil, nil); err != nil {
		return fmt.Errorf("copy object to %s/%s: %w", input.DestinationContainer, input.DestinationObjectName, err)
	}
	return nil
}
