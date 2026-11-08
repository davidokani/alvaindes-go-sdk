package alvaindes

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestClientRequest(t *testing.T) {
	client := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want GET", r.Method)
		}

		if r.URL.Path != "/v1/test" {
			t.Errorf("path = %q, want /v1/test", r.URL.Path)
		}

		if got := r.Header.Get("x-api-key"); got != "test-api-key" {
			t.Errorf("x-api-key = %q, want test-api-key", got)
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{
			"id": "test-123",
		})
	}))

	var result struct {
		ID string `json:"id"`
	}

	err := client.Request(
		context.Background(),
		http.MethodGet,
		"/v1/test",
		nil,
		&result,
	)
	if err != nil {
		t.Fatalf("Client.Request() error = %v", err)
	}

	if result.ID != "test-123" {
		t.Errorf("result.ID = %q, want test-123", result.ID)
	}
}

func TestClientRequestWithOptions(t *testing.T) {
	client := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("X-Custom-Header"); got != "custom-value" {
			t.Errorf("X-Custom-Header = %q, want custom-value", got)
		}

		if got := r.Header.Get("x-api-key"); got != "test-api-key" {
			t.Errorf("x-api-key = %q, want test-api-key", got)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true}`))
	}))

	var result struct {
		Success bool `json:"success"`
	}

	err := client.RequestWithOptions(
		context.Background(),
		http.MethodGet,
		"/v1/test",
		nil,
		&result,
		RequestOptions{
			Headers: http.Header{
				"X-Custom-Header": []string{"custom-value"},
			},
		},
	)
	if err != nil {
		t.Fatalf("Client.RequestWithOptions() error = %v", err)
	}

	if !result.Success {
		t.Error("expected success to be true")
	}
}

func TestAlvaindesErrorMessage(t *testing.T) {
	err := &AlvaindesError{
		Status:  400,
		Message: "Invalid request",
		Code:    "INVALID_REQUEST",
	}

	if got := err.Error(); got != "Invalid request" {
		t.Errorf("Error() = %q, want %q", got, "Invalid request")
	}
}

func TestAlvaindesErrorNilReceiver(t *testing.T) {
	var err *AlvaindesError

	if got := err.Error(); got != "" {
		t.Errorf("nil Error() = %q, want empty string", got)
	}
}
