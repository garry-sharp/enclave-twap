package api

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
)

func AddAuth(req *http.Request, timestamp, method, path, body string) error {
	concattedString := timestamp + method + path + body
	mac := hmac.New(sha256.New, []byte(GetConfig().apiSecret))
	_, err := mac.Write([]byte(concattedString))
	if err != nil {
		return err
	}
	sig := mac.Sum(nil)
	req.Header.Set("ENCLAVE-KEY-ID", GetConfig().apiKey)
	req.Header.Set("ENCLAVE-TIMESTAMP", timestamp)
	req.Header.Set("ENCLAVE-SIGN", hex.EncodeToString(sig))
	return nil
}
