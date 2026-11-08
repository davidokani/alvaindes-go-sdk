package alvaindes

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestFundingSourcesCreate(t *testing.T) {
	client := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %q, want POST", r.Method)
		}

		if got := r.URL.EscapedPath(); got != "/v1/users/user%2F123/accounts/account%2F456/funding-sources" {
			t.Errorf("path = %q, want escaped resource path", got)
		}

		var body CreateFundingSourceParams
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode request body: %v", err)
		}

		if body.Type != FundingSourceACH {
			t.Errorf("Type = %q, want ACH", body.Type)
		}

		if body.OwnerName != "Ada Lovelace" {
			t.Errorf("OwnerName = %q, want Ada Lovelace", body.OwnerName)
		}

		if body.Currency != "USD" {
			t.Errorf("Currency = %q, want USD", body.Currency)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"id": "funding-123",
			"investmentAccountId": "account/456",
			"type": "ACH",
			"status": "PENDING",
			"currency": "USD",
			"brokerReference": null,
			"createdAt": "2026-11-08T00:00:00.000Z",
			"updatedAt": "2026-11-08T00:00:00.000Z"
		}`))
	}))

	result, err := client.FundingSources.Create(
		context.Background(),
		"user/123",
		"account/456",
		CreateFundingSourceParams{
			Type:              FundingSourceACH,
			OwnerName:         "Ada Lovelace",
			BankAccountNumber: "1234567890",
			BankAccountType:   "CHECKING",
			RoutingNumber:     "021000021",
			Currency:          "USD",
		},
	)
	if err != nil {
		t.Fatalf("FundingSources.Create() error = %v", err)
	}

	if result.ID != "funding-123" {
		t.Errorf("result.ID = %q, want funding-123", result.ID)
	}

	if result.Status != FundingSourcePending {
		t.Errorf("result.Status = %q, want PENDING", result.Status)
	}

	if result.BrokerReference != nil {
		t.Errorf("result.BrokerReference = %v, want nil", result.BrokerReference)
	}
}

func TestFundingSourcesCreateRejectsNonACH(t *testing.T) {
	requestMade := false

	client := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestMade = true
		t.Error("unexpected HTTP request for unsupported funding source type")
	}))

	_, err := client.FundingSources.Create(
		context.Background(),
		"user-123",
		"account-456",
		CreateFundingSourceParams{
			Type: FundingSourceWire,
		},
	)
	if err == nil {
		t.Fatal("FundingSources.Create() error = nil, want ACH-only validation error")
	}

	if requestMade {
		t.Error("HTTP request was made for unsupported funding source type")
	}
}

func TestFundingSourcesList(t *testing.T) {
	client := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want GET", r.Method)
		}

		if got := r.URL.EscapedPath(); got != "/v1/users/user%2F123/accounts/account%2F456/funding-sources" {
			t.Errorf("path = %q, want escaped resource path", got)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[
			{
				"id": "funding-123",
				"investmentAccountId": "account/456",
				"type": "ACH",
				"status": "ACTIVE",
				"currency": "USD",
				"brokerReference": "broker-ref-123",
				"createdAt": "2026-11-01T00:00:00.000Z",
				"updatedAt": "2026-11-08T00:00:00.000Z"
			},
			{
				"id": "funding-456",
				"investmentAccountId": "account/456",
				"type": "BANK_ACCOUNT",
				"status": "DISABLED",
				"currency": "NGN",
				"brokerReference": null,
				"createdAt": "2026-11-02T00:00:00.000Z",
				"updatedAt": "2026-11-08T00:00:00.000Z"
			}
		]`))
	}))

	result, err := client.FundingSources.List(
		context.Background(),
		"user/123",
		"account/456",
	)
	if err != nil {
		t.Fatalf("FundingSources.List() error = %v", err)
	}

	if len(result) != 2 {
		t.Fatalf("len(result) = %d, want 2", len(result))
	}

	if result[0].Type != FundingSourceACH {
		t.Errorf("result[0].Type = %q, want ACH", result[0].Type)
	}

	if result[1].Type != FundingSourceBankAccount {
		t.Errorf("result[1].Type = %q, want BANK_ACCOUNT", result[1].Type)
	}

	if result[1].BrokerReference != nil {
		t.Errorf("result[1].BrokerReference = %v, want nil", result[1].BrokerReference)
	}
}

func TestFundingSourcesRefresh(t *testing.T) {
	client := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %q, want POST", r.Method)
		}

		if got := r.URL.EscapedPath(); got != "/v1/users/user%2F123/accounts/account%2F456/funding-sources/funding%2F789/refresh" {
			t.Errorf("path = %q, want escaped refresh path", got)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"id": "funding/789",
			"investmentAccountId": "account/456",
			"type": "ACH",
			"status": "ACTIVE",
			"currency": "USD",
			"brokerReference": "broker-ref-789",
			"createdAt": "2026-11-01T00:00:00.000Z",
			"updatedAt": "2026-11-08T00:00:00.000Z"
		}`))
	}))

	result, err := client.FundingSources.Refresh(
		context.Background(),
		"user/123",
		"account/456",
		"funding/789",
	)
	if err != nil {
		t.Fatalf("FundingSources.Refresh() error = %v", err)
	}

	if result.ID != "funding/789" {
		t.Errorf("result.ID = %q, want funding/789", result.ID)
	}

	if result.Status != FundingSourceActive {
		t.Errorf("result.Status = %q, want ACTIVE", result.Status)
	}
}

func TestFundingSourcesDisable(t *testing.T) {
	client := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("method = %q, want DELETE", r.Method)
		}

		if got := r.URL.EscapedPath(); got != "/v1/users/user%2F123/accounts/account%2F456/funding-sources/funding%2F789" {
			t.Errorf("path = %q, want escaped funding source path", got)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"id": "funding/789",
			"investmentAccountId": "account/456",
			"type": "ACH",
			"status": "DISABLED",
			"currency": "USD",
			"brokerReference": null,
			"createdAt": "2026-11-01T00:00:00.000Z",
			"updatedAt": "2026-11-08T00:00:00.000Z"
		}`))
	}))

	result, err := client.FundingSources.Disable(
		context.Background(),
		"user/123",
		"account/456",
		"funding/789",
	)
	if err != nil {
		t.Fatalf("FundingSources.Disable() error = %v", err)
	}

	if result.Status != FundingSourceDisabled {
		t.Errorf("result.Status = %q, want DISABLED", result.Status)
	}
}
