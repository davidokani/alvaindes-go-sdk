package alvaindes

import (
	"context"
	"net/http"
	"testing"
)

func TestPortfoliosGet(t *testing.T) {
	client := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want GET", r.Method)
		}

		if got := r.URL.EscapedPath(); got != "/v1/users/user%2F123/accounts/account%2F456/portfolio" {
			t.Errorf("path = %q, want escaped portfolio path", got)
		}

		if got := r.Header.Get("x-api-key"); got != "test-api-key" {
			t.Errorf("x-api-key = %q, want test-api-key", got)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"id": "portfolio-123",
			"investmentAccountId": "account/456",
			"accountValue": 12500.75,
			"cashBalance": 2500.25,
			"currency": "USD",
			"status": "ACTIVE",
			"createdAt": "2026-11-08T00:00:00.000Z"
		}`))
	}))

	result, err := client.Portfolios.Get(
		context.Background(),
		"user/123",
		"account/456",
	)
	if err != nil {
		t.Fatalf("Portfolios.Get() error = %v", err)
	}

	if result.ID != "portfolio-123" {
		t.Errorf("result.ID = %q, want portfolio-123", result.ID)
	}

	if result.InvestmentAccountID != "account/456" {
		t.Errorf(
			"result.InvestmentAccountID = %q, want account/456",
			result.InvestmentAccountID,
		)
	}

	if result.AccountValue != 12500.75 {
		t.Errorf("result.AccountValue = %v, want 12500.75", result.AccountValue)
	}

	if result.CashBalance != 2500.25 {
		t.Errorf("result.CashBalance = %v, want 2500.25", result.CashBalance)
	}

	if result.Currency != "USD" {
		t.Errorf("result.Currency = %q, want USD", result.Currency)
	}

	if result.Status != "ACTIVE" {
		t.Errorf("result.Status = %q, want ACTIVE", result.Status)
	}
}
