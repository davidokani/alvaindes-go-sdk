package alvaindes

import (
	"context"
	"net/url"
)

// Get retrieves a user's broker disclosure.
// It returns nil when the broker has no disclosure for the connection.
func (r *BrokerDisclosuresResource) Get(
	ctx context.Context,
	userID string,
	brokerConnectionID string,
) (*BrokerDisclosureResponse, error) {
	path := "/v1/users/" +
		url.PathEscape(userID) +
		"/broker-connections/" +
		url.PathEscape(brokerConnectionID) +
		"/disclosures"

	var result *BrokerDisclosureResponse

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

// Create submits a user's broker disclosure.
func (r *BrokerDisclosuresResource) Create(
	ctx context.Context,
	userID string,
	brokerConnectionID string,
	params CreateBrokerDisclosureParams,
) (*BrokerDisclosureResponse, error) {
	path := "/v1/users/" +
		url.PathEscape(userID) +
		"/broker-connections/" +
		url.PathEscape(brokerConnectionID) +
		"/disclosures"

	var result BrokerDisclosureResponse

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
