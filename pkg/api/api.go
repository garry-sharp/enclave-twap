package api

// The exported functions to be used

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strings"
)

func GetBalances(ctx context.Context, response *APIResponse[[]GetBalancesResponse]) error {
	return getBalances(ctx, response)
}

func GetMarkets(ctx context.Context, response *APIResponse[GetMarketsResponse]) error {
	return getMarkets(ctx, response)
}

func GetBalance(ctx context.Context, asset string, response *APIResponse[GetBalanceResponse]) error {
	return getBalance(ctx, asset, response)
}

func NewMarketBuyOrder(ctx context.Context, market string, amount *big.Float, response *APIResponse[CreateSpotOrderResponse]) error {
	body, err := json.Marshal(SpotOrderRequest{
		Market:    market,
		QuoteSize: amount.String(),
		Side:      BUY,
		Type:      MARKET,
	})
	if err != nil {
		return err
	}
	err = createSpotOrder(ctx, body, response)
	if err == nil && response.Error != "" {
		return fmt.Errorf("error creating order: %s", response.Error)
	}
	return err
}

func NewMarketSellOrder(ctx context.Context, market string, amount *big.Float, response *APIResponse[CreateSpotOrderResponse]) error {
	body, err := json.Marshal(SpotOrderRequest{
		Market: market,
		Size:   amount.String(),
		Side:   SELL,
		Type:   MARKET,
	})
	if err != nil {
		return err
	}
	err = createSpotOrder(ctx, body, response)
	if err == nil && response.Error != "" {
		return fmt.Errorf("error creating order: %s", response.Error)
	}
	return err
}

func IsLoggedIn(ctx context.Context) bool {
	response := APIResponse[string]{Result: ""}
	err := authHello(ctx, &response)
	if err != nil {
		return false
	}
	if response.Error != "" {
		return false
	}
	return true
}

// returns the base name, base increment, quote name, and quote increment. Error if market doesn't exist
func GetSpotMarketDetails(ctx context.Context, market string) (string, *big.Float, string, *big.Float, error) {
	markets := APIResponse[GetMarketsResponse]{}
	err := GetMarkets(ctx, &markets)
	if err != nil {
		return "", nil, "", nil, err
	}

	var base, quote string
	if parts := strings.Split(market, "-"); len(parts) == 2 {
		base = parts[0]
		quote = parts[1]
	} else {
		return "", nil, "", nil, fmt.Errorf("invalid market %s", market)
	}

	for _, m := range markets.Result.SpotMarkets.TradingPairs {
		if m.Pair.Base == base && m.Pair.Quote == quote {
			base := big.NewFloat(0)
			quote := big.NewFloat(0)
			_, baseParseSuccess := base.SetString(m.BaseIncrement)
			_, quoteParseSuccess := quote.SetString(m.QuoteIncrement)

			if !baseParseSuccess || !quoteParseSuccess {
				return "", nil, "", nil, errors.New("unable to parse base or quote increment")
			}
			return m.Pair.Base, base, m.Pair.Quote, quote, nil
		}
	}

	return "", nil, "", nil, fmt.Errorf("market %s does not exist", market)
}

func SufficientSpotBalance(ctx context.Context, asset string, amount *big.Float) (bool, error) {
	balance := APIResponse[GetBalanceResponse]{}
	err := GetBalance(ctx, asset, &balance)
	if err != nil {
		return false, err
	}

	if balance.Error != "" {
		return false, fmt.Errorf("error getting balance: %s", balance.Error)
	}

	balanceValue := big.NewFloat(0)
	_, okay := balanceValue.SetString(balance.Result.FreeBalance)
	if !okay {
		return false, errors.New("unable to parse balance")
	}

	if balanceValue.Cmp(amount) == -1 {
		return false, nil
	}

	return true, nil
}
