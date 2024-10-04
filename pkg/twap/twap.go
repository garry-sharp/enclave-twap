package twap

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/garry-sharp/enclave-assessment/pkg/api"
	"github.com/garry-sharp/enclave-assessment/pkg/logger"
)

func ExecuteTwap(side, amount, duration, market, interval, apiKey, apiSecret, baseURL string) error {

	// Perform initial sanity check on the input arguments
	side = strings.ToLower(side)
	err := ValidateTwapArgs(side, amount, duration, market, interval, apiKey, apiSecret, baseURL)
	if err != nil {
		return err
	}

	// Load API keys and check if user can log in with them
	err = api.Load(apiKey, apiSecret, baseURL)
	if err != nil {
		return err
	}
	timeoutCtx, cancelIsAuthed := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelIsAuthed()
	if loggedIn := api.IsLoggedIn(timeoutCtx); !loggedIn {
		return fmt.Errorf("not logged in")
	}
	logger.Info("API keys valid") //TODO what if the API keys are read only

	// Verify market exists and get the smallest increments

	timeoutCtx, cancelSpotMarketDetails := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelSpotMarketDetails()
	baseName, baseIncrement, quoteName, quoteIncrement, err := api.GetSpotMarketDetails(timeoutCtx, market)
	if err != nil {
		return err
	}
	var increment *big.Float
	if side == "buy" {
		increment = quoteIncrement
	} else {
		increment = baseIncrement
	}
	logger.Info("smallest increment for this market: ", increment)

	// Reduce quantity to the nearest increment
	quantity, okay := big.NewFloat(0).SetString(amount)
	if !okay {
		return fmt.Errorf("unable to parse amount")
	}
	q := RoundDown(quantity, increment)
	if q.String() != quantity.String() {
		logger.Info(fmt.Sprintf("minimum increment for this trading pair is %s, rounding %s amount down to %s", increment.String(), side, q.String()))
	}
	quantity = q

	// Check if user has enough balance to execute the order
	balanceAsset := baseName
	if side == "buy" {
		balanceAsset = quoteName
	}
	timeoutCtx, cancelSufficientBalance := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelSufficientBalance()
	if sufficient, err := api.SufficientSpotBalance(timeoutCtx, balanceAsset, quantity); !sufficient {
		return err
	}

	// Calculate the number of iterations and the quantities to be traded
	_duration, _ := time.ParseDuration(duration)
	_interval, _ := time.ParseDuration(interval)

	iterations := int(_duration / _interval)
	quantities, err := GetQuantities(quantity, increment, iterations)
	if err != nil {
		return err
	}

	// Create a ticker for the timer and set wait group to the number of iterations of the twap
	ticker := time.NewTicker(_interval)
	var wg sync.WaitGroup
	wg.Add(iterations)
	startTime := time.Now()

	// Create a context that can be canceled
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Use sync.Once to ensure cancellation only happens once
	var once sync.Once
	stop := atomic.Bool{}
	successfulIterations := int32(0)

	for i, qty := range quantities {

		// Check if we should stop before waiting for the next interval
		if stop.Load() {
			logger.Info(fmt.Sprintf("skipping iteration %d due to cancellation", i))
			wg.Done()
			continue
		}

		// If the first order, execute immediately; otherwise, wait for the interval or context cancellation
		if i != 0 {
			select {
			case <-ticker.C:
				go executeTrade(i, qty, &wg, ctx, cancel, &stop, &once, &successfulIterations, side, market)
			case <-ctx.Done():
				// Context was canceled while waiting for the ticker
				logger.Info(fmt.Sprintf("skipping iteration %d due to cancellation during wait", i))
				wg.Done()
				continue
			}
		} else {
			go executeTrade(i, qty, &wg, ctx, cancel, &stop, &once, &successfulIterations, side, market)
		}

	}

	wg.Wait()
	ticker.Stop()
	elapsed := time.Since(startTime)
	logger.Info(fmt.Sprintf("TWAP completed in %s", elapsed))
	logger.Info(fmt.Sprintf("completed iterations: %d", atomic.LoadInt32(&successfulIterations)))
	return nil
}

func executeTrade(i int, qty *big.Float, wg *sync.WaitGroup, ctx context.Context, cancel context.CancelFunc, stop *atomic.Bool, once *sync.Once, successfulIterations *int32, side, market string) {
	defer wg.Done()
	errorCount := 0

	for errorCount < 3 {
		// Check if the context has been canceled
		select {
		case <-ctx.Done():
			logger.Info(fmt.Sprintf("order %d aborted due to cancellation", i))
			return
		default:
		}

		if errorCount > 0 {
			time.Sleep(200 * time.Millisecond)
			logger.Info(fmt.Sprintf("retrying order, iteration %d, amount = %s", i, qty.String()))
		}

		logger.Info(fmt.Sprintf("creating order, iteration %d, amount = %s", i, qty.String()))
		response := new(api.APIResponse[api.CreateSpotOrderResponse])
		var err error
		if side == "buy" {
			err = api.NewMarketBuyOrder(ctx, market, qty, response) // Use ctx here to support cancellation
		} else {
			err = api.NewMarketSellOrder(ctx, market, qty, response) // Use ctx here to support cancellation
		}

		if err != nil {
			logger.Error(fmt.Sprintf("error creating order, iteration %d, %v", i, err))
			errorCount++
		} else {
			if response.Error != "" {
				logger.Error(fmt.Sprintf("server error creating order, iteration %d, %v", i, response.Error))
				errorCount++
			} else {
				logger.Info(fmt.Sprintf("%s order created, iteration %d, amount = %s", response.Result.OrderId, i, response.Result.Size))
				atomic.AddInt32(successfulIterations, 1)
				return
			}
		}
	}

	// If the error count exceeds the threshold, cancel all other goroutines
	once.Do(func() {
		logger.Error(fmt.Sprintf("Order %d failed 3 times, canceling all orders", i))
		stop.Store(true)
		cancel()
	})
}
