package alvaindes

import (
	"context"
	"net/http"
	"testing"
)

func TestHoldingsList(t *testing.T) {
	client := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want GET", r.Method)
		}

		if got := r.URL.EscapedPath(); got != "/v1/users/user%2F123/accounts/account%2F456/holdings" {
			t.Errorf("path = %q, want escaped holdings path", got)
		}

		if got := r.Header.Get("x-api-key"); got != "test-api-key" {
			t.Errorf("x-api-key = %q, want test-api-key", got)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[
			{
				"id": "holding-123",
				"symbol": "AAPL",
				"quantity": 10,
				"averagePrice": 175.50,
				"marketValue": 1900.00
			},
			{
				"id": "holding-456",
				"symbol": "MSFT",
				"quantity": 5,
				"averagePrice": 400.00,
				"marketValue": 2100.00
			}
		]`))
	}))

	result, err := client.Holdings.List(
		context.Background(),
		"user/123",
		"account/456",
	)
	if err != nil {
		t.Fatalf("Holdings.List() error = %v", err)
	}

	if len(result) != 2 {
		t.Fatalf("len(result) = %d, want 2", len(result))
	}

	if result[0].ID != "holding-123" {
		t.Errorf("result[0].ID = %q, want holding-123", result[0].ID)
	}

	if result[0].Symbol != "AAPL" {
		t.Errorf("result[0].Symbol = %q, want AAPL", result[0].Symbol)
	}

	if result[0].Quantity != 10 {
		t.Errorf("result[0].Quantity = %v, want 10", result[0].Quantity)
	}

	if result[0].AveragePrice != 175.50 {
		t.Errorf("result[0].AveragePrice = %v, want 175.50", result[0].AveragePrice)
	}

	if result[0].MarketValue != 1900.00 {
		t.Errorf("result[0].MarketValue = %v, want 1900.00", result[0].MarketValue)
	}

	if result[1].ID != "holding-456" {
		t.Errorf("result[1].ID = %q, want holding-456", result[1].ID)
	}

	if result[1].Symbol != "MSFT" {
		t.Errorf("result[1].Symbol = %q, want MSFT", result[1].Symbol)
	}
}
