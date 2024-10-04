package api

//unexported raw http requests

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
)

func authHello(ctx context.Context, response *APIResponse[string]) error {
	path := "/authedHello"
	method := http.MethodGet
	body := ""
	timestamp := GetTimestamp()

	req, err := http.NewRequestWithContext(ctx, method, GetConfig().baseURL+path, nil)
	if err != nil {
		return err
	}

	AddAuth(req, timestamp, method, path, body)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(response)
}

func getMarkets(ctx context.Context, response *APIResponse[GetMarketsResponse]) error {
	path := "/v1/markets"
	method := http.MethodGet

	req, err := http.NewRequestWithContext(ctx, method, GetConfig().baseURL+path, nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(response)
}

func getBalances(ctx context.Context, response *APIResponse[[]GetBalancesResponse]) error {
	path := "/v0/wallet/balances"
	method := http.MethodGet
	body := ""
	timestamp := GetTimestamp()

	req, err := http.NewRequestWithContext(ctx, method, GetConfig().baseURL+path, nil)
	if err != nil {
		return err
	}

	AddAuth(req, timestamp, method, path, body)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(response)
}

func getBalance(ctx context.Context, asset string, response *APIResponse[GetBalanceResponse]) error {
	path := "/v0/get_balance"
	method := http.MethodPost
	body, err := json.Marshal(map[string]string{"symbol": asset})
	if err != nil {
		return err
	}
	timestamp := GetTimestamp()

	req, err := http.NewRequestWithContext(ctx, method, GetConfig().baseURL+path, bytes.NewBuffer([]byte(body)))

	if err != nil {
		return err
	}

	AddAuth(req, timestamp, method, path, string(body))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(response)
}

func createSpotOrder(ctx context.Context, body []byte, response *APIResponse[CreateSpotOrderResponse]) error {
	path := "/v1/orders"
	method := http.MethodPost
	timestamp := GetTimestamp()

	reqObj, err := http.NewRequestWithContext(ctx, method, GetConfig().baseURL+path, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	AddAuth(reqObj, timestamp, method, path, string(body))
	resp, err := http.DefaultClient.Do(reqObj)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(response)
}
