package alvaindes

import (
	"context"
	"net/url"
)

// GetQuote retrieves the latest market quote for an asset symbol.
func (r *MarketDataResource) GetQuote(
	ctx context.Context,
	symbol string,
) (*MarketQuoteResponse, error) {
	path := "/v1/market-data/" + url.PathEscape(symbol)

	var result MarketQuoteResponse

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
