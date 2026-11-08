package alvaindes

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestBrokerAgreementsGetRequirements(t *testing.T) {
	client := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want GET", r.Method)
		}

		if got := r.URL.EscapedPath(); got != "/v1/users/user%2F123/broker-connections/connection%2F456/agreements/requirements" {
			t.Errorf("path = %q, want escaped requirements path", got)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[
			{
				"agreement": "account_agreement",
				"revision": "1.0"
			},
			{
				"agreement": "customer_agreement"
			}
		]`))
	}))

	result, err := client.BrokerAgreements.GetRequirements(
		context.Background(),
		"user/123",
		"connection/456",
	)
	if err != nil {
		t.Fatalf("BrokerAgreements.GetRequirements() error = %v", err)
	}

	if len(result) != 2 {
		t.Fatalf("len(result) = %d, want 2", len(result))
	}

	if result[0].Agreement != "account_agreement" {
		t.Errorf("result[0].Agreement = %q, want account_agreement", result[0].Agreement)
	}

	if result[0].Revision == nil || *result[0].Revision != "1.0" {
		t.Errorf("result[0].Revision = %v, want 1.0", result[0].Revision)
	}

	if result[1].Agreement != "customer_agreement" {
		t.Errorf("result[1].Agreement = %q, want customer_agreement", result[1].Agreement)
	}

	if result[1].Revision != nil {
		t.Errorf("result[1].Revision = %v, want nil", *result[1].Revision)
	}
}

func TestBrokerAgreementsListAccepted(t *testing.T) {
	client := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want GET", r.Method)
		}

		if got := r.URL.EscapedPath(); got != "/v1/users/user%2F123/broker-connections/connection%2F456/agreements" {
			t.Errorf("path = %q, want escaped agreements path", got)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[
			{
				"id": "agreement-123",
				"userId": "user/123",
				"brokerConnectionId": "connection/456",
				"agreement": "account_agreement",
				"revision": "1.0",
				"signedAt": "2026-11-08T10:00:00.000Z",
				"ipAddress": "192.0.2.10",
				"createdAt": "2026-11-08T10:00:00.000Z"
			},
			{
				"id": "agreement-456",
				"userId": "user/123",
				"brokerConnectionId": "connection/456",
				"agreement": "customer_agreement",
				"revision": null,
				"signedAt": "2026-11-08T10:05:00.000Z",
				"ipAddress": "192.0.2.10",
				"createdAt": "2026-11-08T10:05:00.000Z"
			}
		]`))
	}))

	result, err := client.BrokerAgreements.ListAccepted(
		context.Background(),
		"user/123",
		"connection/456",
	)
	if err != nil {
		t.Fatalf("BrokerAgreements.ListAccepted() error = %v", err)
	}

	if len(result) != 2 {
		t.Fatalf("len(result) = %d, want 2", len(result))
	}

	if result[0].ID != "agreement-123" {
		t.Errorf("result[0].ID = %q, want agreement-123", result[0].ID)
	}

	if result[0].Agreement != "account_agreement" {
		t.Errorf("result[0].Agreement = %q, want account_agreement", result[0].Agreement)
	}

	if result[0].Revision == nil || *result[0].Revision != "1.0" {
		t.Errorf("result[0].Revision = %v, want 1.0", result[0].Revision)
	}

	if result[1].Revision != nil {
		t.Errorf("result[1].Revision = %v, want nil", *result[1].Revision)
	}
}

func TestBrokerAgreementsAccept(t *testing.T) {
	revision := "2.0"

	client := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %q, want POST", r.Method)
		}

		if got := r.URL.EscapedPath(); got != "/v1/users/user%2F123/broker-connections/connection%2F456/agreements" {
			t.Errorf("path = %q, want escaped agreements path", got)
		}

		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Errorf("decode request body: %v", err)
		} else {
			if body["agreement"] != "account_agreement" {
				t.Errorf("body agreement = %v, want account_agreement", body["agreement"])
			}

			if body["revision"] != revision {
				t.Errorf("body revision = %v, want %q", body["revision"], revision)
			}
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"id": "agreement-789",
			"userId": "user/123",
			"brokerConnectionId": "connection/456",
			"agreement": "account_agreement",
			"revision": "2.0",
			"signedAt": "2026-11-08T10:00:00.000Z",
			"ipAddress": "192.0.2.10",
			"createdAt": "2026-11-08T10:00:00.000Z"
		}`))
	}))

	result, err := client.BrokerAgreements.Accept(
		context.Background(),
		"user/123",
		"connection/456",
		AcceptBrokerAgreementParams{
			Agreement: "account_agreement",
			Revision:  &revision,
		},
	)
	if err != nil {
		t.Fatalf("BrokerAgreements.Accept() error = %v", err)
	}

	if result.ID != "agreement-789" {
		t.Errorf("result.ID = %q, want agreement-789", result.ID)
	}

	if result.UserID != "user/123" {
		t.Errorf("result.UserID = %q, want user/123", result.UserID)
	}

	if result.BrokerConnectionID != "connection/456" {
		t.Errorf("result.BrokerConnectionID = %q, want connection/456", result.BrokerConnectionID)
	}

	if result.Agreement != "account_agreement" {
		t.Errorf("result.Agreement = %q, want account_agreement", result.Agreement)
	}

	if result.Revision == nil || *result.Revision != revision {
		t.Errorf("result.Revision = %v, want %q", result.Revision, revision)
	}
}

func TestBrokerAgreementsAcceptWithoutOptionalRevision(t *testing.T) {
	client := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Errorf("decode request body: %v", err)
		} else {
			if body["agreement"] != "customer_agreement" {
				t.Errorf("body agreement = %v, want customer_agreement", body["agreement"])
			}

			if _, exists := body["revision"]; exists {
				t.Errorf("revision should be omitted when nil")
			}
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"id": "agreement-999",
			"userId": "user-123",
			"brokerConnectionId": "connection-456",
			"agreement": "customer_agreement",
			"revision": null,
			"signedAt": "2026-11-08T10:00:00.000Z",
			"ipAddress": "192.0.2.10",
			"createdAt": "2026-11-08T10:00:00.000Z"
		}`))
	}))

	result, err := client.BrokerAgreements.Accept(
		context.Background(),
		"user-123",
		"connection-456",
		AcceptBrokerAgreementParams{
			Agreement: "customer_agreement",
		},
	)
	if err != nil {
		t.Fatalf("BrokerAgreements.Accept() error = %v", err)
	}

	if result.Revision != nil {
		t.Errorf("result.Revision = %v, want nil", *result.Revision)
	}
}
