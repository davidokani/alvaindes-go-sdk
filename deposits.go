package alvaindes

import (
	"context"
	"net/http"
	"net/url"

	"github.com/davidokani/alvaindes-go-sdk/internal/transport"
)

// Create submits a deposit for an investment account.
func (r *DepositsResource) Create(
	ctx context.Context,
	userID string,
	accountID string,
	params CreateDepositParams,
	options CreateDepositOptions,
) (*DepositResponse, error) {
	path := "/v1/users/" +
		url.PathEscape(userID) +
		"/accounts/" +
		url.PathEscape(accountID) +
		"/deposits"

	var result DepositResponse

	if options.IdempotencyKey != "" {
		headers := make(http.Header)
		headers.Set("Idempotency-Key", options.IdempotencyKey)

		if err := r.transport.RequestWithOptions(
			ctx,
			"POST",
			path,
			params,
			&result,
			transport.RequestOptions{Headers: headers},
		); err != nil {
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
