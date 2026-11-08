package alvaindes

import (
	"context"
	"errors"
	"net/http"
	"net/url"

	"github.com/davidokani/alvaindes-go-sdk/internal/transport"
)

// Create submits a new order for an investment account.
func (r *OrdersResource) Create(
	ctx context.Context,
	userID string,
	accountID string,
	params CreateOrderParams,
	options CreateOrderOptions,
) (*OrderResponse, error) {
	switch params.Type {
	case OrderTypeMarket:
		if params.LimitPrice != nil {
			return nil, errors.New("limitPrice must not be provided for MARKET orders")
		}
	case OrderTypeLimit:
		if params.LimitPrice == nil {
			return nil, errors.New("limitPrice is required for LIMIT orders")
		}
	default:
		return nil, errors.New("order type must be MARKET or LIMIT")
	}

	path := "/v1/users/" +
		url.PathEscape(userID) +
		"/accounts/" +
		url.PathEscape(accountID) +
		"/orders"

	var result OrderResponse

	if options.IdempotencyKey != "" {
		headers := make(http.Header)
		headers.Set("Idempotency-Key", options.IdempotencyKey)

		err := r.transport.RequestWithOptions(
			ctx,
			"POST",
			path,
			params,
			&result,
			transport.RequestOptions{Headers: headers},
		)
		if err != nil {
			return nil, err
		}
	} else {
		if err := r.transport.Request(
			ctx,
			"POST",
			path,
			params,
			&result,
		); err != nil {
			return nil, err
		}
	}

	return &result, nil
}

// List retrieves all orders for an investment account.
func (r *OrdersResource) List(
	ctx context.Context,
	userID string,
	accountID string,
) ([]OrderResponse, error) {
	path := "/v1/users/" +
		url.PathEscape(userID) +
		"/accounts/" +
		url.PathEscape(accountID) +
		"/orders"

	var result []OrderResponse

	if err := r.transport.Request(
		ctx,
		"GET",
		path,
		nil,
		&result,
	); err != nil {
		return nil, err
	}

	return result, nil
}

// Get retrieves a specific order.
func (r *OrdersResource) Get(
	ctx context.Context,
	userID string,
	accountID string,
	orderID string,
) (*OrderResponse, error) {
	path := "/v1/users/" +
		url.PathEscape(userID) +
		"/accounts/" +
		url.PathEscape(accountID) +
		"/orders/" +
		url.PathEscape(orderID)

	var result OrderResponse

	if err := r.transport.Request(
		ctx,
		"GET",
		path,
		nil,
		&result,
	); err != nil {
		return nil, err
	}

	return &result, nil
}
