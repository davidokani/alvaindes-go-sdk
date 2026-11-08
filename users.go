package alvaindes

import (
	"context"
	"net/url"
)

// Create creates a user within an application.
func (r *UsersResource) Create(
	ctx context.Context,
	applicationID string,
	params CreateUserParams,
) (*UserResponse, error) {
	var result UserResponse

	path := "/v1/applications/" + url.PathEscape(applicationID) + "/users"

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

// Get retrieves a user within an application.
func (r *UsersResource) Get(
	ctx context.Context,
	applicationID string,
	userID string,
) (*UserResponse, error) {
	var result UserResponse

	path := "/v1/applications/" +
		url.PathEscape(applicationID) +
		"/users/" +
		url.PathEscape(userID)

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
