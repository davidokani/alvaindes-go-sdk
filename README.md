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
- An Alvaindes application and the relevant broker connections configured for the operations you intend to perform

## Installation

```bash
go get github.com/davidokani/alvaindes-go-sdk
```

Import the package:

```go
import alvaindes "github.com/davidokani/alvaindes-go-sdk"
```

## Getting started

Create a client using your API key and the Alvaindes API base URL.

```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	alvaindes "github.com/davidokani/alvaindes-go-sdk"
)

func main() {
	apiKey := os.Getenv("ALVAINDES_API_KEY")
	if apiKey == "" {
		log.Fatal("ALVAINDES_API_KEY is required")
	}

	client, err := alvaindes.New(alvaindes.Config{
		APIKey:  apiKey,
		BaseURL: "https://api.alvaindes.com",
		Timeout: 30 * time.Second,
	})
	if err != nil {
		log.Fatalf("create Alvaindes client: %v", err)
	}

	ctx := context.Background()

	_ = ctx
	_ = client

	fmt.Println("Alvaindes client initialized")
}
```

### Configuration

`alvaindes.Config` accepts the following fields:

| Field        | Type            | Description                                                |
| ------------ | --------------- | ---------------------------------------------------------- |
| `APIKey`     | `string`        | Your Alvaindes API key. Required.                          |
| `BaseURL`    | `string`        | API origin. Required.                                      |
| `Timeout`    | `time.Duration` | HTTP request timeout. Defaults to 30 seconds when omitted. |
| `HTTPClient` | `*http.Client`  | Optional custom HTTP client.                               |

The `BaseURL` must be an HTTP or HTTPS origin. Do not append `/v1` or an endpoint path to it.

For example:

```go
client, err := alvaindes.New(alvaindes.Config{
	APIKey:  os.Getenv("ALVAINDES_API_KEY"),
	BaseURL: "https://api.alvaindes.com",
})
if err != nil {
	return err
}
```

The default request timeout is 30 seconds. You can provide a different positive timeout using `time.Duration`.

Keep API keys on trusted server-side infrastructure. Do not expose secret API keys in browser code, mobile applications, public repositories, or client-side bundles.

## Available resources

The SDK exposes the following resource groups through the client:

| Resource            | Description                                                                       |
| ------------------- | --------------------------------------------------------------------------------- |
| `Users`             | Create and retrieve users.                                                        |
| `Accounts`          | Create, list, and retrieve investment accounts.                                   |
| `KYC`               | Submit, retrieve, and update KYC information.                                     |
| `Assets`            | List available assets.                                                            |
| `FundingSources`    | Create, list, refresh, and disable funding sources.                               |
| `Portfolios`        | Retrieve portfolio information.                                                   |
| `Holdings`          | List investment holdings.                                                         |
| `Orders`            | Create, list, and retrieve orders.                                                |
| `Deposits`          | Create deposits.                                                                  |
| `Withdrawals`       | Create withdrawals.                                                               |
| `Transactions`      | List account transactions.                                                        |
| `BrokerAgreements`  | Retrieve agreement requirements, list accepted agreements, and accept agreements. |
| `BrokerDisclosures` | Retrieve and create broker disclosures.                                           |
| `MarketData`        | Retrieve market quotes.                                                           |

All resource methods accept a `context.Context` so requests can be cancelled or controlled with deadlines.

## Usage examples

The examples below assume that `client` has already been initialized and that the required IDs and context are available.

### 1. Users

Create a user under an Alvaindes application.

```go
user, err := client.Users.Create(
	ctx,
	"YOUR_APPLICATION_ID",
	alvaindes.CreateUserParams{
		FirstName:   "Jane",
		LastName:    "Doe",
		Email:       "jane@example.com",
		PhoneNumber: "+2348012345678",
		Country:     "NG",
	},
)
if err != nil {
	return err
}

fmt.Println("User ID:", user.ID)
```

You can optionally provide an external identifier:

```go
externalID := "customer-123"

user, err := client.Users.Create(
	ctx,
	"YOUR_APPLICATION_ID",
	alvaindes.CreateUserParams{
		ExternalID:  &externalID,
		FirstName:   "Jane",
		LastName:    "Doe",
		Email:       "jane@example.com",
		PhoneNumber: "+2348012345678",
		Country:     "NG",
	},
)
```

Retrieve a user:

```go
user, err := client.Users.Get(
	ctx,
	"YOUR_APPLICATION_ID",
	"YOUR_USER_ID",
)
if err != nil {
	return err
}

fmt.Println(user.FirstName, user.LastName)
```

### 2. Investment accounts

Create an investment account for a user.

```go
account, err := client.Accounts.Create(
	ctx,
	"YOUR_USER_ID",
	alvaindes.CreateInvestmentAccountParams{
		Currency:           "USD",
		BrokerConnectionID: "YOUR_BROKER_CONNECTION_ID",
	},
)
if err != nil {
	return err
}

fmt.Println("Account ID:", account.ID)
```

List a user's investment accounts:

```go
accounts, err := client.Accounts.List(ctx, "YOUR_USER_ID")
if err != nil {
	return err
}

for _, account := range accounts {
	fmt.Println(account.ID, account.Currency, account.Status)
}
```

Retrieve a specific account:

```go
account, err := client.Accounts.Get(
	ctx,
	"YOUR_USER_ID",
	"YOUR_ACCOUNT_ID",
)
if err != nil {
	return err
}

fmt.Println("Account number:", account.AccountNumber)
```

### 3. Know Your Customer (KYC)

Submit KYC information for a user.

```go
nationality := "NG"
documentNumber := "A12345678"

kyc, err := client.KYC.Submit(
	ctx,
	"YOUR_USER_ID",
	alvaindes.SubmitKycParams{
		CountryOfResidence: "NG",
		Nationality:        &nationality,
		DateOfBirth:        "1995-06-15",
		AddressLine1:       "12 Example Street",
		AddressCity:        "Lagos",
		DocumentType:       alvaindes.KycDocumentNationalID,
		DocumentNumber:     &documentNumber,
	},
)
if err != nil {
	return err
}

fmt.Println("KYC status:", kyc.Status)
```

Supported document types:

- `alvaindes.KycDocumentPassport`
- `alvaindes.KycDocumentNationalID`
- `alvaindes.KycDocumentDriversLicense`
- `alvaindes.KycDocumentVotersCard`

Optional KYC fields use pointers. A nil pointer is omitted from the JSON request.

Retrieve KYC information:

```go
kyc, err := client.KYC.Get(ctx, "YOUR_USER_ID")
if err != nil {
	return err
}

fmt.Println("KYC status:", kyc.Status)
```

Update existing KYC information:

```go
city := "Abuja"

kyc, err := client.KYC.Update(
	ctx,
	"YOUR_USER_ID",
	alvaindes.UpdateKycParams{
		AddressCity: &city,
	},
)
if err != nil {
	return err
}
```

`UpdateKycParams` supports partial updates. Only fields supplied with non-nil pointers are included in the request.

### 4. Assets

List available assets, optionally filtering by broker connection or search query.

```go
assets, err := client.Assets.List(
	ctx,
	alvaindes.ListAssetsParams{
		ConnectionID: "YOUR_BROKER_CONNECTION_ID",
		Query:        "AAPL",
	},
)
if err != nil {
	return err
}

for _, asset := range assets {
	fmt.Println(asset.Symbol, asset.Name, asset.Exchange)
}
```

### 5. Funding sources

Create a funding source.

The current Go SDK supports ACH funding-source creation.

```go
fundingSource, err := client.FundingSources.Create(
	ctx,
	"YOUR_USER_ID",
	"YOUR_ACCOUNT_ID",
	alvaindes.CreateFundingSourceParams{
		Type:              alvaindes.FundingSourceACH,
		OwnerName:         "Jane Doe",
		BankAccountNumber: "1234567890",
		BankAccountType:   "CHECKING",
		RoutingNumber:     "021000021",
		Currency:          "USD",
	},
)
if err != nil {
	return err
}

fmt.Println("Funding source ID:", fundingSource.ID)
```

List funding sources:

```go
fundingSources, err := client.FundingSources.List(
	ctx,
	"YOUR_USER_ID",
	"YOUR_ACCOUNT_ID",
)
if err != nil {
	return err
}

for _, source := range fundingSources {
	fmt.Println(source.ID, source.Type, source.Status)
}
```

Refresh a funding source:

```go
fundingSource, err := client.FundingSources.Refresh(
	ctx,
	"YOUR_USER_ID",
	"YOUR_ACCOUNT_ID",
	"YOUR_FUNDING_SOURCE_ID",
)
if err != nil {
	return err
}
```

Disable a funding source:

```go
fundingSource, err := client.FundingSources.Disable(
	ctx,
	"YOUR_USER_ID",
	"YOUR_ACCOUNT_ID",
	"YOUR_FUNDING_SOURCE_ID",
)
if err != nil {
	return err
}
```

Funding-source types and statuses are represented by:

- `FundingSourceACH`
- `FundingSourceBankAccount`
- `FundingSourceWire`
- `FundingSourceOther`
- `FundingSourcePending`
- `FundingSourceActive`
- `FundingSourceDisabled`

The presence of a type constant does not mean that creation through the SDK is currently supported for that type.

### 6. Portfolios

Retrieve a user's portfolio for an investment account.

```go
portfolio, err := client.Portfolios.Get(
	ctx,
	"YOUR_USER_ID",
	"YOUR_ACCOUNT_ID",
)
if err != nil {
	return err
}

fmt.Println("Account value:", portfolio.AccountValue)
fmt.Println("Cash balance:", portfolio.CashBalance)
fmt.Println("Currency:", portfolio.Currency)
```

### 7. Holdings

List holdings in an investment account.

```go
holdings, err := client.Holdings.List(
	ctx,
	"YOUR_USER_ID",
	"YOUR_ACCOUNT_ID",
)
if err != nil {
	return err
}

for _, holding := range holdings {
	fmt.Println(
		holding.Symbol,
		holding.Quantity,
		holding.AveragePrice,
		holding.MarketValue,
	)
}
```

### 8. Orders

Create a market order.

```go
order, err := client.Orders.Create(
	ctx,
	"YOUR_USER_ID",
	"YOUR_ACCOUNT_ID",
	alvaindes.CreateOrderParams{
		Symbol:   "AAPL",
		Side:     alvaindes.OrderSideBuy,
		Quantity: 1,
		Type:     alvaindes.OrderTypeMarket,
	},
	alvaindes.CreateOrderOptions{
		IdempotencyKey: "buy-aapl-001",
	},
)
if err != nil {
	return err
}

fmt.Println("Order ID:", order.ID)
fmt.Println("Order status:", order.Status)
```

Create a limit order by providing a positive limit price:

```go
limitPrice := 190.00

order, err := client.Orders.Create(
	ctx,
	"YOUR_USER_ID",
	"YOUR_ACCOUNT_ID",
	alvaindes.CreateOrderParams{
		Symbol:     "AAPL",
		Side:       alvaindes.OrderSideBuy,
		Quantity:   1,
		Type:       alvaindes.OrderTypeLimit,
		LimitPrice: &limitPrice,
	},
	alvaindes.CreateOrderOptions{
		IdempotencyKey: "buy-aapl-limit-001",
	},
)
if err != nil {
	return err
}
```

Order sides:

- `alvaindes.OrderSideBuy`
- `alvaindes.OrderSideSell`

Order types:

- `alvaindes.OrderTypeMarket`
- `alvaindes.OrderTypeLimit`

A market order must not include `LimitPrice`. A limit order requires a positive `LimitPrice`.

List orders:

```go
orders, err := client.Orders.List(
	ctx,
	"YOUR_USER_ID",
	"YOUR_ACCOUNT_ID",
)
if err != nil {
	return err
}

for _, order := range orders {
	fmt.Println(order.ID, order.Symbol, order.Status)
}
```

Retrieve an order:

```go
order, err := client.Orders.Get(
	ctx,
	"YOUR_USER_ID",
	"YOUR_ACCOUNT_ID",
	"YOUR_ORDER_ID",
)
if err != nil {
	return err
}

fmt.Println("Order status:", order.Status)
```

### 9. Deposits

Create a deposit.

```go
deposit, err := client.Deposits.Create(
	ctx,
	"YOUR_USER_ID",
	"YOUR_ACCOUNT_ID",
	alvaindes.CreateDepositParams{
		Amount:   100,
		Currency: "USD",
	},
	alvaindes.CreateDepositOptions{
		IdempotencyKey: "deposit-001",
	},
)
if err != nil {
	return err
}

fmt.Println("Deposit ID:", deposit.ID)
fmt.Println("Deposit status:", deposit.Status)
```

You can optionally specify a funding source:

```go
fundingSourceID := "YOUR_FUNDING_SOURCE_ID"

deposit, err := client.Deposits.Create(
	ctx,
	"YOUR_USER_ID",
	"YOUR_ACCOUNT_ID",
	alvaindes.CreateDepositParams{
		Amount:          100,
		Currency:        "USD",
		FundingSourceID: &fundingSourceID,
	},
	alvaindes.CreateDepositOptions{
		IdempotencyKey: "deposit-002",
	},
)
if err != nil {
	return err
}
```

When `FundingSourceID` is omitted, the API may resolve an eligible funding source according to its funding-source selection rules.

### 10. Withdrawals

Create a withdrawal.

```go
withdrawal, err := client.Withdrawals.Create(
	ctx,
	"YOUR_USER_ID",
	"YOUR_ACCOUNT_ID",
	alvaindes.CreateWithdrawalParams{
		Amount:   50,
		Currency: "USD",
	},
	alvaindes.CreateWithdrawalOptions{
		IdempotencyKey: "withdrawal-001",
	},
)
if err != nil {
	return err
}

fmt.Println("Withdrawal ID:", withdrawal.ID)
fmt.Println("Withdrawal status:", withdrawal.Status)
```

You can optionally specify a funding source:

```go
fundingSourceID := "YOUR_FUNDING_SOURCE_ID"

withdrawal, err := client.Withdrawals.Create(
	ctx,
	"YOUR_USER_ID",
	"YOUR_ACCOUNT_ID",
	alvaindes.CreateWithdrawalParams{
		Amount:          50,
		Currency:        "USD",
		FundingSourceID: &fundingSourceID,
	},
	alvaindes.CreateWithdrawalOptions{
		IdempotencyKey: "withdrawal-002",
	},
)
if err != nil {
	return err
}
```

Withdrawals are subject to the API's account, available-cash, and funding validations.

### 11. Transactions

List transactions for an investment account.

```go
transactions, err := client.Transactions.List(
	ctx,
	"YOUR_USER_ID",
	"YOUR_ACCOUNT_ID",
)
if err != nil {
	return err
}

for _, transaction := range transactions {
	fmt.Println(
		transaction.ID,
		transaction.Type,
		transaction.Amount,
		transaction.Currency,
		transaction.Status,
	)
}
```

### 12. Broker agreements

Retrieve the agreement requirements for a broker connection.

```go
requirements, err := client.BrokerAgreements.GetRequirements(
	ctx,
	"YOUR_USER_ID",
	"YOUR_BROKER_CONNECTION_ID",
)
if err != nil {
	return err
}

for _, requirement := range requirements {
	fmt.Println("Agreement:", requirement.Agreement)

	if requirement.Revision != nil {
		fmt.Println("Revision:", *requirement.Revision)
	}
}
```

Accept an agreement:

```go
agreement, err := client.BrokerAgreements.Accept(
	ctx,
	"YOUR_USER_ID",
	"YOUR_BROKER_CONNECTION_ID",
	alvaindes.AcceptBrokerAgreementParams{
		Agreement: "AGREEMENT_NAME_FROM_API",
	},
)
if err != nil {
	return err
}

fmt.Println("Accepted agreement:", agreement.Agreement)
```

If the API supplies a revision, include that revision when accepting the agreement:

```go
revision := "REVISION_FROM_API"

agreement, err := client.BrokerAgreements.Accept(
	ctx,
	"YOUR_USER_ID",
	"YOUR_BROKER_CONNECTION_ID",
	alvaindes.AcceptBrokerAgreementParams{
		Agreement: "AGREEMENT_NAME_FROM_API",
		Revision:  &revision,
	},
)
if err != nil {
	return err
}
```

List accepted agreements:

```go
accepted, err := client.BrokerAgreements.ListAccepted(
	ctx,
	"YOUR_USER_ID",
	"YOUR_BROKER_CONNECTION_ID",
)
if err != nil {
	return err
}

for _, agreement := range accepted {
	fmt.Println(agreement.Agreement, agreement.SignedAt)
}
```

Use agreement names and revisions returned by the API. Do not assume illustrative values are current requirements.

### 13. Broker disclosures

Retrieve a broker disclosure:

```go
disclosure, err := client.BrokerDisclosures.Get(
	ctx,
	"YOUR_USER_ID",
	"YOUR_BROKER_CONNECTION_ID",
)
if err != nil {
	return err
}

if disclosure == nil {
	fmt.Println("No broker disclosure is available")
} else {
	fmt.Println("Employment status:", disclosure.EmploymentStatus)
}
```

The method returns `nil` when the broker has no disclosure.

Create a broker disclosure:

```go
disclosure, err := client.BrokerDisclosures.Create(
	ctx,
	"YOUR_USER_ID",
	"YOUR_BROKER_CONNECTION_ID",
	alvaindes.CreateBrokerDisclosureParams{
		EmploymentStatus:            "EMPLOYED",
		IsControlPerson:             false,
		IsAffiliatedExchangeOrFINRA: false,
		IsPoliticallyExposed:        false,
		ImmediateFamilyExposed:      false,
	},
)
if err != nil {
	return err
}

fmt.Println("Disclosure ID:", disclosure.ID)
```

`EmploymentSector` is optional and can be supplied using a string pointer.

### 14. Market data

Retrieve a market quote:

```go
quote, err := client.MarketData.GetQuote(ctx, "AAPL")
if err != nil {
	return err
}

fmt.Println("Symbol:", quote.Symbol)
fmt.Println("Price:", quote.Price)
fmt.Println("Currency:", quote.Currency)
fmt.Println("Timestamp:", quote.Timestamp)
```

Market-data availability depends on the broker connection and provider configuration.

## Idempotency

The SDK supports optional idempotency keys for:

- Orders
- Deposits
- Withdrawals

Provide the key through the corresponding options struct:

```go
alvaindes.CreateOrderOptions{
	IdempotencyKey: "unique-operation-key",
}
```

When supplied, the SDK sends the key using the `Idempotency-Key` HTTP header.

Recommended practices:

- Generate a unique key for each logical financial operation.
- Reuse the same key when retrying the same operation.
- Never reuse a key for a different operation or a different request payload.
- If a request times out or its outcome is uncertain, check the operation's status before submitting a new operation with a different key.

The SDK does not automatically retry financial operations. This avoids silently repeating operations when the outcome of a previous request is uncertain.

## Error handling

The SDK returns errors for unsuccessful API requests, transport failures, and other request errors.

API errors can be inspected using `AlvaindesError`.

```go
package main

import (
	"context"
	"errors"
	"fmt"

	alvaindes "github.com/davidokani/alvaindes-go-sdk"
)

func handleRequestError(err error) {
	if err == nil {
		return
	}

	var apiErr *alvaindes.AlvaindesError

	if errors.As(err, &apiErr) {
		fmt.Println("API error:", apiErr.Message)
		fmt.Println("HTTP status:", apiErr.Status)
		fmt.Println("Error code:", apiErr.Code)
		fmt.Println("Response:", apiErr.Response)
		return
	}

	fmt.Println("Request error:", err)
}

func example(ctx context.Context, client *alvaindes.Client) {
	_, err := client.Users.Get(
		ctx,
		"YOUR_APPLICATION_ID",
		"YOUR_USER_ID",
	)

	handleRequestError(err)
}
```

`AlvaindesError` exposes:

| Field      | Description                               |
| ---------- | ----------------------------------------- |
| `Status`   | HTTP status code.                         |
| `Message`  | Error message.                            |
| `Code`     | API error code, when provided.            |
| `Response` | Decoded response payload, when available. |

Network failures, timeouts, and context cancellation may not have an HTTP response status.

## Context, timeouts, and cancellation

All resource methods accept `context.Context`.

You can use a context deadline to limit the duration of an individual operation:

```go
ctx, cancel := context.WithTimeout(
	context.Background(),
	10*time.Second,
)
defer cancel()

user, err := client.Users.Get(
	ctx,
	"YOUR_APPLICATION_ID",
	"YOUR_USER_ID",
)
if err != nil {
	return err
}

fmt.Println(user.ID)
```

You can also cancel requests using a cancellable context.

The configured client timeout and context deadline both apply to requests.

## Custom API requests

The SDK provides generic request methods for public, versioned API endpoints that may not yet have dedicated resource methods.

### Request

```go
var result map[string]any

err := client.Request(
	ctx,
	"GET",
	"/v1/YOUR_PUBLIC_ENDPOINT",
	nil,
	&result,
)
if err != nil {
	return err
}
```

### RequestWithOptions

Use `RequestWithOptions` when you need to provide additional HTTP headers.

```go
headers := make(http.Header)
headers.Set("X-Custom-Header", "example-value")

var result map[string]any

err := client.RequestWithOptions(
	ctx,
	"GET",
	"/v1/YOUR_PUBLIC_ENDPOINT",
	nil,
	&result,
	alvaindes.RequestOptions{
		Headers: headers,
	},
)
if err != nil {
	return err
}
```

For requests with a JSON body, provide the request body as the fourth argument and a result pointer as the fifth argument.

These methods are intended for public, versioned customer-facing API endpoints. They do not provide access to internal, administrative, or control-plane endpoints.

## Authentication

The SDK authenticates requests using your Alvaindes API key through the `x-api-key` HTTP header.

The API key determines the environment and permissions associated with the request.

The SDK does not automatically switch between environments. Use the API key and configuration appropriate for the environment you intend to access.

Keep API keys secret and store them in environment variables or a secure secrets manager.

## API versioning

The SDK targets the public Alvaindes API under the `/v1` route.

The SDK's resource methods construct requests against the versioned API. When using the generic request methods, provide the appropriate public versioned path.

Internal, administrative, and control-plane endpoints are outside the scope of this SDK.

## Development

Clone the repository:

```bash
git clone https://github.com/davidokani/alvaindes-go-sdk.git
cd alvaindes-go-sdk
```

Run the test suite:

```bash
go test ./...
```

Run tests with the race detector:

```bash
go test -race ./...
```

Run static analysis:

```bash
go vet ./...
```

Check formatting:

```bash
gofmt -l .
```

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
