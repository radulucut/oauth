package oauth

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestClient_Facebook_Success(t *testing.T) {
	// Create a mock server that returns a successful Facebook response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify the request
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		if r.URL.Query().Get("access_token") != "test_token" {
			t.Errorf("Expected access_token=test_token, got %s", r.URL.Query().Get("access_token"))
		}
		if r.URL.Query().Get("fields") != "email,name" {
			t.Errorf("Expected fields=email,name, got %s", r.URL.Query().Get("fields"))
		}

		// Return mock Facebook user data
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"id": "123456789",
			"name": "John Doe",
			"email": "john.doe@example.com"
		}`))
	}))
	defer server.Close()

	// Create client with mock server URL
	c := &client{
		facebookURL: server.URL,
		httpClient:  &http.Client{Timeout: 10 * time.Second},
	}

	// Test the Facebook function
	ctx := context.Background()
	payload, err := c.Facebook(ctx, "test_token")

	// Assertions
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if payload == nil {
		t.Fatal("Expected payload, got nil")
	}
	if payload.Id != "123456789" {
		t.Errorf("Expected ID '123456789', got '%s'", payload.Id)
	}
	if payload.Name != "John Doe" {
		t.Errorf("Expected name 'John Doe', got '%s'", payload.Name)
	}
	if payload.Email != "john.doe@example.com" {
		t.Errorf("Expected email 'john.doe@example.com', got '%s'", payload.Email)
	}
}

func TestClient_Facebook_Error(t *testing.T) {
	// Create a mock server that returns an error response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Return mock Facebook error response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{
			"error": {
				"message": "Invalid OAuth access token.",
				"type": "OAuthException",
				"code": 190
			}
		}`))
	}))
	defer server.Close()

	// Create client with mock server URL
	c := &client{
		facebookURL: server.URL,
		httpClient:  &http.Client{Timeout: 10 * time.Second},
	}

	// Test the Facebook function
	ctx := context.Background()
	payload, err := c.Facebook(ctx, "invalid_token")

	// Assertions
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
	if payload != nil {
		t.Fatalf("Expected nil payload, got %+v", payload)
	}

	// Check if error is of type FacebookError
	facebookErr, ok := err.(*FacebookError)
	if !ok {
		t.Fatalf("Expected FacebookError, got %T", err)
	}
	if facebookErr.Err.Message != "Invalid OAuth access token." {
		t.Errorf("Expected error message 'Invalid OAuth access token.', got '%s'", facebookErr.Err.Message)
	}
	if facebookErr.Err.Type != "OAuthException" {
		t.Errorf("Expected error type 'OAuthException', got '%s'", facebookErr.Err.Type)
	}
	if facebookErr.Err.Code != 190 {
		t.Errorf("Expected error code 190, got %d", facebookErr.Err.Code)
	}
	if facebookErr.Error() != "Invalid OAuth access token." {
		t.Errorf("Expected Error() to return 'Invalid OAuth access token.', got '%s'", facebookErr.Error())
	}
}

func TestClient_Facebook_HTTPError(t *testing.T) {
	// Create client with malformed URL to test HTTP errors
	c := &client{
		facebookURL: "://invalid-url",
		httpClient:  &http.Client{Timeout: 1 * time.Second},
	}

	// Test the Facebook function
	ctx := context.Background()
	payload, err := c.Facebook(ctx, "test_token")

	// Assertions
	if err == nil {
		t.Fatal("Expected error due to invalid URL, got nil")
	}
	if payload != nil {
		t.Fatalf("Expected nil payload, got %+v", payload)
	}
}

func TestClient_Facebook_ContextCancellation(t *testing.T) {
	// Create a mock server that delays response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"id": "123", "name": "Test User"}`))
	}))
	defer server.Close()

	// Create client with mock server URL
	c := &client{
		facebookURL: server.URL,
		httpClient:  &http.Client{Timeout: 10 * time.Second},
	}

	// Create context that cancels immediately
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	// Test the Facebook function
	payload, err := c.Facebook(ctx, "test_token")

	// Assertions
	if err == nil {
		t.Fatal("Expected context cancellation error, got nil")
	}
	if payload != nil {
		t.Fatalf("Expected nil payload, got %+v", payload)
	}
}
