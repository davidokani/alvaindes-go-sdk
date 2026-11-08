package alvaindes

// -----------------------------------------------------------------------------
// Users
// -----------------------------------------------------------------------------

type CreateUserParams struct {
	ExternalID  *string `json:"externalId,omitempty"`
	FirstName   string  `json:"firstName"`
	LastName    string  `json:"lastName"`
	Email       string  `json:"email"`
	PhoneNumber string  `json:"phoneNumber"`
	Country     string  `json:"country"`
}

type UserResponse struct {
	ID          string  `json:"id"`
	ExternalID  *string `json:"externalId"`
	FirstName   string  `json:"firstName"`
	LastName    string  `json:"lastName"`
	Email       string  `json:"email"`
	PhoneNumber string  `json:"phoneNumber"`
	Status      string  `json:"status"`
	CreatedAt   string  `json:"createdAt"`
}

// -----------------------------------------------------------------------------
// Accounts
// -----------------------------------------------------------------------------

type CreateInvestmentAccountParams struct {
	Currency           string `json:"currency"`
	BrokerConnectionID string `json:"brokerConnectionId"`
}

type InvestmentAccountResponse struct {
	ID            string `json:"id"`
	UserID        string `json:"userId"`
	AccountNumber string `json:"accountNumber"`
	Status        string `json:"status"`
	Currency      string `json:"currency"`
	CreatedAt     string `json:"createdAt"`
}

// -----------------------------------------------------------------------------
// KYC
// -----------------------------------------------------------------------------

type KycDocumentType string

const (
	KycDocumentPassport       KycDocumentType = "PASSPORT"
	KycDocumentNationalID     KycDocumentType = "NATIONAL_ID"
	KycDocumentDriversLicense KycDocumentType = "DRIVERS_LICENSE"
	KycDocumentVotersCard     KycDocumentType = "VOTERS_CARD"
)

type SubmitKycParams struct {
	CountryOfResidence string          `json:"countryOfResidence"`
	Nationality        *string         `json:"nationality,omitempty"`
	DateOfBirth        string          `json:"dateOfBirth"`
	AddressLine1       string          `json:"addressLine1"`
	AddressLine2       *string         `json:"addressLine2,omitempty"`
	AddressCity        string          `json:"addressCity"`
	AddressState       *string         `json:"addressState,omitempty"`
	AddressPostalCode  *string         `json:"addressPostalCode,omitempty"`
	DocumentType       KycDocumentType `json:"documentType"`
	DocumentNumber     *string         `json:"documentNumber,omitempty"`
	DocumentURL        *string         `json:"documentUrl,omitempty"`
	DocumentExpiryDate *string         `json:"documentExpiryDate,omitempty"`
	DocumentCountry    *string         `json:"documentCountry,omitempty"`
}

// UpdateKycParams represents a partial KYC update.
// Nil fields are omitted from the JSON request.
type UpdateKycParams struct {
	CountryOfResidence *string          `json:"countryOfResidence,omitempty"`
	Nationality        *string          `json:"nationality,omitempty"`
	DateOfBirth        *string          `json:"dateOfBirth,omitempty"`
	AddressLine1       *string          `json:"addressLine1,omitempty"`
	AddressLine2       *string          `json:"addressLine2,omitempty"`
	AddressCity        *string          `json:"addressCity,omitempty"`
	AddressState       *string          `json:"addressState,omitempty"`
	AddressPostalCode  *string          `json:"addressPostalCode,omitempty"`
	DocumentType       *KycDocumentType `json:"documentType,omitempty"`
	DocumentNumber     *string          `json:"documentNumber,omitempty"`
	DocumentURL        *string          `json:"documentUrl,omitempty"`
	DocumentExpiryDate *string          `json:"documentExpiryDate,omitempty"`
	DocumentCountry    *string          `json:"documentCountry,omitempty"`
}

type KycResponse struct {
	ID                 string           `json:"id"`
	UserID             string           `json:"userId"`
	Status             string           `json:"status"`
	CountryOfResidence *string          `json:"countryOfResidence"`
	Nationality        *string          `json:"nationality"`
	DateOfBirth        *string          `json:"dateOfBirth"`
	AddressLine1       *string          `json:"addressLine1"`
	AddressLine2       *string          `json:"addressLine2"`
	AddressCity        *string          `json:"addressCity"`
	AddressState       *string          `json:"addressState"`
	AddressPostalCode  *string          `json:"addressPostalCode"`
	DocumentType       *KycDocumentType `json:"documentType"`
	DocumentNumber     *string          `json:"documentNumber"`
	DocumentURL        *string          `json:"documentUrl"`
	DocumentExpiryDate *string          `json:"documentExpiryDate"`
	DocumentCountry    *string          `json:"documentCountry"`
	ApprovedAt         *string          `json:"approvedAt"`
	RejectedReason     *string          `json:"rejectedReason"`
	CreatedAt          string           `json:"createdAt"`
	UpdatedAt          string           `json:"updatedAt"`
}

// -----------------------------------------------------------------------------
// Assets
// -----------------------------------------------------------------------------

type AssetResponse struct {
	Symbol   string `json:"symbol"`
	Name     string `json:"name"`
	Exchange string `json:"exchange"`
	Currency string `json:"currency"`
	Type     string `json:"type"`
	Status   string `json:"status"`
}

type ListAssetsParams struct {
	ConnectionID string `json:"connectionId"`
	Query        string `json:"query"`
}

// -----------------------------------------------------------------------------
// Funding sources
// -----------------------------------------------------------------------------

type FundingSourceType string

const (
	FundingSourceACH         FundingSourceType = "ACH"
	FundingSourceBankAccount FundingSourceType = "BANK_ACCOUNT"
	FundingSourceWire        FundingSourceType = "WIRE"
	FundingSourceOther       FundingSourceType = "OTHER"
)

type FundingSourceStatus string

const (
	FundingSourcePending  FundingSourceStatus = "PENDING"
	FundingSourceActive   FundingSourceStatus = "ACTIVE"
	FundingSourceDisabled FundingSourceStatus = "DISABLED"
)

type CreateFundingSourceParams struct {
	Type              FundingSourceType `json:"type"`
	OwnerName         string            `json:"ownerName"`
	BankAccountNumber string            `json:"bankAccountNumber"`
	BankAccountType   string            `json:"bankAccountType"`
	RoutingNumber     string            `json:"routingNumber"`
	Currency          string            `json:"currency"`
}

type FundingSourceResponse struct {
	ID                  string              `json:"id"`
	InvestmentAccountID string              `json:"investmentAccountId"`
	Type                FundingSourceType   `json:"type"`
	Status              FundingSourceStatus `json:"status"`
	Currency            string              `json:"currency"`
	BrokerReference     *string             `json:"brokerReference"`
	CreatedAt           string              `json:"createdAt"`
	UpdatedAt           string              `json:"updatedAt"`
}

// -----------------------------------------------------------------------------
// Portfolios
// -----------------------------------------------------------------------------

type PortfolioResponse struct {
	ID                  string  `json:"id"`
	InvestmentAccountID string  `json:"investmentAccountId"`
	AccountValue        float64 `json:"accountValue"`
	CashBalance         float64 `json:"cashBalance"`
	Currency            string  `json:"currency"`
	Status              string  `json:"status"`
	CreatedAt           string  `json:"createdAt"`
}

// -----------------------------------------------------------------------------
// Holdings
// -----------------------------------------------------------------------------

type HoldingResponse struct {
	ID           string  `json:"id"`
	Symbol       string  `json:"symbol"`
	Quantity     float64 `json:"quantity"`
	AveragePrice float64 `json:"averagePrice"`
	MarketValue  float64 `json:"marketValue"`
}

// -----------------------------------------------------------------------------
// Orders
// -----------------------------------------------------------------------------

type OrderSide string

const (
	OrderSideBuy  OrderSide = "BUY"
	OrderSideSell OrderSide = "SELL"
)

type OrderType string

const (
	OrderTypeMarket OrderType = "MARKET"
	OrderTypeLimit  OrderType = "LIMIT"
)

type CreateOrderParams struct {
	Symbol     string    `json:"symbol"`
	Side       OrderSide `json:"side"`
	Quantity   float64   `json:"quantity"`
	Type       OrderType `json:"type"`
	LimitPrice *float64  `json:"limitPrice,omitempty"`
}

type CreateOrderOptions struct {
	IdempotencyKey string
}

type OrderResponse struct {
	ID            string  `json:"id"`
	Symbol        string  `json:"symbol"`
	Side          string  `json:"side"`
	Quantity      float64 `json:"quantity"`
	Status        string  `json:"status"`
	BrokerOrderID *string `json:"brokerOrderId"`
	CreatedAt     string  `json:"createdAt"`
}

// -----------------------------------------------------------------------------
// Deposits
// -----------------------------------------------------------------------------

type CreateDepositParams struct {
	Amount          float64 `json:"amount"`
	Currency        string  `json:"currency"`
	FundingSourceID *string `json:"fundingSourceId,omitempty"`
}

type CreateDepositOptions struct {
	IdempotencyKey string
}

type DepositResponse struct {
	ID        string  `json:"id"`
	Amount    float64 `json:"amount"`
	Currency  string  `json:"currency"`
	Status    string  `json:"status"`
	Reference string  `json:"reference"`
	CreatedAt string  `json:"createdAt"`
}

// -----------------------------------------------------------------------------
// Withdrawals
// -----------------------------------------------------------------------------

type CreateWithdrawalParams struct {
	Amount          float64 `json:"amount"`
	Currency        string  `json:"currency"`
	FundingSourceID *string `json:"fundingSourceId,omitempty"`
}

type CreateWithdrawalOptions struct {
	IdempotencyKey string
}

type WithdrawalResponse struct {
	ID        string  `json:"id"`
	Amount    float64 `json:"amount"`
	Currency  string  `json:"currency"`
	Status    string  `json:"status"`
	Reference string  `json:"reference"`
	CreatedAt string  `json:"createdAt"`
}

// -----------------------------------------------------------------------------
// Transactions
// -----------------------------------------------------------------------------

type TransactionResponse struct {
	ID          string  `json:"id"`
	Type        string  `json:"type"`
	Amount      float64 `json:"amount"`
	Currency    string  `json:"currency"`
	Description *string `json:"description,omitempty"`
	Status      string  `json:"status"`
	CreatedAt   string  `json:"createdAt"`
}

// -----------------------------------------------------------------------------
// Broker agreements
// -----------------------------------------------------------------------------

type BrokerAgreementRequirement struct {
	Agreement string  `json:"agreement"`
	Revision  *string `json:"revision,omitempty"`
}

type AcceptBrokerAgreementParams struct {
	Agreement string  `json:"agreement"`
	Revision  *string `json:"revision,omitempty"`
}

type BrokerAgreementResponse struct {
	ID                 string  `json:"id"`
	UserID             string  `json:"userId"`
	BrokerConnectionID string  `json:"brokerConnectionId"`
	Agreement          string  `json:"agreement"`
	Revision           *string `json:"revision"`
	SignedAt           string  `json:"signedAt"`
	IPAddress          string  `json:"ipAddress"`
	CreatedAt          string  `json:"createdAt"`
}

// -----------------------------------------------------------------------------
// Broker disclosures
// -----------------------------------------------------------------------------

type CreateBrokerDisclosureParams struct {
	EmploymentStatus            string  `json:"employmentStatus"`
	EmploymentSector            *string `json:"employmentSector,omitempty"`
	IsControlPerson             bool    `json:"isControlPerson"`
	IsAffiliatedExchangeOrFINRA bool    `json:"isAffiliatedExchangeOrFinra"`
	IsPoliticallyExposed        bool    `json:"isPoliticallyExposed"`
	ImmediateFamilyExposed      bool    `json:"immediateFamilyExposed"`
}

type BrokerDisclosureResponse struct {
	ID                          string  `json:"id"`
	UserID                      string  `json:"userId"`
	BrokerConnectionID          string  `json:"brokerConnectionId"`
	EmploymentStatus            string  `json:"employmentStatus"`
	EmploymentSector            *string `json:"employmentSector"`
	IsControlPerson             bool    `json:"isControlPerson"`
	IsAffiliatedExchangeOrFINRA bool    `json:"isAffiliatedExchangeOrFinra"`
	IsPoliticallyExposed        bool    `json:"isPoliticallyExposed"`
	ImmediateFamilyExposed      bool    `json:"immediateFamilyExposed"`
	CreatedAt                   string  `json:"createdAt"`
	UpdatedAt                   string  `json:"updatedAt"`
}

// -----------------------------------------------------------------------------
// Market data
// -----------------------------------------------------------------------------

type MarketQuoteResponse struct {
	Symbol    string  `json:"symbol"`
	Price     float64 `json:"price"`
	Currency  string  `json:"currency"`
	Timestamp string  `json:"timestamp"`
}
