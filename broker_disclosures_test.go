package alvaindes

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestBrokerDisclosuresGet(t *testing.T) {
	client := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want GET", r.Method)
		}

		if got := r.URL.EscapedPath(); got != "/v1/users/user%2F123/broker-connections/connection%2F456/disclosures" {
			t.Errorf("path = %q, want escaped disclosures path", got)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"id": "disclosure-123",
			"userId": "user/123",
			"brokerConnectionId": "connection/456",
			"employmentStatus": "employed",
			"employmentSector": "technology",
			"isControlPerson": false,
			"isAffiliatedExchangeOrFinra": false,
			"isPoliticallyExposed": false,
			"immediateFamilyExposed": false,
			"createdAt": "2026-11-08T10:00:00.000Z",
			"updatedAt": "2026-11-08T10:00:00.000Z"
		}`))
	}))

	result, err := client.BrokerDisclosures.Get(
		context.Background(),
		"user/123",
		"connection/456",
	)
	if err != nil {
		t.Fatalf("BrokerDisclosures.Get() error = %v", err)
	}

	if result == nil {
		t.Fatal("BrokerDisclosures.Get() returned nil, want disclosure")
	}

	if result.ID != "disclosure-123" {
		t.Errorf("result.ID = %q, want disclosure-123", result.ID)
	}

	if result.UserID != "user/123" {
		t.Errorf("result.UserID = %q, want user/123", result.UserID)
	}

	if result.BrokerConnectionID != "connection/456" {
		t.Errorf("result.BrokerConnectionID = %q, want connection/456", result.BrokerConnectionID)
	}

	if result.EmploymentStatus != "employed" {
		t.Errorf("result.EmploymentStatus = %q, want employed", result.EmploymentStatus)
	}

	if result.EmploymentSector == nil || *result.EmploymentSector != "technology" {
		t.Errorf("result.EmploymentSector = %v, want technology", result.EmploymentSector)
	}

	if result.IsControlPerson {
		t.Error("result.IsControlPerson = true, want false")
	}

	if result.IsAffiliatedExchangeOrFINRA {
		t.Error("result.IsAffiliatedExchangeOrFINRA = true, want false")
	}

	if result.IsPoliticallyExposed {
		t.Error("result.IsPoliticallyExposed = true, want false")
	}

	if result.ImmediateFamilyExposed {
		t.Error("result.ImmediateFamilyExposed = true, want false")
	}
}

func TestBrokerDisclosuresGetReturnsNilWhenNoDisclosureExists(t *testing.T) {
	client := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want GET", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte("null"))
	}))

	result, err := client.BrokerDisclosures.Get(
		context.Background(),
		"user-123",
		"connection-456",
	)
	if err != nil {
		t.Fatalf("BrokerDisclosures.Get() error = %v", err)
	}

	if result != nil {
		t.Errorf("BrokerDisclosures.Get() = %+v, want nil", result)
	}
}

func TestBrokerDisclosuresCreate(t *testing.T) {
	employmentSector := "technology"

	params := CreateBrokerDisclosureParams{
		EmploymentStatus:            "employed",
		EmploymentSector:            &employmentSector,
		IsControlPerson:             true,
		IsAffiliatedExchangeOrFINRA: false,
		IsPoliticallyExposed:        true,
		ImmediateFamilyExposed:      false,
	}

	client := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %q, want POST", r.Method)
		}

		if got := r.URL.EscapedPath(); got != "/v1/users/user%2F123/broker-connections/connection%2F456/disclosures" {
			t.Errorf("path = %q, want escaped disclosures path", got)
		}

		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Errorf("decode request body: %v", err)
		} else {
			expected := map[string]any{
				"employmentStatus":            "employed",
				"employmentSector":            "technology",
				"isControlPerson":             true,
				"isAffiliatedExchangeOrFinra": false,
				"isPoliticallyExposed":        true,
				"immediateFamilyExposed":      false,
			}

			for key, want := range expected {
				if body[key] != want {
					t.Errorf("body[%q] = %v, want %v", key, body[key], want)
				}
			}
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"id": "disclosure-789",
			"userId": "user/123",
			"brokerConnectionId": "connection/456",
			"employmentStatus": "employed",
			"employmentSector": "technology",
			"isControlPerson": true,
			"isAffiliatedExchangeOrFinra": false,
			"isPoliticallyExposed": true,
			"immediateFamilyExposed": false,
			"createdAt": "2026-11-08T10:00:00.000Z",
			"updatedAt": "2026-11-08T10:00:00.000Z"
		}`))
	}))

	result, err := client.BrokerDisclosures.Create(
		context.Background(),
		"user/123",
		"connection/456",
		params,
	)
	if err != nil {
		t.Fatalf("BrokerDisclosures.Create() error = %v", err)
	}

	if result.ID != "disclosure-789" {
		t.Errorf("result.ID = %q, want disclosure-789", result.ID)
	}

	if result.UserID != "user/123" {
		t.Errorf("result.UserID = %q, want user/123", result.UserID)
	}

	if result.BrokerConnectionID != "connection/456" {
		t.Errorf("result.BrokerConnectionID = %q, want connection/456", result.BrokerConnectionID)
	}

	if result.EmploymentSector == nil || *result.EmploymentSector != employmentSector {
		t.Errorf("result.EmploymentSector = %v, want %q", result.EmploymentSector, employmentSector)
	}

	if !result.IsControlPerson {
		t.Error("result.IsControlPerson = false, want true")
	}

	if result.IsAffiliatedExchangeOrFINRA {
		t.Error("result.IsAffiliatedExchangeOrFINRA = true, want false")
	}

	if !result.IsPoliticallyExposed {
		t.Error("result.IsPoliticallyExposed = false, want true")
	}

	if result.ImmediateFamilyExposed {
		t.Error("result.ImmediateFamilyExposed = true, want false")
	}
}

func TestBrokerDisclosuresCreateWithoutOptionalEmploymentSector(t *testing.T) {
	client := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Errorf("decode request body: %v", err)
		} else {
			if _, exists := body["employmentSector"]; exists {
				t.Error("employmentSector should be omitted when nil")
			}

			if body["employmentStatus"] != "self_employed" {
				t.Errorf("body employmentStatus = %v, want self_employed", body["employmentStatus"])
			}
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"id": "disclosure-999",
			"userId": "user-123",
			"brokerConnectionId": "connection-456",
			"employmentStatus": "self_employed",
			"employmentSector": null,
			"isControlPerson": false,
			"isAffiliatedExchangeOrFinra": false,
			"isPoliticallyExposed": false,
			"immediateFamilyExposed": false,
			"createdAt": "2026-11-08T10:00:00.000Z",
			"updatedAt": "2026-11-08T10:00:00.000Z"
		}`))
	}))

	result, err := client.BrokerDisclosures.Create(
		context.Background(),
		"user-123",
		"connection-456",
		CreateBrokerDisclosureParams{
			EmploymentStatus: "self_employed",
		},
	)
	if err != nil {
		t.Fatalf("BrokerDisclosures.Create() error = %v", err)
	}

	if result.EmploymentSector != nil {
		t.Errorf("result.EmploymentSector = %v, want nil", *result.EmploymentSector)
	}
}
