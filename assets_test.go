package alvaindes

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestAssetsList(t *testing.T) {
	client := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want GET", r.Method)
		}

		if got := r.URL.Path; got != "/v1/assets" {
			t.Errorf("path = %q, want /v1/assets", got)
		}

		if got := r.URL.Query().Get("connectionId"); got != "conn/123" {
			t.Errorf("connectionId = %q, want conn/123", got)
		}

		if got := r.URL.Query().Get("query"); got != "AAPL & MSFT" {
			t.Errorf("query = %q, want AAPL & MSFT", got)
		}

		if got := r.Header.Get("x-api-key"); got != "test-api-key" {
			t.Errorf("x-api-key = %q, want test-api-key", got)
		}

		if got := r.URL.RawQuery; got != "connectionId=conn%2F123&query=AAPL+%26+MSFT" {
			t.Errorf("RawQuery = %q, want correctly encoded query parameters", got)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[
			{
				"symbol": "AAPL",
				"name": "Apple Inc.",
				"exchange": "NASDAQ",
				"currency": "USD",
				"type": "STOCK",
				"status": "ACTIVE"
			},
			{
				"symbol": "MSFT",
				"name": "Microsoft Corporation",
				"exchange": "NASDAQ",
				"currency": "USD",
				"type": "STOCK",
				"status": "ACTIVE"
			}
		]`))
	}))

	result, err := client.Assets.List(
		context.Background(),
		ListAssetsParams{
			ConnectionID: "conn/123",
			Query:        "AAPL & MSFT",
		},
	)
	if err != nil {
		t.Fatalf("Assets.List() error = %v", err)
	}

	if len(result) != 2 {
		t.Fatalf("len(result) = %d, want 2", len(result))
	}

	if result[0].Symbol != "AAPL" {
		t.Errorf("result[0].Symbol = %q, want AAPL", result[0].Symbol)
	}

	if result[0].Name != "Apple Inc." {
		t.Errorf("result[0].Name = %q, want Apple Inc.", result[0].Name)
	}

	if result[1].Symbol != "MSFT" {
		t.Errorf("result[1].Symbol = %q, want MSFT", result[1].Symbol)
	}

	if result[1].Status != "ACTIVE" {
		t.Errorf("result[1].Status = %q, want ACTIVE", result[1].Status)
	}
}

func TestAssetsListEmptyResponse(t *testing.T) {
	client := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want GET", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[]`))
	}))

	result, err := client.Assets.List(
		context.Background(),
		ListAssetsParams{
			ConnectionID: "connection-123",
			Query:        "UNKNOWN",
		},
	)
	if err != nil {
		t.Fatalf("Assets.List() error = %v", err)
	}

	if len(result) != 0 {
		t.Errorf("len(result) = %d, want 0", len(result))
	}

	// Confirm the response was decoded as an empty JSON array.
	encoded, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("marshal result: %v", err)
	}

	if string(encoded) != "[]" {
		t.Errorf("encoded result = %s, want []", encoded)
	}
}
