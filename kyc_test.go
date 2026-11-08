package alvaindes

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestKycSubmit(t *testing.T) {
	client := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %q, want POST", r.Method)
		}

		if got := r.URL.EscapedPath(); got != "/v1/users/user%2F123/kyc" {
			t.Errorf("path = %q, want escaped user ID", got)
		}

		if got := r.Header.Get("x-api-key"); got != "test-api-key" {
			t.Errorf("x-api-key = %q, want test-api-key", got)
		}

		var body SubmitKycParams
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode request body: %v", err)
		}

		if body.CountryOfResidence != "NG" {
			t.Errorf("CountryOfResidence = %q, want NG", body.CountryOfResidence)
		}

		if body.DateOfBirth != "1995-05-20" {
			t.Errorf("DateOfBirth = %q, want 1995-05-20", body.DateOfBirth)
		}

		if body.DocumentType != KycDocumentPassport {
			t.Errorf("DocumentType = %q, want PASSPORT", body.DocumentType)
		}

		if body.AddressCity != "Lagos" {
			t.Errorf("AddressCity = %q, want Lagos", body.AddressCity)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"id": "kyc-123",
			"userId": "user/123",
			"status": "PENDING",
			"countryOfResidence": "NG",
			"nationality": "NG",
			"dateOfBirth": "1995-05-20",
			"addressLine1": "10 Marina Road",
			"addressLine2": null,
			"addressCity": "Lagos",
			"addressState": "Lagos",
			"addressPostalCode": null,
			"documentType": "PASSPORT",
			"documentNumber": "A12345678",
			"documentUrl": null,
			"documentExpiryDate": null,
			"documentCountry": "NG",
			"approvedAt": null,
			"rejectedReason": null,
			"createdAt": "2026-11-08T00:00:00.000Z",
			"updatedAt": "2026-11-08T00:00:00.000Z"
		}`))
	}))

	result, err := client.KYC.Submit(
		context.Background(),
		"user/123",
		SubmitKycParams{
			CountryOfResidence: "NG",
			DateOfBirth:        "1995-05-20",
			AddressLine1:       "10 Marina Road",
			AddressCity:        "Lagos",
			DocumentType:       KycDocumentPassport,
		},
	)
	if err != nil {
		t.Fatalf("KYC.Submit() error = %v", err)
	}

	if result.ID != "kyc-123" {
		t.Errorf("result.ID = %q, want kyc-123", result.ID)
	}

	if result.Status != "PENDING" {
		t.Errorf("result.Status = %q, want PENDING", result.Status)
	}

	if result.AddressLine2 != nil {
		t.Errorf("result.AddressLine2 = %v, want nil", result.AddressLine2)
	}
}

func TestKycGet(t *testing.T) {
	client := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want GET", r.Method)
		}

		if got := r.URL.EscapedPath(); got != "/v1/users/user%2F456/kyc" {
			t.Errorf("path = %q, want escaped user ID", got)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"id": "kyc-456",
			"userId": "user/456",
			"status": "APPROVED",
			"countryOfResidence": "NG",
			"nationality": "NG",
			"dateOfBirth": "1990-01-15",
			"addressLine1": "25 Broad Street",
			"addressLine2": null,
			"addressCity": "Lagos",
			"addressState": "Lagos",
			"addressPostalCode": null,
			"documentType": "NATIONAL_ID",
			"documentNumber": null,
			"documentUrl": null,
			"documentExpiryDate": null,
			"documentCountry": "NG",
			"approvedAt": "2026-11-07T10:00:00.000Z",
			"rejectedReason": null,
			"createdAt": "2026-11-01T00:00:00.000Z",
			"updatedAt": "2026-11-07T10:00:00.000Z"
		}`))
	}))

	result, err := client.KYC.Get(context.Background(), "user/456")
	if err != nil {
		t.Fatalf("KYC.Get() error = %v", err)
	}

	if result.ID != "kyc-456" {
		t.Errorf("result.ID = %q, want kyc-456", result.ID)
	}

	if result.Status != "APPROVED" {
		t.Errorf("result.Status = %q, want APPROVED", result.Status)
	}

	if result.ApprovedAt == nil {
		t.Fatal("result.ApprovedAt = nil, want approval timestamp")
	}
}

func TestKycUpdate(t *testing.T) {
	client := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("method = %q, want PATCH", r.Method)
		}

		if got := r.URL.EscapedPath(); got != "/v1/users/user%2F789/kyc" {
			t.Errorf("path = %q, want escaped user ID", got)
		}

		var body map[string]json.RawMessage
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode request body: %v", err)
		}

		if _, ok := body["addressCity"]; !ok {
			t.Error("request body is missing addressCity")
		}

		if _, ok := body["addressLine1"]; ok {
			t.Error("request body unexpectedly contains addressLine1")
		}

		var addressCity string
		if err := json.Unmarshal(body["addressCity"], &addressCity); err != nil {
			t.Fatalf("decode addressCity: %v", err)
		}

		if addressCity != "Abuja" {
			t.Errorf("addressCity = %q, want Abuja", addressCity)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"id": "kyc-789",
			"userId": "user/789",
			"status": "PENDING",
			"countryOfResidence": "NG",
			"nationality": "NG",
			"dateOfBirth": "1992-03-10",
			"addressLine1": "10 Example Street",
			"addressLine2": null,
			"addressCity": "Abuja",
			"addressState": "FCT",
			"addressPostalCode": null,
			"documentType": "PASSPORT",
			"documentNumber": null,
			"documentUrl": null,
			"documentExpiryDate": null,
			"documentCountry": "NG",
			"approvedAt": null,
			"rejectedReason": null,
			"createdAt": "2026-11-01T00:00:00.000Z",
			"updatedAt": "2026-11-08T00:00:00.000Z"
		}`))
	}))

	addressCity := "Abuja"

	result, err := client.KYC.Update(
		context.Background(),
		"user/789",
		UpdateKycParams{
			AddressCity: &addressCity,
		},
	)
	if err != nil {
		t.Fatalf("KYC.Update() error = %v", err)
	}

	if result.ID != "kyc-789" {
		t.Errorf("result.ID = %q, want kyc-789", result.ID)
	}

	if result.AddressCity == nil || *result.AddressCity != "Abuja" {
		t.Errorf("result.AddressCity = %v, want Abuja", result.AddressCity)
	}
}
