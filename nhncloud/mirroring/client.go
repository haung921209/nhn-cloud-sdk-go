// Package mirroring provides Traffic Mirroring service client
package mirroring

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client represents a Traffic Mirroring API client
type Client struct {
	baseURL    string
	token      string
	tenantID   string
	httpClient *http.Client
	debug      bool
}

// NewClient creates a new Traffic Mirroring client
func NewClient(region, token, tenantID string, httpClient *http.Client, debug bool) *Client {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 30 * time.Second}
	}
	baseURL := fmt.Sprintf("https://%s-api-network-infrastructure.nhncloudservice.com", region)
	return &Client{
		baseURL:    baseURL,
		token:      token,
		tenantID:   tenantID,
		httpClient: httpClient,
		debug:      debug,
	}
}

// doRequest performs an HTTP request
func (c *Client) doRequest(ctx context.Context, method, path string, body interface{}) ([]byte, error) {
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewReader(jsonData)
		if c.debug {
			fmt.Printf("[DEBUG] Request body: %s\n", string(jsonData))
		}
	}

	fullURL := c.baseURL + path
	if c.debug {
		fmt.Printf("[DEBUG] %s %s\n", method, fullURL)
	}

	req, err := http.NewRequestWithContext(ctx, method, fullURL, reqBody)
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

// ================================
// Session Operations
// ================================

// ListSessions lists all mirroring sessions
func (c *Client) ListSessions(ctx context.Context) (*ListSessionsOutput, error) {
	data, err := c.doRequest(ctx, "GET", "/v2.0/mirroring/sessions", nil)
	if err != nil {
		return nil, err
	}

	var result ListSessionsOutput
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// GetSession gets a session by ID
func (c *Client) GetSession(ctx context.Context, sessionID string) (*SessionOutput, error) {
	path := fmt.Sprintf("/v2.0/mirroring/sessions/%s", sessionID)
	data, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var result SessionOutput
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// CreateSession creates a new session
func (c *Client) CreateSession(ctx context.Context, input *CreateSessionInput) (*SessionOutput, error) {
	request := map[string]interface{}{
		"mirroring_session": input,
	}
	data, err := c.doRequest(ctx, "POST", "/v2.0/mirroring/sessions", request)
	if err != nil {
		return nil, err
	}

	var result SessionOutput
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// UpdateSession updates a session
func (c *Client) UpdateSession(ctx context.Context, sessionID string, input *UpdateSessionInput) (*SessionOutput, error) {
	path := fmt.Sprintf("/v2.0/mirroring/sessions/%s", sessionID)
	request := map[string]interface{}{
		"mirroring_session": input,
	}
	data, err := c.doRequest(ctx, "PUT", path, request)
	if err != nil {
		return nil, err
	}

	var result SessionOutput
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// DeleteSession deletes a session
func (c *Client) DeleteSession(ctx context.Context, sessionID string) error {
	path := fmt.Sprintf("/v2.0/mirroring/sessions/%s", sessionID)
	_, err := c.doRequest(ctx, "DELETE", path, nil)
	return err
}

// ================================
// Filter Group Operations
// ================================

// ListFilterGroups lists all filter groups
func (c *Client) ListFilterGroups(ctx context.Context) (*ListFilterGroupsOutput, error) {
	data, err := c.doRequest(ctx, "GET", "/v2.0/mirroring/filtergroups", nil)
	if err != nil {
		return nil, err
	}

	var result ListFilterGroupsOutput
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// GetFilterGroup gets a filter group by ID
func (c *Client) GetFilterGroup(ctx context.Context, filterGroupID string) (*FilterGroupOutput, error) {
	path := fmt.Sprintf("/v2.0/mirroring/filtergroups/%s", filterGroupID)
	data, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var result FilterGroupOutput
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// CreateFilterGroup creates a new filter group
func (c *Client) CreateFilterGroup(ctx context.Context, input *CreateFilterGroupInput) (*FilterGroupOutput, error) {
	request := map[string]interface{}{
		"mirroring_filtergroup": input,
	}
	data, err := c.doRequest(ctx, "POST", "/v2.0/mirroring/filtergroups", request)
	if err != nil {
		return nil, err
	}

	var result FilterGroupOutput
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// UpdateFilterGroup updates a filter group
func (c *Client) UpdateFilterGroup(ctx context.Context, filterGroupID string, input *UpdateFilterGroupInput) (*FilterGroupOutput, error) {
	path := fmt.Sprintf("/v2.0/mirroring/filtergroups/%s", filterGroupID)
	request := map[string]interface{}{
		"mirroring_filtergroup": input,
	}
	data, err := c.doRequest(ctx, "PUT", path, request)
	if err != nil {
		return nil, err
	}

	var result FilterGroupOutput
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// DeleteFilterGroup deletes a filter group
func (c *Client) DeleteFilterGroup(ctx context.Context, filterGroupID string) error {
	path := fmt.Sprintf("/v2.0/mirroring/filtergroups/%s", filterGroupID)
	_, err := c.doRequest(ctx, "DELETE", path, nil)
	return err
}

// ================================
// Filter Operations
// ================================

// ListFilters lists all filters
func (c *Client) ListFilters(ctx context.Context) (*ListFiltersOutput, error) {
	data, err := c.doRequest(ctx, "GET", "/v2.0/mirroring/filters", nil)
	if err != nil {
		return nil, err
	}

	var result ListFiltersOutput
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// GetFilter gets a filter by ID
func (c *Client) GetFilter(ctx context.Context, filterID string) (*FilterOutput, error) {
	path := fmt.Sprintf("/v2.0/mirroring/filters/%s", filterID)
	data, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var result FilterOutput
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// CreateFilter creates a new filter
func (c *Client) CreateFilter(ctx context.Context, input *CreateFilterInput) (*FilterOutput, error) {
	request := map[string]interface{}{
		"mirroring_filter": input,
	}
	data, err := c.doRequest(ctx, "POST", "/v2.0/mirroring/filters", request)
	if err != nil {
		return nil, err
	}

	var result FilterOutput
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// UpdateFilter updates a filter
func (c *Client) UpdateFilter(ctx context.Context, filterID string, input *UpdateFilterInput) (*FilterOutput, error) {
	path := fmt.Sprintf("/v2.0/mirroring/filters/%s", filterID)
	request := map[string]interface{}{
		"mirroring_filter": input,
	}
	data, err := c.doRequest(ctx, "PUT", path, request)
	if err != nil {
		return nil, err
	}

	var result FilterOutput
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// DeleteFilter deletes a filter
func (c *Client) DeleteFilter(ctx context.Context, filterID string) error {
	path := fmt.Sprintf("/v2.0/mirroring/filters/%s", filterID)
	_, err := c.doRequest(ctx, "DELETE", path, nil)
	return err
}
