package alvaindes

import (
	"context"
	"fmt"
	"net/url"
)

// Create creates an ACH funding source for an investment account.
func (r *FundingSourcesResource) Create(
	ctx context.Context,
	userID string,
	accountID string,
	params CreateFundingSourceParams,
) (*FundingSourceResponse, error) {
	if params.Type != FundingSourceACH {
		return nil, fmt.Errorf("funding source creation only supports ACH")
	}

	var result FundingSourceResponse

	path := "/v1/users/" +
		url.PathEscape(userID) +
		"/accounts/" +
		url.PathEscape(accountID) +
		"/funding-sources"

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

// List retrieves the funding sources belonging to an investment account.
func (r *FundingSourcesResource) List(
	ctx context.Context,
	userID string,
	accountID string,
) ([]FundingSourceResponse, error) {
	var result []FundingSourceResponse

	path := "/v1/users/" +
		url.PathEscape(userID) +
		"/accounts/" +
		url.PathEscape(accountID) +
		"/funding-sources"

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

// Refresh refreshes a funding source with the broker.
func (r *FundingSourcesResource) Refresh(
	ctx context.Context,
	userID string,
	accountID string,
	fundingSourceID string,
) (*FundingSourceResponse, error) {
	var result FundingSourceResponse

	path := "/v1/users/" +
		url.PathEscape(userID) +
		"/accounts/" +
		url.PathEscape(accountID) +
		"/funding-sources/" +
		url.PathEscape(fundingSourceID) +
		"/refresh"

	if err := r.transport.Request(
		ctx,
		"POST",
		path,
		nil,
		&result,
	); err != nil {
		return nil, err
	}

	return &result, nil
}

// Disable disables a funding source.
func (r *FundingSourcesResource) Disable(
	ctx context.Context,
	userID string,
	accountID string,
	fundingSourceID string,
) (*FundingSourceResponse, error) {
	var result FundingSourceResponse

	path := "/v1/users/" +
		url.PathEscape(userID) +
		"/accounts/" +
		url.PathEscape(accountID) +
		"/funding-sources/" +
		url.PathEscape(fundingSourceID)

	if err := r.transport.Request(
		ctx,
		"DELETE",
		path,
		nil,
		&result,
	); err != nil {
		return nil, err
	}

	return &result, nil
}
