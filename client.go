package alvaindes

import (
	"context"
	"net/http"

	"github.com/davidokani/alvaindes-go-sdk/internal/transport"
)

// RequestOptions contains optional settings for a direct API request.
type RequestOptions struct {
	Headers http.Header
}

// Client is the main Alvaindes SDK client.
type Client struct {
	transport *transport.Transport

	Users             *UsersResource
	Accounts          *AccountsResource
	KYC               *KycResource
	Assets            *AssetsResource
	FundingSources    *FundingSourcesResource
	Portfolios        *PortfoliosResource
	Holdings          *HoldingsResource
	Orders            *OrdersResource
	Deposits          *DepositsResource
	Withdrawals       *WithdrawalsResource
	Transactions      *TransactionsResource
	BrokerAgreements  *BrokerAgreementsResource
	BrokerDisclosures *BrokerDisclosuresResource
	MarketData        *MarketDataResource
}

// New creates an Alvaindes SDK client.
func New(config Config) (*Client, error) {
	resolved, err := resolveConfig(config)
	if err != nil {
		return nil, err
	}

	httpTransport := transport.New(transport.Config{
		APIKey:     resolved.apiKey,
		BaseURL:    resolved.baseURL,
		Timeout:    resolved.timeout,
		HTTPClient: resolved.httpClient,
	})

	client := &Client{
		transport: httpTransport,
	}

	client.Users = &UsersResource{transport: httpTransport}
	client.Accounts = &AccountsResource{transport: httpTransport}
	client.KYC = &KycResource{transport: httpTransport}
	client.Assets = &AssetsResource{transport: httpTransport}
	client.FundingSources = &FundingSourcesResource{transport: httpTransport}
	client.Portfolios = &PortfoliosResource{transport: httpTransport}
	client.Holdings = &HoldingsResource{transport: httpTransport}
	client.Orders = &OrdersResource{transport: httpTransport}
	client.Deposits = &DepositsResource{transport: httpTransport}
	client.Withdrawals = &WithdrawalsResource{transport: httpTransport}
	client.Transactions = &TransactionsResource{transport: httpTransport}
	client.BrokerAgreements = &BrokerAgreementsResource{transport: httpTransport}
	client.BrokerDisclosures = &BrokerDisclosuresResource{transport: httpTransport}
	client.MarketData = &MarketDataResource{transport: httpTransport}

	return client, nil
}

// Request executes a direct API request.
// Use a typed result pointer to decode a JSON response.
func (c *Client) Request(
	ctx context.Context,
	method string,
	path string,
	body any,
	result any,
) error {
	return c.transport.Request(ctx, method, path, body, result)
}

// RequestWithOptions executes a direct API request with custom headers.
func (c *Client) RequestWithOptions(
	ctx context.Context,
	method string,
	path string,
	body any,
	result any,
	options RequestOptions,
) error {
	return c.transport.RequestWithOptions(
		ctx,
		method,
		path,
		body,
		result,
		transport.RequestOptions{
			Headers: options.Headers,
		},
	)
}
