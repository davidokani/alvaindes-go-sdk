package alvaindes

import (
	"context"
	"net/url"
)

// Get retrieves the portfolio for an investment account.
func (r *PortfoliosResource) Get(
	ctx context.Context,
	userID string,
	accountID string,
) (*PortfolioResponse, error) {
	var result PortfolioResponse

	path := "/v1/users/" +
		url.PathEscape(userID) +
		"/accounts/" +
		url.PathEscape(accountID) +
		"/portfolio"

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
