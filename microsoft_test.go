package oauth

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestClient_Microsoft_Success(t *testing.T) {
	// Create a mock server that returns a successful Microsoft response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify the request
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		authHeader := r.Header.Get("Authorization")
		if authHeader != "Bearer test_token" {
			t.Errorf("Expected Authorization header 'Bearer test_token', got '%s'", authHeader)
		}

		// Return mock Microsoft user data
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"id": "123456789",
			"mail": "john.doe@example.com",
			"displayName": "John Doe",
			"givenName": "John",
			"surname": "Doe",
			"preferredLanguage": "en-US"
		}`))
	}))
	defer server.Close()

	// Create client with mock server URL
	c := &client{
		microsoftURL: server.URL,
		httpClient:   &http.Client{Timeout: 10 * time.Second},
	}

	// Test the Microsoft function
	ctx := context.Background()
	payload, err := c.Microsoft(ctx, "test_token")

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
	if payload.Email != "john.doe@example.com" {
		t.Errorf("Expected email 'john.doe@example.com', got '%s'", payload.Email)
	}
	if payload.DisplayName != "John Doe" {
		t.Errorf("Expected display name 'John Doe', got '%s'", payload.DisplayName)
	}
	if payload.GivenName != "John" {
		t.Errorf("Expected given name 'John', got '%s'", payload.GivenName)
	}
	if payload.Surname != "Doe" {
		t.Errorf("Expected surname 'Doe', got '%s'", payload.Surname)
	}
	if payload.PreferredLanguage != "en-US" {
		t.Errorf("Expected preferred language 'en-US', got '%s'", payload.PreferredLanguage)
	}
}

func TestClient_Microsoft_Error(t *testing.T) {
	// Create a mock server that returns an error response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Return mock Microsoft error response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{
			"error": {
				"code": "InvalidAuthenticationToken",
				"message": "Access token is invalid"
			}
		}`))
	}))
	defer server.Close()

	// Create client with mock server URL
	c := &client{
		microsoftURL: server.URL,
		httpClient:   &http.Client{Timeout: 10 * time.Second},
	}

	// Test the Microsoft function
	ctx := context.Background()
	payload, err := c.Microsoft(ctx, "invalid_token")

	// Assertions
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
	if payload != nil {
		t.Fatalf("Expected nil payload, got %+v", payload)
	}

	// Check if error is of type MicrosoftError
	microsoftErr, ok := err.(*MicrosoftError)
	if !ok {
		t.Fatalf("Expected MicrosoftError, got %T", err)
	}
	if microsoftErr.Err.Code != "InvalidAuthenticationToken" {
		t.Errorf("Expected error code 'InvalidAuthenticationToken', got '%s'", microsoftErr.Err.Code)
	}
	if microsoftErr.Err.Message != "Access token is invalid" {
		t.Errorf("Expected error message 'Access token is invalid', got '%s'", microsoftErr.Err.Message)
	}
	if microsoftErr.Error() != "Microsoft error" {
		t.Errorf("Expected Error() to return 'Microsoft error', got '%s'", microsoftErr.Error())
	}
}

func TestClient_Microsoft_HTTPError(t *testing.T) {
	// Create client with malformed URL to test HTTP errors
	c := &client{
		microsoftURL: "://invalid-url",
		httpClient:   &http.Client{Timeout: 1 * time.Second},
	}

	// Test the Microsoft function
	ctx := context.Background()
	payload, err := c.Microsoft(ctx, "test_token")

	// Assertions
	if err == nil {
		t.Fatal("Expected error due to invalid URL, got nil")
	}
	if payload != nil {
		t.Fatalf("Expected nil payload, got %+v", payload)
	}
}

func TestClient_Microsoft_ContextCancellation(t *testing.T) {
	// Create a mock server that delays response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"id": "123"}`))
	}))
	defer server.Close()

	// Create client with mock server URL
	c := &client{
		microsoftURL: server.URL,
		httpClient:   &http.Client{Timeout: 10 * time.Second},
	}

	// Create context that cancels immediately
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	// Test the Microsoft function
	payload, err := c.Microsoft(ctx, "test_token")

	// Assertions
	if err == nil {
		t.Fatal("Expected context cancellation error, got nil")
	}
	if payload != nil {
		t.Fatalf("Expected nil payload, got %+v", payload)
	}
}
