package alvaindes

import (
	"context"
	"net/url"
)

// List retrieves the transactions for an investment account.
func (r *TransactionsResource) List(
	ctx context.Context,
	userID string,
	accountID string,
) ([]TransactionResponse, error) {
	path := "/v1/users/" +
		url.PathEscape(userID) +
		"/accounts/" +
		url.PathEscape(accountID) +
		"/transactions"

	var result []TransactionResponse

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
