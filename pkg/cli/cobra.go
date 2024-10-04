package cli

import (
	"log"
	"os"

	"github.com/garry-sharp/enclave-assessment/pkg/twap"

	"github.com/garry-sharp/enclave-assessment/pkg/logger"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

func getTwapCommand() *cobra.Command {
	var (
		side      string
		amount    string
		duration  string
		market    string
		interval  string
		apiKey    string
		apiSecret string
		baseURL   string
	)

	var twapCmd = &cobra.Command{
		Use:   "twap",
		Short: "Run a TWAP trade",
		Long: `
_______         ___    ____  
|_   _\ \      / / \  |  _ \ 
  | |  \ \ /\ / / _ \ | |_) |
  | |   \ V  V / ___ \|  __/ 
  |_|    \_/\_/_/   \_\_|`,
		Run: func(cmd *cobra.Command, args []string) {
			err := twap.ExecuteTwap(side, amount, duration, market, interval, apiKey, apiSecret, baseURL)
			if err != nil {
				logger.Error("Failed to execute TWAP trade", err)
				os.Exit(1)
			}
		},
	}

	twapCmd.Flags().StringVarP(&side, "side", "s", getEnv("TRADE_SIDE", ""), "The side the trade should run on (buy or sell)")
	twapCmd.Flags().StringVarP(&amount, "amount", "a", getEnv("AMOUNT", ""), "Amount to be bought or sold. Denominated in the quote currency if a buy and the base currency if a sell")
	twapCmd.Flags().StringVarP(&duration, "duration", "d", getEnv("DURATION", ""), "The length of time the TWAP will take place over, expressed as a number and then a unit e.g 20m for twenty minutes\nValid time units are “ns”, “us” (or “µs”), “ms”, “s”, “m”, “h”")
	twapCmd.Flags().StringVarP(&market, "market", "m", getEnv("MARKET", ""), "The market to run the trade on. Denominated in the base and quote currency separated by a hyphen e.g AVAX-USDC")
	twapCmd.Flags().StringVarP(&interval, "interval", "i", getEnv("INTERVAL", ""), "How often the TWAP will run, this must divide perfectly into the duration, expressed as a number and then a unit e.g 30s for thirty seconds\nA maximum of 1000 intervals are allowed per execution\nValid time units are “ms”, “s”, “m”, “h”, 500ms is the smallest interval")
	twapCmd.Flags().StringVar(&apiKey, "api-key", getEnv("API_KEY", ""), "The Enclave.markets API key")
	twapCmd.Flags().StringVar(&apiSecret, "api-secret", getEnv("API_SECRET", ""), "The Enclave.markets API key")
	twapCmd.Flags().StringVar(&baseURL, "base-url", getEnv("BASE_URL", "https://api-sandbox.enclave.market"), "The base url for the Enclave.markets API")
	return twapCmd
}

func LoadCLI() (*cobra.Command, error) {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using default values")
	}

	return getTwapCommand(), nil
}

// Helper function to get environment variables with a fallback default
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
