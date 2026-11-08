package alvaindes

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestDepositsCreate(t *testing.T) {
	fundingSourceID := "funding-source-123"

	client := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %q, want POST", r.Method)
		}

		if got := r.URL.EscapedPath(); got != "/v1/users/user%2F123/accounts/account%2F456/deposits" {
			t.Errorf("path = %q, want escaped deposits path", got)
		}

		if got := r.Header.Get("Idempotency-Key"); got != "deposit-request-123" {
			t.Errorf("Idempotency-Key = %q, want deposit-request-123", got)
		}

		if got := r.Header.Get("x-api-key"); got != "test-api-key" {
			t.Errorf("x-api-key = %q, want test-api-key", got)
		}

		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Errorf("decode request body: %v", err)
		} else {
			if body["amount"] != float64(250.50) {
				t.Errorf("body amount = %v, want 250.50", body["amount"])
			}

			if body["currency"] != "USD" {
				t.Errorf("body currency = %v, want USD", body["currency"])
			}

			if body["fundingSourceId"] != fundingSourceID {
				t.Errorf("body fundingSourceId = %v, want %q", body["fundingSourceId"], fundingSourceID)
			}
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"id": "deposit-123",
			"amount": 250.50,
			"currency": "USD",
			"status": "PENDING",
			"reference": "deposit-ref-123",
			"createdAt": "2026-11-08T00:00:00.000Z"
		}`))
	}))

	result, err := client.Deposits.Create(
		context.Background(),
		"user/123",
		"account/456",
		CreateDepositParams{
			Amount:          250.50,
			Currency:        "USD",
			FundingSourceID: &fundingSourceID,
		},
		CreateDepositOptions{
			IdempotencyKey: "deposit-request-123",
		},
	)
	if err != nil {
		t.Fatalf("Deposits.Create() error = %v", err)
	}

	if result.ID != "deposit-123" {
		t.Errorf("result.ID = %q, want deposit-123", result.ID)
	}

	if result.Amount != 250.50 {
		t.Errorf("result.Amount = %v, want 250.50", result.Amount)
	}

	if result.Currency != "USD" {
		t.Errorf("result.Currency = %q, want USD", result.Currency)
	}

	if result.Status != "PENDING" {
		t.Errorf("result.Status = %q, want PENDING", result.Status)
	}

	if result.Reference != "deposit-ref-123" {
		t.Errorf("result.Reference = %q, want deposit-ref-123", result.Reference)
	}
}

func TestDepositsCreateWithoutOptionalFundingSourceOrIdempotencyKey(t *testing.T) {
	client := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Idempotency-Key"); got != "" {
			t.Errorf("Idempotency-Key = %q, want empty", got)
		}

		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Errorf("decode request body: %v", err)
		} else {
			if _, exists := body["fundingSourceId"]; exists {
				t.Errorf("fundingSourceId should be omitted when nil")
			}

			if body["amount"] != float64(100) {
				t.Errorf("body amount = %v, want 100", body["amount"])
			}

			if body["currency"] != "NGN" {
				t.Errorf("body currency = %v, want NGN", body["currency"])
			}
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"id": "deposit-456",
			"amount": 100,
			"currency": "NGN",
			"status": "PENDING",
			"reference": "deposit-ref-456",
			"createdAt": "2026-11-08T00:00:00.000Z"
		}`))
	}))

	result, err := client.Deposits.Create(
		context.Background(),
		"user-123",
		"account-456",
		CreateDepositParams{
			Amount:   100,
			Currency: "NGN",
		},
		CreateDepositOptions{},
	)
	if err != nil {
		t.Fatalf("Deposits.Create() error = %v", err)
	}

	if result.ID != "deposit-456" {
		t.Errorf("result.ID = %q, want deposit-456", result.ID)
	}
}
