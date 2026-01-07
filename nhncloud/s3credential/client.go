// Package s3credential provides S3 Credential service client
package s3credential

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client represents a S3 Credential API client
type Client struct {
	baseURL    string
	token      string
	httpClient *http.Client
	debug      bool
}

// NewClient creates a new S3 Credential client
func NewClient(region, token string, httpClient *http.Client, debug bool) *Client {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 30 * time.Second}
	}
	baseURL := fmt.Sprintf("https://api-identity-infrastructure.%s.nhncloudservice.com", region)
	return &Client{
		baseURL:    baseURL,
		token:      token,
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

// ListCredentials lists all S3 credentials for a user
func (c *Client) ListCredentials(ctx context.Context, userID string) (*ListCredentialsOutput, error) {
	path := fmt.Sprintf("/v2.0/users/%s/credentials/OS-EC2", userID)
	data, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var credentials []S3Credential
	if err := json.Unmarshal(data, &credentials); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &ListCredentialsOutput{Credentials: credentials}, nil
}

// CreateCredential creates a new S3 credential
func (c *Client) CreateCredential(ctx context.Context, apiUserID string, input *CreateCredentialInput) (*CredentialOutput, error) {
	path := fmt.Sprintf("/v2.0/users/%s/credentials/OS-EC2", apiUserID)
	data, err := c.doRequest(ctx, "POST", path, input)
	if err != nil {
		return nil, err
	}

	var result CredentialOutput
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// DeleteCredential deletes an S3 credential
func (c *Client) DeleteCredential(ctx context.Context, userID, accessKey string) error {
	path := fmt.Sprintf("/v2.0/users/%s/credentials/OS-EC2/%s", userID, accessKey)
	_, err := c.doRequest(ctx, "DELETE", path, nil)
	return err
}
