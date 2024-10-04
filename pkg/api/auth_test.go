package api

import (
	"net/http"
	"testing"
)

func TestAddAuth(t *testing.T) {
	apiKey := "test"
	apiSecret := "password"
	time := "123456789"
	method := http.MethodGet
	path := "/test"

	Load(apiKey, apiSecret, "http://localhost:8080")

	req, _ := http.NewRequest(http.MethodGet, "http://localhost:8080/test", nil)
	err := AddAuth(req, time, method, path, "")
	if err != nil {
		t.Errorf("failed")
	}

	if req.Header.Get("ENCLAVE-KEY-ID") != apiKey {
		t.Errorf("expected %s, got %s", apiKey, req.Header.Get("ENCLAVE-KEY-ID"))
	}
	if req.Header.Get("ENCLAVE-TIMESTAMP") != time {
		t.Errorf("expected %s, got %s", time, req.Header.Get("ENCLAVE-TIMESTAMP"))
	}
	expectedSig := "9d2341ccc394feaceead8ff3b6a61a3ce6cf66050807a36055f793f37109aaac"
	if req.Header.Get("ENCLAVE-SIGN") != expectedSig {
		t.Errorf("expected %s, got %s", expectedSig, req.Header.Get("ENCLAVE-SIGN"))
	}
}
