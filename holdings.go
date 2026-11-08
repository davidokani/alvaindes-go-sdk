package alvaindes

import (
	"context"
	"net/url"
)

// List retrieves the holdings for an investment account.
func (r *HoldingsResource) List(
	ctx context.Context,
	userID string,
	accountID string,
) ([]HoldingResponse, error) {
	var result []HoldingResponse

	path := "/v1/users/" +
		url.PathEscape(userID) +
		"/accounts/" +
		url.PathEscape(accountID) +
		"/holdings"

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
