package alvaindes

import (
	"context"
	"net/http"
	"testing"
)

func TestTransactionsList(t *testing.T) {
	client := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want GET", r.Method)
		}

		if got := r.URL.EscapedPath(); got != "/v1/users/user%2F123/accounts/account%2F456/transactions" {
			t.Errorf("path = %q, want escaped transactions path", got)
		}

		if got := r.Header.Get("x-api-key"); got != "test-api-key" {
			t.Errorf("x-api-key = %q, want test-api-key", got)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[
			{
				"id": "transaction-123",
				"type": "DEPOSIT",
				"amount": 500,
				"currency": "USD",
				"description": "Initial deposit",
				"status": "COMPLETED",
				"createdAt": "2026-11-08T00:00:00.000Z"
			},
			{
				"id": "transaction-456",
				"type": "WITHDRAWAL",
				"amount": 100,
				"currency": "NGN",
				"status": "PENDING",
				"createdAt": "2026-11-08T01:00:00.000Z"
			}
		]`))
	}))

	result, err := client.Transactions.List(
		context.Background(),
		"user/123",
		"account/456",
	)
	if err != nil {
		t.Fatalf("Transactions.List() error = %v", err)
	}

	if len(result) != 2 {
		t.Fatalf("len(result) = %d, want 2", len(result))
	}

	first := result[0]

	if first.ID != "transaction-123" {
		t.Errorf("result[0].ID = %q, want transaction-123", first.ID)
	}

	if first.Type != "DEPOSIT" {
		t.Errorf("result[0].Type = %q, want DEPOSIT", first.Type)
	}

	if first.Amount != 500 {
		t.Errorf("result[0].Amount = %v, want 500", first.Amount)
	}

	if first.Currency != "USD" {
		t.Errorf("result[0].Currency = %q, want USD", first.Currency)
	}

	if first.Description == nil || *first.Description != "Initial deposit" {
		t.Errorf("result[0].Description = %v, want Initial deposit", first.Description)
	}

	if first.Status != "COMPLETED" {
		t.Errorf("result[0].Status = %q, want COMPLETED", first.Status)
	}

	second := result[1]

	if second.ID != "transaction-456" {
		t.Errorf("result[1].ID = %q, want transaction-456", second.ID)
	}

	if second.Type != "WITHDRAWAL" {
		t.Errorf("result[1].Type = %q, want WITHDRAWAL", second.Type)
	}

	if second.Description != nil {
		t.Errorf("result[1].Description = %v, want nil when omitted", *second.Description)
	}
}
