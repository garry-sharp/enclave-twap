package api

type Side string

const (
	BUY  Side = "buy"
	SELL Side = "sell"
)

type OrderType string

const (
	LIMIT  OrderType = "limit"
	MARKET OrderType = "market"
)

type TimeInForce string

const (
	GTC TimeInForce = "GTC"
	IOC TimeInForce = "IOC"
)

type SpotOrderRequest struct {
	ClientOrderId string      `json:"clientOrderId,omitempty"`
	Market        string      `json:"market"`
	Price         string      `json:"price,omitempty"`
	QuoteSize     string      `json:"quoteSize,omitempty"`
	Side          Side        `json:"side"`
	Size          string      `json:"size,omitempty"`
	Type          OrderType   `json:"type,omitempty"`
	TimeInForce   TimeInForce `json:"timeInForce,omitempty"`
	PostOnly      bool        `json:"postOnly,omitempty"`
}
