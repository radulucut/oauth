package oauth

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestClient_Google_Success(t *testing.T) {
	// Create a mock server that returns a successful Google response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify the request
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		if r.URL.Query().Get("access_token") != "test_token" {
			t.Errorf("Expected access_token=test_token, got %s", r.URL.Query().Get("access_token"))
		}

		// Return mock Google user data
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"sub": "123456789",
			"name": "John Doe",
			"given_name": "John",
			"family_name": "Doe",
			"picture": "https://example.com/photo.jpg",
			"email": "john.doe@example.com",
			"email_verified": true,
			"locale": "en"
		}`))
	}))
	defer server.Close()

	// Create client with mock server URL
	c := &client{
		googleURL:  server.URL,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}

	// Test the Google function
	ctx := context.Background()
	payload, err := c.Google(ctx, "test_token")

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
	if payload.GivenName != "John" {
		t.Errorf("Expected given name 'John', got '%s'", payload.GivenName)
	}
	if payload.FamilyName != "Doe" {
		t.Errorf("Expected family name 'Doe', got '%s'", payload.FamilyName)
	}
	if payload.PictureURL != "https://example.com/photo.jpg" {
		t.Errorf("Expected picture URL 'https://example.com/photo.jpg', got '%s'", payload.PictureURL)
	}
	if payload.Email != "john.doe@example.com" {
		t.Errorf("Expected email 'john.doe@example.com', got '%s'", payload.Email)
	}
	if !payload.EmailVerified {
		t.Error("Expected email verified to be true")
	}
	if payload.Locale != "en" {
		t.Errorf("Expected locale 'en', got '%s'", payload.Locale)
	}
}

func TestClient_Google_Error(t *testing.T) {
	// Create a mock server that returns an error response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Return mock Google error response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{
			"error": "invalid_token",
			"error_description": "Invalid OAuth access token."
		}`))
	}))
	defer server.Close()

	// Create client with mock server URL
	c := &client{
		googleURL:  server.URL,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}

	// Test the Google function
	ctx := context.Background()
	payload, err := c.Google(ctx, "invalid_token")

	// Assertions
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
	if payload != nil {
		t.Fatalf("Expected nil payload, got %+v", payload)
	}

	// Check if error is of type GoogleError
	googleErr, ok := err.(*GoogleError)
	if !ok {
		t.Fatalf("Expected GoogleError, got %T", err)
	}
	if googleErr.Name != "invalid_token" {
		t.Errorf("Expected error name 'invalid_token', got '%s'", googleErr.Name)
	}
	if googleErr.Description != "Invalid OAuth access token." {
		t.Errorf("Expected error description 'Invalid OAuth access token.', got '%s'", googleErr.Description)
	}
	if googleErr.Error() != "Invalid OAuth access token." {
		t.Errorf("Expected Error() to return 'Invalid OAuth access token.', got '%s'", googleErr.Error())
	}
}

func TestClient_Google_HTTPError(t *testing.T) {
	// Create client with malformed URL to test HTTP errors
	c := &client{
		googleURL:  "://invalid-url",
		httpClient: &http.Client{Timeout: 1 * time.Second},
	}

	// Test the Google function
	ctx := context.Background()
	payload, err := c.Google(ctx, "test_token")

	// Assertions
	if err == nil {
		t.Fatal("Expected error due to invalid URL, got nil")
	}
	if payload != nil {
		t.Fatalf("Expected nil payload, got %+v", payload)
	}
}

func TestClient_Google_ContextCancellation(t *testing.T) {
	// Create a mock server that delays response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"sub": "123"}`))
	}))
	defer server.Close()

	// Create client with mock server URL
	c := &client{
		googleURL:  server.URL,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}

	// Create context that cancels immediately
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	// Test the Google function
	payload, err := c.Google(ctx, "test_token")

	// Assertions
	if err == nil {
		t.Fatal("Expected context cancellation error, got nil")
	}
	if payload != nil {
		t.Fatalf("Expected nil payload, got %+v", payload)
	}
}
