package alvaindes

import (
	"context"
	"net/url"
)

// GetRequirements retrieves the agreements required by a broker connection.
func (r *BrokerAgreementsResource) GetRequirements(
	ctx context.Context,
	userID string,
	brokerConnectionID string,
) ([]BrokerAgreementRequirement, error) {
	path := "/v1/users/" +
		url.PathEscape(userID) +
		"/broker-connections/" +
		url.PathEscape(brokerConnectionID) +
		"/agreements/requirements"

	var result []BrokerAgreementRequirement

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

// ListAccepted retrieves the agreements accepted by a user for a broker connection.
func (r *BrokerAgreementsResource) ListAccepted(
	ctx context.Context,
	userID string,
	brokerConnectionID string,
) ([]BrokerAgreementResponse, error) {
	path := "/v1/users/" +
		url.PathEscape(userID) +
		"/broker-connections/" +
		url.PathEscape(brokerConnectionID) +
		"/agreements"

	var result []BrokerAgreementResponse

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

// Accept submits acceptance of a broker agreement.
func (r *BrokerAgreementsResource) Accept(
	ctx context.Context,
	userID string,
	brokerConnectionID string,
	params AcceptBrokerAgreementParams,
) (*BrokerAgreementResponse, error) {
	path := "/v1/users/" +
		url.PathEscape(userID) +
		"/broker-connections/" +
		url.PathEscape(brokerConnectionID) +
		"/agreements"

	var result BrokerAgreementResponse

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
