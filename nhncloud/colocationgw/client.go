// Package colocationgw provides Colocation Gateway service client
package colocationgw

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client represents a Colocation Gateway API client
type Client struct {
	baseURL    string
	token      string
	tenantID   string
	httpClient *http.Client
	debug      bool
}

// NewClient creates a new Colocation Gateway client
func NewClient(region, token, tenantID string, httpClient *http.Client, debug bool) *Client {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 30 * time.Second}
	}
	baseURL := fmt.Sprintf("https://kr1-api-network-infrastructure.nhncloudservice.com")
	return &Client{
		baseURL:    baseURL,
		token:      token,
		tenantID:   tenantID,
		httpClient: httpClient,
		debug:      debug,
	}
}

// doRequest performs an HTTP request
func (c *Client) doRequest(ctx context.Context, method, path string) ([]byte, error) {
	fullURL := c.baseURL + path
	if c.debug {
		fmt.Printf("[DEBUG] %s %s\n", method, fullURL)
	}

	req, err := http.NewRequestWithContext(ctx, method, fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Auth-Token", c.token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if c.debug {
		fmt.Printf("[DEBUG] Response status: %d\n", resp.StatusCode)
		fmt.Printf("[DEBUG] Response body: %s\n", string(respBody))
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

// List lists all colocation gateways
func (c *Client) List(ctx context.Context) (*ListOutput, error) {
	data, err := c.doRequest(ctx, "GET", "/v2.0/gateways/colocationgateways")
	if err != nil {
		return nil, err
	}

	var result ListOutput
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// Get gets a colocation gateway by ID
func (c *Client) Get(ctx context.Context, gatewayID string) (*GetOutput, error) {
	path := fmt.Sprintf("/v2.0/gateways/colocationgateways/%s", gatewayID)
	data, err := c.doRequest(ctx, "GET", path)
	if err != nil {
		return nil, err
	}

	var result GetOutput
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}
