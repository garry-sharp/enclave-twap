package api

import (
	"context"
	"encoding/json"
	"testing"
)

func TestEndpointgetMarkets(t *testing.T) {
	setup()
	markets := APIResponse[GetMarketsResponse]{}
	ctx := context.Background()
	err := getMarkets(ctx, &markets)
	if err != nil {
		t.Error(err)
	}
	if markets.Error != "" {
		t.Error(markets.Error)
	}
}

func TestEndpointgetBalances(t *testing.T) {
	setup()
	balances := APIResponse[[]GetBalancesResponse]{}
	ctx := context.Background()
	err := getBalances(ctx, &balances)
	if err != nil {
		t.Error(err)
	}
	if balances.Error != "" {
		t.Error(balances.Error)
	}
}

func TestEndpointgetBalance(t *testing.T) {
	setup()
	balance := APIResponse[GetBalanceResponse]{}
	ctx := context.Background()
	err := getBalance(ctx, "AVAX", &balance)
	if err != nil {
		t.Error(err)
	}
	if balance.Error != "" {
		t.Error(balance.Error)
	}

}

func TestEndpointcreateSpotOrder(t *testing.T) {
	setup()
	body1, _ := json.Marshal(SpotOrderRequest{
		Market:    "AVAX-USDC",
		QuoteSize: "0.01",
		Side:      "buy",
		Type:      "market",
	})

	body2, _ := json.Marshal(SpotOrderRequest{
		Market: "AVAX-USDC",
		Size:   "0.001",
		Side:   "sell",
		Type:   "market",
	})

	response := APIResponse[CreateSpotOrderResponse]{}
	ctx := context.Background()
	err := createSpotOrder(ctx, body1, &response)

	if err != nil {
		t.Error(err)
	}
	if response.Error != "" {
		t.Error(response.Error)
	}

	err = createSpotOrder(ctx, body2, &response)
	if err != nil {
		t.Error(err)
	}
	if response.Error != "" {
		t.Error(response.Error)
	}
}
