package alvaindes

import (
	"context"
	"net/url"
)

// Submit submits KYC information for a user.
func (r *KycResource) Submit(
	ctx context.Context,
	userID string,
	params SubmitKycParams,
) (*KycResponse, error) {
	var result KycResponse

	path := "/v1/users/" + url.PathEscape(userID) + "/kyc"

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

// Get retrieves the KYC information for a user.
func (r *KycResource) Get(
	ctx context.Context,
	userID string,
) (*KycResponse, error) {
	var result KycResponse

	path := "/v1/users/" + url.PathEscape(userID) + "/kyc"

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

// Update partially updates the KYC information for a user.
func (r *KycResource) Update(
	ctx context.Context,
	userID string,
	params UpdateKycParams,
) (*KycResponse, error) {
	var result KycResponse

	path := "/v1/users/" + url.PathEscape(userID) + "/kyc"

	if err := r.transport.Request(
		ctx,
		"PATCH",
		path,
		params,
		&result,
	); err != nil {
		return nil, err
	}

	return &result, nil
}
