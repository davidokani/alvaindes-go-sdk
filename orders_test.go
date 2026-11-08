package alvaindes

import (
	"context"
	"net/http"
	"testing"
)

func TestOrdersCreate(t *testing.T) {
	client := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %q, want POST", r.Method)
		}

		if got := r.URL.EscapedPath(); got != "/v1/users/user%2F123/accounts/account%2F456/orders" {
			t.Errorf("path = %q, want escaped orders path", got)
		}

		if got := r.Header.Get("Idempotency-Key"); got != "order-request-123" {
			t.Errorf("Idempotency-Key = %q, want order-request-123", got)
		}

		if got := r.Header.Get("x-api-key"); got != "test-api-key" {
			t.Errorf("x-api-key = %q, want test-api-key", got)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"id": "order-123",
			"symbol": "AAPL",
			"side": "BUY",
			"quantity": 10,
			"status": "PENDING",
			"brokerOrderId": null,
			"createdAt": "2026-11-08T00:00:00.000Z"
		}`))
	}))

	result, err := client.Orders.Create(
		context.Background(),
		"user/123",
		"account/456",
		CreateOrderParams{
			Symbol:   "AAPL",
			Side:     OrderSideBuy,
			Quantity: 10,
			Type:     OrderTypeMarket,
		},
		CreateOrderOptions{
			IdempotencyKey: "order-request-123",
		},
	)
	if err != nil {
		t.Fatalf("Orders.Create() error = %v", err)
	}

	if result.ID != "order-123" {
		t.Errorf("result.ID = %q, want order-123", result.ID)
	}

	if result.Symbol != "AAPL" {
		t.Errorf("result.Symbol = %q, want AAPL", result.Symbol)
	}

	if result.Quantity != 10 {
		t.Errorf("result.Quantity = %v, want 10", result.Quantity)
	}

	if result.BrokerOrderID != nil {
		t.Errorf("result.BrokerOrderID = %v, want nil", *result.BrokerOrderID)
	}
}

func TestOrdersCreateLimitOrder(t *testing.T) {
	client := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %q, want POST", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"id": "order-limit-123",
			"symbol": "MSFT",
			"side": "SELL",
			"quantity": 5,
			"status": "PENDING",
			"brokerOrderId": "broker-order-123",
			"createdAt": "2026-11-08T00:00:00.000Z"
		}`))
	}))

	limitPrice := 420.50

	result, err := client.Orders.Create(
		context.Background(),
		"user-123",
		"account-456",
		CreateOrderParams{
			Symbol:     "MSFT",
			Side:       OrderSideSell,
			Quantity:   5,
			Type:       OrderTypeLimit,
			LimitPrice: &limitPrice,
		},
		CreateOrderOptions{},
	)
	if err != nil {
		t.Fatalf("Orders.Create() error = %v", err)
	}

	if result.ID != "order-limit-123" {
		t.Errorf("result.ID = %q, want order-limit-123", result.ID)
	}

	if result.BrokerOrderID == nil || *result.BrokerOrderID != "broker-order-123" {
		t.Errorf("result.BrokerOrderID = %v, want broker-order-123", result.BrokerOrderID)
	}
}

func TestOrdersCreateRejectsInvalidLimitPriceCombinations(t *testing.T) {
	requestCount := 0

	client := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++
	}))

	limitPrice := 100.00

	tests := []struct {
		name   string
		params CreateOrderParams
	}{
		{
			name: "market order with limit price",
			params: CreateOrderParams{
				Symbol:     "AAPL",
				Side:       OrderSideBuy,
				Quantity:   1,
				Type:       OrderTypeMarket,
				LimitPrice: &limitPrice,
			},
		},
		{
			name: "limit order without limit price",
			params: CreateOrderParams{
				Symbol:   "AAPL",
				Side:     OrderSideBuy,
				Quantity: 1,
				Type:     OrderTypeLimit,
			},
		},
		{
			name: "unsupported order type",
			params: CreateOrderParams{
				Symbol:   "AAPL",
				Side:     OrderSideBuy,
				Quantity: 1,
				Type:     OrderType("STOP"),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := client.Orders.Create(
				context.Background(),
				"user-123",
				"account-456",
				test.params,
				CreateOrderOptions{},
			)

			if err == nil {
				t.Fatal("Orders.Create() error = nil, want validation error")
			}
		})
	}

	if requestCount != 0 {
		t.Errorf("requestCount = %d, want 0 for invalid orders", requestCount)
	}
}

func TestOrdersList(t *testing.T) {
	client := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want GET", r.Method)
		}

		if got := r.URL.EscapedPath(); got != "/v1/users/user%2F123/accounts/account%2F456/orders" {
			t.Errorf("path = %q, want escaped orders path", got)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[
			{
				"id": "order-123",
				"symbol": "AAPL",
				"side": "BUY",
				"quantity": 10,
				"status": "FILLED",
				"brokerOrderId": "broker-123",
				"createdAt": "2026-11-08T00:00:00.000Z"
			},
			{
				"id": "order-456",
				"symbol": "MSFT",
				"side": "SELL",
				"quantity": 5,
				"status": "PENDING",
				"brokerOrderId": null,
				"createdAt": "2026-11-08T00:00:00.000Z"
			}
		]`))
	}))

	result, err := client.Orders.List(
		context.Background(),
		"user/123",
		"account/456",
	)
	if err != nil {
		t.Fatalf("Orders.List() error = %v", err)
	}

	if len(result) != 2 {
		t.Fatalf("len(result) = %d, want 2", len(result))
	}

	if result[0].ID != "order-123" {
		t.Errorf("result[0].ID = %q, want order-123", result[0].ID)
	}

	if result[1].ID != "order-456" {
		t.Errorf("result[1].ID = %q, want order-456", result[1].ID)
	}

	if result[1].BrokerOrderID != nil {
		t.Errorf("result[1].BrokerOrderID = %v, want nil", *result[1].BrokerOrderID)
	}
}

func TestOrdersGet(t *testing.T) {
	client := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want GET", r.Method)
		}

		if got := r.URL.EscapedPath(); got != "/v1/users/user%2F123/accounts/account%2F456/orders/order%2F789" {
			t.Errorf("path = %q, want escaped order path", got)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"id": "order/789",
			"symbol": "NVDA",
			"side": "BUY",
			"quantity": 2,
			"status": "FILLED",
			"brokerOrderId": "broker-789",
			"createdAt": "2026-11-08T00:00:00.000Z"
		}`))
	}))

	result, err := client.Orders.Get(
		context.Background(),
		"user/123",
		"account/456",
		"order/789",
	)
	if err != nil {
		t.Fatalf("Orders.Get() error = %v", err)
	}

	if result.ID != "order/789" {
		t.Errorf("result.ID = %q, want order/789", result.ID)
	}

	if result.Symbol != "NVDA" {
		t.Errorf("result.Symbol = %q, want NVDA", result.Symbol)
	}
}
