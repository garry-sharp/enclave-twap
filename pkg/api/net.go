package api

import (
	"fmt"
	"time"
)

type Config struct {
	apiKey    string
	apiSecret string
	baseURL   string
}

var config *Config

func Load(apiKey, apiSecret, baseURL string) error {
	if apiKey == "" || apiSecret == "" {
		return fmt.Errorf("variables apiKey and apiSecret must be set")
	}

	if baseURL == "" {
		return fmt.Errorf("variable baseURL must be set")
	}

	config = &Config{
		apiKey:    apiKey,
		apiSecret: apiSecret,
		baseURL:   baseURL,
	}

	return nil
}

func GetConfig() *Config {
	return config
}

func GetTimestamp() string {
	return fmt.Sprintf("%d", time.Now().UnixMilli())
}
