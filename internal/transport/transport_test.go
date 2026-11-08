package transport

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func newTestTransport(server *httptest.Server) *Transport {
	return New(Config{
		APIKey:  "test_api_key",
		BaseURL: server.URL,
		Timeout: 2 * time.Second,
	})
}

func TestTransportSendsAuthenticationAndJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}

		if r.URL.Path != "/v1/accounts" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		if r.Header.Get("x-api-key") != "test_api_key" {
			t.Error("missing or incorrect API key")
		}

		if r.Header.Get("Accept") != "application/json" {
			t.Error("incorrect Accept header")
		}

		if r.Header.Get("Content-Type") != "application/json" {
			t.Error("incorrect Content-Type header")
		}

		var body map[string]any

		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Errorf("unable to decode request body: %v", err)
		}

		if body["currency"] != "NGN" {
			t.Errorf("unexpected currency: %v", body["currency"])
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		_, _ = w.Write([]byte(`{"id":"account_123","currency":"NGN"}`))
	}))
	defer server.Close()

	transport := newTestTransport(server)

	var result map[string]any

	err := transport.Request(
		context.Background(),
		http.MethodPost,
		"/v1/accounts",
		map[string]any{"currency": "NGN"},
		&result,
	)

	if err != nil {
		t.Fatalf("unexpected request error: %v", err)
	}

	if result["id"] != "account_123" {
		t.Errorf("unexpected response: %v", result)
	}
}

func TestTransportHandlesPlainTextResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("Alvaindes API is operational"))
	}))
	defer server.Close()

	transport := newTestTransport(server)

	var result string

	err := transport.Request(
		context.Background(),
		http.MethodGet,
		"/health",
		nil,
		&result,
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "Alvaindes API is operational" {
		t.Errorf("unexpected response: %s", result)
	}
}

func TestTransportPreservesAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)

		_, _ = w.Write([]byte(
			`{"statusCode":400,"message":["Invalid account","Currency required"],"code":"VALIDATION_ERROR"}`,
		))
	}))
	defer server.Close()

	transport := newTestTransport(server)

	err := transport.Request(
		context.Background(),
		http.MethodGet,
		"/v1/accounts",
		nil,
		nil,
	)

	var apiError *AlvaindesError

	if !errors.As(err, &apiError) {
		t.Fatalf("expected AlvaindesError, got %T", err)
	}

	if apiError.Status != http.StatusBadRequest {
		t.Errorf("unexpected status: %d", apiError.Status)
	}

	if apiError.Code != "VALIDATION_ERROR" {
		t.Errorf("unexpected code: %s", apiError.Code)
	}

	expectedMessage := "Invalid account; Currency required"

	if apiError.Message != expectedMessage {
		t.Errorf("expected %q, got %q", expectedMessage, apiError.Message)
	}

	if apiError.Response == nil {
		t.Error("expected original response body to be preserved")
	}
}

func TestTransportRejectsInvalidPaths(t *testing.T) {
	transport := New(Config{
		APIKey:  "test_api_key",
		BaseURL: "https://api.alvaindes.com",
		Timeout: time.Second,
	})

	invalidPaths := []string{
		"v1/accounts",
		"//example.com/accounts",
	}

	for _, path := range invalidPaths {
		err := transport.Request(
			context.Background(),
			http.MethodGet,
			path,
			nil,
			nil,
		)

		if err == nil {
			t.Errorf("expected error for path %q", path)
		}
	}
}

func TestTransportRejectsUnsupportedMethods(t *testing.T) {
	transport := New(Config{
		APIKey:  "test_api_key",
		BaseURL: "https://api.alvaindes.com",
		Timeout: time.Second,
	})

	err := transport.Request(
		context.Background(),
		"OPTIONS",
		"/v1/accounts",
		nil,
		nil,
	)

	if err == nil {
		t.Fatal("expected unsupported method error")
	}
}

func TestTransportHandlesCallerCancellation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		<-r.Context().Done()
	}))
	defer server.Close()

	transport := newTestTransport(server)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := transport.Request(
		ctx,
		http.MethodGet,
		"/v1/accounts",
		nil,
		nil,
	)

	var apiError *AlvaindesError

	if !errors.As(err, &apiError) {
		t.Fatalf("expected AlvaindesError, got %T", err)
	}

	if apiError.Code != "REQUEST_ABORTED" {
		t.Errorf("expected REQUEST_ABORTED, got %s", apiError.Code)
	}
}

func TestTransportHandlesTimeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		<-r.Context().Done()
	}))
	defer server.Close()

	transport := New(Config{
		APIKey:  "test_api_key",
		BaseURL: server.URL,
		Timeout: 50 * time.Millisecond,
	})

	err := transport.Request(
		context.Background(),
		http.MethodGet,
		"/v1/accounts",
		nil,
		nil,
	)

	var apiError *AlvaindesError

	if !errors.As(err, &apiError) {
		t.Fatalf("expected AlvaindesError, got %T", err)
	}

	if apiError.Code != "TIMEOUT" {
		t.Errorf("expected TIMEOUT, got %s", apiError.Code)
	}
}

func TestTransportSupportsCustomHeadersWithoutOverridingManagedHeaders(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Idempotency-Key") != "idem_123" {
			t.Errorf("expected custom Idempotency-Key header")
		}

		if r.Header.Get("x-api-key") != "test_api_key" {
			t.Errorf("custom headers must not override the SDK API key")
		}

		if r.Header.Get("Accept") != "application/json" {
			t.Errorf("custom headers must not override Accept")
		}

		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("custom headers must not override Content-Type")
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true}`))
	}))
	defer server.Close()

	transport := newTestTransport(server)

	headers := make(http.Header)
	headers.Set("Idempotency-Key", "idem_123")
	headers.Set("x-api-key", "incorrect_key")
	headers.Set("Accept", "text/plain")
	headers.Set("Content-Type", "text/plain")

	var result map[string]any

	err := transport.RequestWithOptions(
		context.Background(),
		http.MethodPost,
		"/v1/orders",
		map[string]any{"symbol": "AAPL"},
		&result,
		RequestOptions{Headers: headers},
	)
	if err != nil {
		t.Fatalf("unexpected request error: %v", err)
	}

	if result["success"] != true {
		t.Errorf("unexpected response: %v", result)
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(request *http.Request) (*http.Response, error) {
	return f(request)
}

func TestTransportUsesInjectedHTTPClient(t *testing.T) {
	called := false

	httpClient := &http.Client{
		Transport: roundTripFunc(func(request *http.Request) (*http.Response, error) {
			called = true

			if request.URL.String() != "https://api.example.com/v1/accounts" {
				t.Errorf("unexpected request URL: %s", request.URL)
			}

			return &http.Response{
				StatusCode: http.StatusOK,
				Header:     make(http.Header),
				Body:       io.NopCloser(strings.NewReader(`{"id":"account_custom"}`)),
				Request:    request,
			}, nil
		}),
	}

	transport := New(Config{
		APIKey:     "test_api_key",
		BaseURL:    "https://api.example.com",
		Timeout:    2 * time.Second,
		HTTPClient: httpClient,
	})

	var result map[string]any

	err := transport.Request(
		context.Background(),
		http.MethodGet,
		"/v1/accounts",
		nil,
		&result,
	)
	if err != nil {
		t.Fatalf("unexpected request error: %v", err)
	}

	if !called {
		t.Fatal("expected the injected HTTP client to handle the request")
	}

	if result["id"] != "account_custom" {
		t.Errorf("unexpected response: %v", result)
	}
}
