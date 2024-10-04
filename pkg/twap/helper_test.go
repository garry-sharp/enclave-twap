package twap

import (
	"math/big"
	"testing"
)

func TestValidateTwapArgs(t *testing.T) {
	tests := []struct {
		name      string
		side      string
		amount    string
		duration  string
		market    string
		interval  string
		apiKey    string
		apiSecret string
		baseUrl   string
		expectErr bool
	}{
		// Valid cases
		{
			name:      "Valid buy order",
			side:      "buy",
			amount:    "100.0",
			duration:  "10m",
			market:    "BTC-USD",
			interval:  "1m",
			apiKey:    "valid_key",
			apiSecret: "valid_secret",
			baseUrl:   "https://api.enclave.market",
			expectErr: false,
		},
		{
			name:      "Valid sell order",
			side:      "sell",
			amount:    "50.0",
			duration:  "15m",
			market:    "ETH-USDC",
			interval:  "5m",
			apiKey:    "valid_key",
			apiSecret: "valid_secret",
			baseUrl:   "https://api-sandbox.enclave.market",
			expectErr: false,
		},

		// Invalid cases
		{
			name:      "Invalid side",
			side:      "hold",
			amount:    "100.0",
			duration:  "10m",
			market:    "BTC-USD",
			interval:  "1m",
			apiKey:    "valid_key",
			apiSecret: "valid_secret",
			baseUrl:   "https://api.enclave.market",
			expectErr: true,
		},
		{
			name:      "Invalid duration",
			side:      "buy",
			amount:    "100.0",
			duration:  "invalid_duration",
			market:    "BTC-USD",
			interval:  "1m",
			apiKey:    "valid_key",
			apiSecret: "valid_secret",
			baseUrl:   "https://api.enclave.market",
			expectErr: true,
		},
		{
			name:      "Interval greater than duration",
			side:      "buy",
			amount:    "100.0",
			duration:  "5m",
			market:    "BTC-USD",
			interval:  "10m",
			apiKey:    "valid_key",
			apiSecret: "valid_secret",
			baseUrl:   "https://api.enclave.market",
			expectErr: true,
		},
		{
			name:      "Interval not dividing perfectly into duration",
			side:      "buy",
			amount:    "100.0",
			duration:  "10m",
			market:    "BTC-USD",
			interval:  "3m",
			apiKey:    "valid_key",
			apiSecret: "valid_secret",
			baseUrl:   "https://api.enclave.market",
			expectErr: true,
		},
		{
			name:      "Smaller than 500ms interval",
			side:      "buy",
			amount:    "100.0",
			duration:  "5s",
			market:    "BTC-USD",
			interval:  "100ms",
			apiKey:    "valid_key",
			apiSecret: "valid_secret",
			baseUrl:   "https://api.enclave.market",
			expectErr: true,
		},
		{
			name:      "More than 1000 intervals",
			side:      "buy",
			amount:    "100.0",
			duration:  "10m",
			market:    "BTC-USD",
			interval:  "500ms",
			apiKey:    "valid_key",
			apiSecret: "valid_secret",
			baseUrl:   "https://api.enclave.market",
			expectErr: true,
		},
		{
			name:      "Invalid amount",
			side:      "buy",
			amount:    "invalid_amount",
			duration:  "10m",
			market:    "BTC-USD",
			interval:  "1m",
			apiKey:    "valid_key",
			apiSecret: "valid_secret",
			baseUrl:   "https://api.enclave.market",
			expectErr: true,
		},
		{
			name:      "Invalid market format",
			side:      "buy",
			amount:    "100.0",
			duration:  "10m",
			market:    "BTCUSD",
			interval:  "1m",
			apiKey:    "valid_key",
			apiSecret: "valid_secret",
			baseUrl:   "https://api.enclave.market",
			expectErr: true,
		},
		{
			name:      "Missing API key",
			side:      "buy",
			amount:    "100.0",
			duration:  "10m",
			market:    "BTC-USD",
			interval:  "1m",
			apiKey:    "",
			apiSecret: "valid_secret",
			baseUrl:   "https://api.enclave.market",
			expectErr: true,
		},
		{
			name:      "Missing API secret",
			side:      "buy",
			amount:    "100.0",
			duration:  "10m",
			market:    "BTC-USD",
			interval:  "1m",
			apiKey:    "valid_key",
			apiSecret: "",
			baseUrl:   "https://api.enclave.market",
			expectErr: true,
		},
		{
			name:      "Invalid base URL",
			side:      "buy",
			amount:    "100.0",
			duration:  "10m",
			market:    "BTC-USD",
			interval:  "1m",
			apiKey:    "valid_key",
			apiSecret: "valid_secret",
			baseUrl:   "https://invalid-url.com",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateTwapArgs(tt.side, tt.amount, tt.duration, tt.market, tt.interval, tt.apiKey, tt.apiSecret, tt.baseUrl)
			if (err != nil) != tt.expectErr {
				t.Errorf("expected error: %v, got: %v", tt.expectErr, err)
			}
		})
	}
}

func TestRoundDown(t *testing.T) {
	res1 := RoundDown(big.NewFloat(1.23456), big.NewFloat(0.01))
	res2 := RoundDown(big.NewFloat(1.23456), big.NewFloat(0.1))
	res3 := RoundDown(big.NewFloat(1.23456), big.NewFloat(0.00001))

	if res1.String() != "1.23" {
		t.Errorf("expected 1.23, got: %s", res1.String())
	}

	if res2.String() != "1.2" {
		t.Errorf("expected 1.2, got: %s", res2.String())
	}

	if res3.String() != "1.23456" {
		t.Errorf("expected 1.23456, got: %s", res3.String())
	}
}

func TestGetQuantities(t *testing.T) {
	total := big.NewFloat(0)
	total.SetString("13.9171")

	increment := big.NewFloat(0)
	increment.SetString("0.01")

	total = RoundDown(total, increment)

	quantities, err := GetQuantities(total, increment, 30)
	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}

	if len(quantities) != 30 {
		t.Errorf("expected 30 quantities, got: %d", len(quantities))
	}

	checkTotal := big.NewFloat(0)
	for _, q := range quantities {
		checkTotal = checkTotal.Add(checkTotal, q)
		if q.String() != "0.46" && q.String() != "0.47" {
			t.Errorf("expected 0.4639, got: %s", q.String())
		}
	}
	if checkTotal.String() != "13.91" {
		t.Errorf("expected 13.91, got: %s", total.String())
	}
}
