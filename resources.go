package alvaindes

import "github.com/davidokani/alvaindes-go-sdk/internal/transport"

// Resource handles share the client's configured HTTP transport.

type UsersResource struct {
	transport *transport.Transport
}

type AccountsResource struct {
	transport *transport.Transport
}

type KycResource struct {
	transport *transport.Transport
}

type AssetsResource struct {
	transport *transport.Transport
}

type FundingSourcesResource struct {
	transport *transport.Transport
}

type PortfoliosResource struct {
	transport *transport.Transport
}

type HoldingsResource struct {
	transport *transport.Transport
}

type OrdersResource struct {
	transport *transport.Transport
}

type DepositsResource struct {
	transport *transport.Transport
}

type WithdrawalsResource struct {
	transport *transport.Transport
}

type TransactionsResource struct {
	transport *transport.Transport
}

type BrokerAgreementsResource struct {
	transport *transport.Transport
}

type BrokerDisclosuresResource struct {
	transport *transport.Transport
}

type MarketDataResource struct {
	transport *transport.Transport
}
