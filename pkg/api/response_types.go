package api

type APIResponse[T any] struct {
	Success   bool   `json:"success"`
	Result    T      `json:"result"`
	Error     string `json:"error,omitempty"`
	ErrorCode string `json:"error_code,omitempty"`
}

// region Generics
type CurrencyPair struct {
	Base  string `json:"base"`
	Quote string `json:"quote"`
}

// region Cross Markets
type V1CrossMarketResult struct {
	Pair          CurrencyPair `json:"pair"`
	Disabled      bool         `json:"disabled"`
	DecimalPlaces int          `json:"decimalPlaces"`
}

type CrossMarkets struct {
	TradingPairs []V1CrossMarketResult `json:"tradingPairs"`
}

// region Spot
type V1SpotMarketResult struct {
	Pair           CurrencyPair `json:"pair"`
	Disabled       bool         `json:"disabled"`
	BaseIncrement  string       `json:"baseIncrement"`
	QuoteIncrement string       `json:"quoteIncrement"`
}

type SpotMarkets struct {
	TradingPairs []V1SpotMarketResult `json:"tradingPairs"`
}

// region Markets

type GetMarketsResponse struct {
	CrossMarkets CrossMarkets `json:"cross"`
	SpotMarkets  SpotMarkets  `json:"spot"`
}

//region Balance

type GetBalanceResponse struct {
	AccountId       string `json:"accountId"`
	FreeBalance     string `json:"freeBalance"`
	ReservedBalance string `json:"reservedBalance"`
	Symbol          string `json:"symbol"`
	TotalBalance    string `json:"totalBalance"`
}

//region WalletBalances

type GetBalancesResponse struct {
	Coin     string `json:"coin"`
	Free     string `json:"free"`
	Reserved string `json:"reserved"`
	Total    string `json:"total"`
	UsdValue string `json:"usdValue"`
}

//region Spot

type CreateSpotOrderResponse struct {
	CanceledAt    string `json:"canceledAt,omitempty"`
	ClientOrderId string `json:"clientOrderId"`
	CreatedAt     string `json:"createdAt,omitempty"`
	Fee           string `json:"fee"`
	FilledAt      string `json:"filledAt,omitempty"`
	FilledCost    string `json:"filledCost"`
	FilledSize    string `json:"filledSize"`
	Market        string `json:"market"`
	OrderId       string `json:"orderId"`
	Price         string `json:"price"`
	Side          string `json:"side"`
	Size          string `json:"size"`
	Status        string `json:"status"`
	Type          string `json:"type"`
	TimeInForce   string `json:"timeInForce"`
	CancelReason  string `json:"cancelReason"`
}
