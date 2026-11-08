package alvaindes

import (
	"context"
	"net/http"
	"testing"
)

func TestMarketDataGetQuote(t *testing.T) {
	client := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want GET", r.Method)
		}

		if got := r.URL.EscapedPath(); got != "/v1/market-data/AAPL" {
			t.Errorf("path = %q, want /v1/market-data/AAPL", got)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"symbol": "AAPL",
			"price": 225.50,
			"currency": "USD",
			"timestamp": "2026-11-08T10:30:00.000Z"
		}`))
	}))

	result, err := client.MarketData.GetQuote(context.Background(), "AAPL")
	if err != nil {
		t.Fatalf("MarketData.GetQuote() error = %v", err)
	}

	if result == nil {
		t.Fatal("MarketData.GetQuote() returned nil")
	}

	if result.Symbol != "AAPL" {
		t.Errorf("result.Symbol = %q, want AAPL", result.Symbol)
	}

	if result.Price != 225.50 {
		t.Errorf("result.Price = %v, want 225.50", result.Price)
	}

	if result.Currency != "USD" {
		t.Errorf("result.Currency = %q, want USD", result.Currency)
	}

	if result.Timestamp != "2026-11-08T10:30:00.000Z" {
		t.Errorf("result.Timestamp = %q, want 2026-11-08T10:30:00.000Z", result.Timestamp)
	}
}

func TestMarketDataGetQuoteEscapesSymbol(t *testing.T) {
	client := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.EscapedPath(); got != "/v1/market-data/BRK%2FB" {
			t.Errorf("path = %q, want /v1/market-data/BRK%%2FB", got)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"symbol": "BRK/B",
			"price": 500.25,
			"currency": "USD",
			"timestamp": "2026-11-08T10:30:00.000Z"
		}`))
	}))

	result, err := client.MarketData.GetQuote(context.Background(), "BRK/B")
	if err != nil {
		t.Fatalf("MarketData.GetQuote() error = %v", err)
	}

	if result.Symbol != "BRK/B" {
		t.Errorf("result.Symbol = %q, want BRK/B", result.Symbol)
	}
}
