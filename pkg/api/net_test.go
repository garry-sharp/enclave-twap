package api

import "testing"

func TestLoad(t *testing.T) {
	err1 := Load("", "", "")
	if err1.Error() != "variables apiKey and apiSecret must be set" {
		t.Error("Expected apiKey and apiSecret must be set")
	}

	err2 := Load("test", "password", "")
	if err2.Error() != "variable baseURL must be set" {
		t.Error("variable baseURL must be set")
	}
}
