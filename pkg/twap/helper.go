package twap

import (
	"fmt"
	"math/big"
	"regexp"
	"strings"
	"time"
)

// simple sanity check on the input arguments. Returns an error if anything isn't supported.
func ValidateTwapArgs(side, amount, duration, market, interval, apiKey, apiSecret, baseUrl string) error {
	if strings.ToLower(side) != "buy" && strings.ToLower(side) != "sell" {
		return fmt.Errorf("side must be either buy or sell")
	}

	_duration, err := time.ParseDuration(duration)
	if err != nil {
		return fmt.Errorf("duration must be a valid time duration, received: %s", duration)
	}

	_interval, err := time.ParseDuration(interval)
	if err != nil {
		return fmt.Errorf("interval must be a valid time duration, received: %s", interval)
	}

	if _interval > _duration {
		return fmt.Errorf("interval must be less than the duration")
	}

	if _duration%_interval != 0 {
		return fmt.Errorf("interval must divide perfectly into the duration")
	}

	if _duration/_interval > 1000 {
		return fmt.Errorf("maximum of 1000 intervals is allowed per execution")
	}

	if _interval < time.Millisecond*500 {
		return fmt.Errorf("minimum interval allowed is 500ms")
	}

	_, amountParsed := big.NewFloat(0).SetString(amount)
	if !amountParsed {
		return fmt.Errorf("amount must be a valid number, received: %s", amount)
	}

	marketRegexp, _ := regexp.Compile(`^[A-Z0-9a-z]*-[A-Z0-9a-z]*$`)
	if !marketRegexp.MatchString(market) {
		return fmt.Errorf("market must be in the format BASE-QUOTE e.g AVAX-USDC")
	}

	if apiKey == "" {
		return fmt.Errorf("api-key must be provided")
	}

	if apiSecret == "" {
		return fmt.Errorf("api-secret must be provided")
	}

	if baseUrl != "https://api.enclave.market" && baseUrl != "https://api-staging.enclavemarket.dev" && baseUrl != "https://api-sandbox.enclave.market" {
		return fmt.Errorf("base-url must be one of https://api.enclave.market, https://api-staging.enclavemarket.dev, https://api-sandbox.enclave.market")
	}

	return nil
}

// Helper function to round down a value to the nearest increment
func RoundDown(value, increment *big.Float) *big.Float {
	// Divide the value by the increment
	quotient := new(big.Float).Quo(value, increment)

	// Floor the quotient by converting it to an integer
	floored, _ := quotient.Int(nil)

	// Multiply the floored quotient by the increment to get the rounded down value
	result := new(big.Float).Mul(new(big.Float).SetInt(floored), increment)

	return result
}

// Gets a spread of quantities that sum up to the given amount
func GetQuantities(amount *big.Float, increment *big.Float, segments int) ([]*big.Float, error) {
	if segments <= 0 {
		return nil, fmt.Errorf("segments must be greater than zero")
	}

	// Calculate the base quantity per segment
	baseQuantity := new(big.Float).Quo(amount, big.NewFloat(float64(segments)))
	baseQuantity = RoundDown(baseQuantity, increment)

	// Calculate the total base quantity
	totalBase := new(big.Float).Mul(baseQuantity, big.NewFloat(float64(segments)))

	// Calculate the remaining amount to distribute
	remaining := new(big.Float).Sub(amount, totalBase)

	// Distribute quantities into the slice
	quantities := make([]*big.Float, segments)
	for i := 0; i < segments; i++ {
		quantities[i] = new(big.Float).Set(baseQuantity)
	}

	// Spread the remaining amount across the segments
	for i := 0; remaining.Cmp(big.NewFloat(0)) > 0 && i < segments; i++ {
		// Add the increment to the current segment if the remaining amount is enough
		if remaining.Cmp(increment) >= 0 {
			quantities[i].Add(quantities[i], increment)
			remaining.Sub(remaining, increment)
		} else {
			// Add whatever is left in the remaining amount
			quantities[i].Add(quantities[i], remaining)
			remaining.Sub(remaining, remaining)
		}
	}

	return quantities, nil
}
