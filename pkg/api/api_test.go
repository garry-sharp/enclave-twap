package api

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"os"
	"testing"

	"github.com/garry-sharp/enclave-assessment/pkg/logger"

	"github.com/joho/godotenv"
)

func setup() {
	Load(os.Getenv("API_KEY"), os.Getenv("API_SECRET"), os.Getenv("BASE_URL"))
}

func TestMain(m *testing.M) {
	l, _ := logger.New()
	logger.SetLogger(l)

	err := godotenv.Load("../../.env")
	if err != nil {
		fmt.Println(err)
		return
	}

	err = Load(os.Getenv("API_KEY"), os.Getenv("API_SECRET"), os.Getenv("BASE_URL"))
	if err != nil {
		fmt.Println(err)
		return
	}

	m.Run()
}

func TestGetBalances(t *testing.T) {
	setup()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	resp := APIResponse[[]GetBalancesResponse]{}
	err := getBalances(ctx, &resp)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
	ctx = context.Background()
	err = getBalances(ctx, &resp)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
func TestGetMarkets(t *testing.T) {
	setup()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	resp := APIResponse[GetMarketsResponse]{}
	err := getMarkets(ctx, &resp)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
	ctx = context.Background()
	err = getMarkets(ctx, &resp)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
func TestGetBalance(t *testing.T) {
	setup()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	resp := APIResponse[GetBalanceResponse]{}
	err := getBalance(ctx, "AVAX", &resp)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
	ctx = context.Background()
	err = getBalance(ctx, "AVAX", &resp)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
func TestNewMarketBuyOrder(t *testing.T) {
	setup()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	resp := APIResponse[CreateSpotOrderResponse]{}
	err := NewMarketBuyOrder(ctx, "AVAX-USDC", big.NewFloat(0.01), &resp)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
	ctx = context.Background()
	err = NewMarketBuyOrder(ctx, "AVAX-USDC", big.NewFloat(0.01), &resp)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	err = NewMarketBuyOrder(ctx, "AVAX-USDQ", big.NewFloat(0.01), &resp)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
	if resp.Error == "" {
		t.Errorf("expected error, got nil")
	}
}
func TestNewMarketSellOrder(t *testing.T) {
	setup()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	resp := APIResponse[CreateSpotOrderResponse]{}
	err := NewMarketSellOrder(ctx, "AVAX-USDC", big.NewFloat(0.001), &resp)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
	ctx = context.Background()
	err = NewMarketSellOrder(ctx, "AVAX-USDC", big.NewFloat(0.001), &resp)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	err = NewMarketSellOrder(ctx, "AVAX-USDQ", big.NewFloat(0.001), &resp)
	if err == nil {
		t.Errorf("expected error, got %v", err)
	}
	if resp.Error == "" {
		t.Errorf("expected error, got nil")
	}
}
func TestIsLoggedIn(t *testing.T) {
	defer Load(os.Getenv("API_KEY"), os.Getenv("API_SECRET"), os.Getenv("BASE_URL"))
	if IsLoggedIn(context.Background()) == false {
		t.Errorf("expected true, got false")
	}
	Load("abc", "def", os.Getenv("BASE_URL"))
	if IsLoggedIn(context.Background()) == true {
		t.Errorf("expected false, got true")
	}
}

func TestGetSpotMarketDetails(t *testing.T) {
	setup()
	base, baseIncrement, quote, quoteIncrement, err := GetSpotMarketDetails(context.Background(), "AVAX-USDC")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if base != "AVAX" && quote != "USDC" && baseIncrement.String() != "0.0001" && quoteIncrement.String() != "0.01" {
		t.Errorf("unexpected values: %v, %v, %v, %v", base, baseIncrement, quote, quoteIncrement)
	}

	base, baseIncrement, quote, quoteIncrement, err = GetSpotMarketDetails(context.Background(), "AVAX-USDQ")
	if err == nil {
		t.Errorf("expected error, got nil")
	}
	if base != "" && quote != "" && baseIncrement != nil && quoteIncrement != nil {
		t.Errorf("unexpected values: %v, %v, %v, %v", base, baseIncrement, quote, quoteIncrement)
	}
}
func TestSufficientSpotBalance(t *testing.T) {
	setup()
	res1, err1 := SufficientSpotBalance(context.Background(), "AVAX", big.NewFloat(0.001))
	res2, err2 := SufficientSpotBalance(context.Background(), "AVAX", big.NewFloat(1000000))
	if err := errors.Join(err1, err2); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if res1 != true {
		t.Errorf("expected true, got: %v", res1)
	}
	if res2 != false {
		t.Errorf("expected false, got: %v", res2)
	}
}
