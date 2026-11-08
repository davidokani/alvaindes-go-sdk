package alvaindes

import (
	"context"
	"net/url"
)

// Create creates an investment account for a user.
func (r *AccountsResource) Create(
	ctx context.Context,
	userID string,
	params CreateInvestmentAccountParams,
) (*InvestmentAccountResponse, error) {
	var result InvestmentAccountResponse

	path := "/v1/users/" + url.PathEscape(userID) + "/accounts"

	if err := r.transport.Request(
		ctx,
		"POST",
		path,
		params,
		&result,
	); err != nil {
		return nil, err
	}

	return &result, nil
}

// List retrieves the investment accounts belonging to a user.
func (r *AccountsResource) List(
	ctx context.Context,
	userID string,
) ([]InvestmentAccountResponse, error) {
	var result []InvestmentAccountResponse

	path := "/v1/users/" + url.PathEscape(userID) + "/accounts"

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

// Get retrieves a specific investment account belonging to a user.
func (r *AccountsResource) Get(
	ctx context.Context,
	userID string,
	accountID string,
) (*InvestmentAccountResponse, error) {
	var result InvestmentAccountResponse

	path := "/v1/users/" +
		url.PathEscape(userID) +
		"/accounts/" +
		url.PathEscape(accountID)

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
