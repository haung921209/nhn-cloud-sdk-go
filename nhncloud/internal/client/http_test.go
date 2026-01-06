package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockTokenProvider struct {
	token string
}

func (m *mockTokenProvider) GetToken(ctx context.Context) (string, error) {
	return m.token, nil
}

func (m *mockTokenProvider) SetAuthHeader(req *http.Request, token string) {
	req.Header.Set("Authorization", "Bearer "+token)
}

func TestClientGET(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/test" {
			t.Errorf("expected /test, got %s", r.URL.Path)
		}
		if r.Header.Get("Authorization") != "Bearer test-token" {
			t.Errorf("expected Bearer test-token, got %s", r.Header.Get("Authorization"))
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": "success"})
	}))
	defer server.Close()

	client := NewClient(server.URL, &mockTokenProvider{token: "test-token"})

	var result map[string]string
	err := client.GET(context.Background(), "/test", &result)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result["message"] != "success" {
		t.Errorf("expected success, got %s", result["message"])
	}
}

func TestClientPOST(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("expected application/json, got %s", r.Header.Get("Content-Type"))
		}

		var body map[string]string
		json.NewDecoder(r.Body).Decode(&body)

		if body["name"] != "test" {
			t.Errorf("expected test, got %s", body["name"])
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"id": "123", "name": body["name"]})
	}))
	defer server.Close()

	client := NewClient(server.URL, &mockTokenProvider{token: "test-token"})

	var result map[string]string
	err := client.POST(context.Background(), "/resources", map[string]string{"name": "test"}, &result)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result["id"] != "123" {
		t.Errorf("expected 123, got %s", result["id"])
	}
}

func TestClientErrorResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"message": "resource not found", "error_code": "NOT_FOUND"})
	}))
	defer server.Close()

	client := NewClient(server.URL, &mockTokenProvider{token: "test-token"})

	var result map[string]string
	err := client.GET(context.Background(), "/notfound", &result)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected APIError, got %T", err)
	}

	if apiErr.StatusCode != 404 {
		t.Errorf("expected 404, got %d", apiErr.StatusCode)
	}

	if apiErr.Message != "resource not found" {
		t.Errorf("expected 'resource not found', got %s", apiErr.Message)
	}

	if apiErr.ErrorCode != "NOT_FOUND" {
		t.Errorf("expected 'NOT_FOUND', got %s", apiErr.ErrorCode)
	}
}

func TestClientNHNCloudErrorResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"header": map[string]interface{}{
				"resultCode":    400,
				"resultMessage": "Invalid parameter",
				"isSuccessful":  false,
			},
		})
	}))
	defer server.Close()

	client := NewClient(server.URL, &mockTokenProvider{token: "test-token"})

	var result map[string]string
	err := client.POST(context.Background(), "/invalid", nil, &result)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected APIError, got %T", err)
	}

	if apiErr.Message != "Invalid parameter" {
		t.Errorf("expected 'Invalid parameter', got %s", apiErr.Message)
	}
}

func TestClientWithOptions(t *testing.T) {
	customHTTPClient := &http.Client{}
	client := NewClient(
		"https://api.example.com",
		nil,
		WithHTTPClient(customHTTPClient),
		WithDebug(true),
		WithUserAgent("test-agent/1.0"),
	)

	if client.HTTPClient != customHTTPClient {
		t.Error("HTTPClient option not applied")
	}

	if !client.Debug {
		t.Error("Debug option not applied")
	}

	if client.UserAgent != "test-agent/1.0" {
		t.Errorf("UserAgent expected 'test-agent/1.0', got %s", client.UserAgent)
	}
}

func TestClientDELETE(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := NewClient(server.URL, &mockTokenProvider{token: "test-token"})

	err := client.DELETE(context.Background(), "/resource/123", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestClientPUT(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}

		var body map[string]string
		json.NewDecoder(r.Body).Decode(&body)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(body)
	}))
	defer server.Close()

	client := NewClient(server.URL, &mockTokenProvider{token: "test-token"})

	var result map[string]string
	err := client.PUT(context.Background(), "/resource/123", map[string]string{"name": "updated"}, &result)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result["name"] != "updated" {
		t.Errorf("expected 'updated', got %s", result["name"])
	}
}

func TestAPIErrorString(t *testing.T) {
	tests := []struct {
		name     string
		err      *APIError
		expected string
	}{
		{
			name:     "with error code",
			err:      &APIError{StatusCode: 404, Message: "not found", ErrorCode: "NOT_FOUND"},
			expected: "API Error 404 [NOT_FOUND]: not found",
		},
		{
			name:     "without error code",
			err:      &APIError{StatusCode: 500, Message: "internal error"},
			expected: "API Error 500: internal error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.expected {
				t.Errorf("Error() = %v, want %v", got, tt.expected)
			}
		})
	}
}
