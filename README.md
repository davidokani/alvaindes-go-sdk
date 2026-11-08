# Alvaindes Go SDK

<p align="center">
  <img src="./assets/alvaindes.png" alt="Alvaindes Logo" width="140">
</p>

<h1 align="center">Alvaindes SDK</h1>

<p align="center">
 Official Go SDK for integrating with the Alvaindes brokerage infrastructure API.
</p>

<p align="center">
  The SDK provides a broker-agnostic interface for working with users, investment accounts, KYC, assets, funding sources, portfolios, holdings, orders, deposits, withdrawals, transactions, broker agreements, broker disclosures, and market data.
</p>

---

## Requirements

- Go 1.22 or later
- An Alvaindes API key
- Your Alvaindes API base URL

## Installation

```bash
go get github.com/davidokani/alvaindes-go-sdk@v0.1.0
```

## Getting started

Create a client using your API key and API base URL:

```go
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	alvaindes "github.com/davidokani/alvaindes-go-sdk"
)

func main() {
	client, err := alvaindes.New(alvaindes.Config{
		APIKey:  os.Getenv("ALVAINDES_API_KEY"),
		BaseURL: os.Getenv("ALVAINDES_BASE_URL"),
	})
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	user, err := client.Users.Create(ctx, "YOUR_APPLICATION_ID", alvaindes.CreateUserParams{
		FirstName:   "Jane",
		LastName:    "Doe",
		Email:       "jane@example.com",
		PhoneNumber: "+2348012345678",
		Country:     "NG",
	})
	if err != nil {
		var apiErr *alvaindes.AlvaindesError
		if errors.As(err, &apiErr) {
			log.Printf(
				"Alvaindes API error: status=%d code=%s message=%s",
				apiErr.Status,
				apiErr.Code,
				apiErr.Message,
			)
			return
		}

		log.Fatal(err)
	}

	fmt.Printf("Created user: %+v\n", user)
}
```

Set the required environment variables before running your application:

```bash
export ALVAINDES_API_KEY="your_api_key"
export ALVAINDES_BASE_URL="https://api.alvaindes.com"
```

Use the API base URL assigned to your Alvaindes environment. The URL must be the API origin only; do not append `/v1` or another path prefix.

Keep your API key on a trusted server. Do not expose it in browser-side applications, commit it to source control, or include it in client-side mobile or web bundles.

## Configuration

The SDK accepts the following configuration:

| Field        | Description                                                   |
| ------------ | ------------------------------------------------------------- |
| `APIKey`     | Your Alvaindes API key. Required.                             |
| `BaseURL`    | API origin. Required.                                         |
| `Timeout`    | Request timeout as a `time.Duration`. Defaults to 30 seconds. |
| `HTTPClient` | Optional custom `*http.Client`.                               |

Example:

```go
client, err := alvaindes.New(alvaindes.Config{
	APIKey:  os.Getenv("ALVAINDES_API_KEY"),
	BaseURL: os.Getenv("ALVAINDES_BASE_URL"),
	Timeout: 45 * time.Second,
})
if err != nil {
	log.Fatal(err)
}
```

The base URL must use HTTP or HTTPS and must not contain URL credentials, a query, a fragment, or a path prefix. A configured timeout must be positive and represent a whole number of milliseconds.

## Available resources

The SDK exposes the following resource groups through the client:

| Resource           | Methods                                     |
| ------------------ | ------------------------------------------- |
| Users              | `Create`, `Get`                             |
| Accounts           | `Create`, `List`, `Get`                     |
| KYC                | `Submit`, `Get`, `Update`                   |
| Assets             | `List`                                      |
| Funding Sources    | `Create`, `List`, `Refresh`, `Disable`      |
| Portfolios         | `Get`                                       |
| Holdings           | `List`                                      |
| Orders             | `Create`, `List`, `Get`                     |
| Deposits           | `Create`                                    |
| Withdrawals        | `Create`                                    |
| Transactions       | `List`                                      |
| Broker Agreements  | `GetRequirements`, `ListAccepted`, `Accept` |
| Broker Disclosures | `Get`, `Create`                             |
| Market Data        | `GetQuote`                                  |

All resource methods accept a `context.Context`, allowing callers to control cancellation and deadlines.

## Examples

### Create an investment account

```go
account, err := client.Accounts.Create(
	ctx,
	"USER_ID",
	alvaindes.CreateInvestmentAccountParams{
		Currency:           "USD",
		BrokerConnectionID: "YOUR_BROKER_CONNECTION_ID",
	},
)
if err != nil {
	log.Fatal(err)
}

fmt.Printf("Investment account: %+v\n", account)
```

### Create an order

```go
order, err := client.Orders.Create(
	ctx,
	"USER_ID",
	"ACCOUNT_ID",
	alvaindes.CreateOrderParams{
		Symbol:   "AAPL",
		Side:     alvaindes.OrderSideBuy,
		Quantity: 1,
		Type:     alvaindes.OrderTypeMarket,
	},
	alvaindes.CreateOrderOptions{
		IdempotencyKey: "your-unique-request-key",
	},
)
if err != nil {
	log.Fatal(err)
}

fmt.Printf("Order: %+v\n", order)
```

The SDK supports `MARKET` and `LIMIT` orders. A `LIMIT` order requires a limit price; a `MARKET` order must not include one. Use a unique idempotency key when you want supported order-creation requests to be safely retried.

### Create a deposit

```go
deposit, err := client.Deposits.Create(
	ctx,
	"USER_ID",
	"ACCOUNT_ID",
	alvaindes.CreateDepositParams{
		Amount:   100,
		Currency: "USD",
	},
	alvaindes.CreateDepositOptions{
		IdempotencyKey: "your-unique-request-key",
	},
)
if err != nil {
	log.Fatal(err)
}

fmt.Printf("Deposit: %+v\n", deposit)
```

The withdrawal resource follows the same options pattern and also supports an idempotency key. Deposit and withdrawal requests can optionally include a funding source ID.

### Retrieve a market quote

```go
quote, err := client.MarketData.GetQuote(ctx, "AAPL")
if err != nil {
	log.Fatal(err)
}

fmt.Printf("Market quote: %+v\n", quote)
```

## Error handling

The SDK exports `AlvaindesError` for API errors. It provides:

| Field      | Description                   |
| ---------- | ----------------------------- |
| `Status`   | HTTP status code              |
| `Message`  | Error message                 |
| `Code`     | API error code, when provided |
| `Response` | Decoded response payload      |

Use `errors.As` to inspect API errors:

```go
var apiErr *alvaindes.AlvaindesError

if errors.As(err, &apiErr) {
	fmt.Println("HTTP status:", apiErr.Status)
	fmt.Println("Error code:", apiErr.Code)
	fmt.Println("Message:", apiErr.Message)
}
```

Errors that are not API errors can be handled separately as ordinary Go errors.

## Notes

- Funding source creation currently supports ACH. Other funding source types return an error.
- Broker disclosure retrieval can return `nil` when the connected broker has no disclosure.
- Optional request fields are omitted when not supplied.
- Order, deposit, and withdrawal creation support optional idempotency keys.

## License

This project is licensed under the MIT License. See [LICENSE](./LICENSE).
