package alvaindes

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newTestClient(t *testing.T, handler http.Handler) *Client {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	client, err := New(Config{
		APIKey:  "test-api-key",
		BaseURL: server.URL,
	})
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	return client
}

func TestUsersCreate(t *testing.T) {
	client := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %q, want POST", r.Method)
		}

		if r.URL.EscapedPath() != "/v1/applications/app%2F123/users" {
			t.Errorf("path = %q, want escaped application ID", r.URL.EscapedPath())
		}

		if got := r.Header.Get("x-api-key"); got != "test-api-key" {
			t.Errorf("x-api-key = %q, want test-api-key", got)
		}

		if got := r.Header.Get("Content-Type"); got != "application/json" {
			t.Errorf("Content-Type = %q, want application/json", got)
		}

		var body CreateUserParams
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode request body: %v", err)
		}

		if body.FirstName != "Ada" || body.Email != "ada@example.com" {
			t.Errorf("unexpected request body: %+v", body)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"id": "user-123",
			"externalId": null,
			"firstName": "Ada",
			"lastName": "Lovelace",
			"email": "ada@example.com",
			"phoneNumber": "+2348000000000",
			"status": "ACTIVE",
			"createdAt": "2026-11-08T00:00:00.000Z"
		}`))
	}))

	result, err := client.Users.Create(context.Background(), "app/123", CreateUserParams{
		FirstName:   "Ada",
		LastName:    "Lovelace",
		Email:       "ada@example.com",
		PhoneNumber: "+2348000000000",
		Country:     "NG",
	})
	if err != nil {
		t.Fatalf("Users.Create() error = %v", err)
	}

	if result.ID != "user-123" {
		t.Errorf("result.ID = %q, want user-123", result.ID)
	}

	if result.ExternalID != nil {
		t.Errorf("result.ExternalID = %v, want nil", result.ExternalID)
	}
}

func TestUsersGet(t *testing.T) {
	client := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want GET", r.Method)
		}

		if got := r.URL.EscapedPath(); got != "/v1/applications/app-123/users/user%2F456" {
			t.Errorf("path = %q, want escaped user ID", got)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"id": "user-456",
			"externalId": "external-1",
			"firstName": "Ada",
			"lastName": "Lovelace",
			"email": "ada@example.com",
			"phoneNumber": "+2348000000000",
			"status": "ACTIVE",
			"createdAt": "2026-11-08T00:00:00.000Z"
		}`))
	}))

	result, err := client.Users.Get(context.Background(), "app-123", "user/456")
	if err != nil {
		t.Fatalf("Users.Get() error = %v", err)
	}

	if result.ID != "user-456" {
		t.Errorf("result.ID = %q, want user-456", result.ID)
	}

	if result.ExternalID == nil || *result.ExternalID != "external-1" {
		t.Errorf("unexpected external ID: %v", result.ExternalID)
	}
}

func TestAccountsCreate(t *testing.T) {
	client := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %q, want POST", r.Method)
		}

		if got := r.URL.EscapedPath(); got != "/v1/users/user-123/accounts" {
			t.Errorf("path = %q, want /v1/users/user-123/accounts", got)
		}

		var body CreateInvestmentAccountParams
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode request body: %v", err)
		}

		if body.Currency != "USD" || body.BrokerConnectionID != "connection-123" {
			t.Errorf("unexpected request body: %+v", body)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"id": "account-123",
			"userId": "user-123",
			"accountNumber": "123456",
			"status": "ACTIVE",
			"currency": "USD",
			"createdAt": "2026-11-08T00:00:00.000Z"
		}`))
	}))

	result, err := client.Accounts.Create(context.Background(), "user-123", CreateInvestmentAccountParams{
		Currency:           "USD",
		BrokerConnectionID: "connection-123",
	})
	if err != nil {
		t.Fatalf("Accounts.Create() error = %v", err)
	}

	if result.ID != "account-123" {
		t.Errorf("result.ID = %q, want account-123", result.ID)
	}
}

func TestAccountsList(t *testing.T) {
	client := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want GET", r.Method)
		}

		if got := r.URL.EscapedPath(); got != "/v1/users/user-123/accounts" {
			t.Errorf("path = %q, want /v1/users/user-123/accounts", got)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[
			{
				"id": "account-123",
				"userId": "user-123",
				"accountNumber": "123456",
				"status": "ACTIVE",
				"currency": "USD",
				"createdAt": "2026-11-08T00:00:00.000Z"
			}
		]`))
	}))

	result, err := client.Accounts.List(context.Background(), "user-123")
	if err != nil {
		t.Fatalf("Accounts.List() error = %v", err)
	}

	if len(result) != 1 || result[0].ID != "account-123" {
		t.Errorf("unexpected account list: %+v", result)
	}
}

func TestAccountsGet(t *testing.T) {
	client := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want GET", r.Method)
		}

		if got := r.URL.EscapedPath(); got != "/v1/users/user-123/accounts/account%2F456" {
			t.Errorf("path = %q, want escaped account ID", got)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"id": "account-456",
			"userId": "user-123",
			"accountNumber": "654321",
			"status": "ACTIVE",
			"currency": "USD",
			"createdAt": "2026-11-08T00:00:00.000Z"
		}`))
	}))

	result, err := client.Accounts.Get(context.Background(), "user-123", "account/456")
	if err != nil {
		t.Fatalf("Accounts.Get() error = %v", err)
	}

	if result.ID != "account-456" {
		t.Errorf("result.ID = %q, want account-456", result.ID)
	}
}
