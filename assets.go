package alvaindes

import (
	"context"
	"net/url"
)

// List retrieves assets available through a broker connection.
func (r *AssetsResource) List(
	ctx context.Context,
	params ListAssetsParams,
) ([]AssetResponse, error) {
	var result []AssetResponse

	query := url.Values{}
	query.Set("connectionId", params.ConnectionID)
	query.Set("query", params.Query)

	path := "/v1/assets?" + query.Encode()

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
