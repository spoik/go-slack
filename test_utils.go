package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func createRequest(t *testing.T, method string, url string) *http.Request {
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		t.Fatal("Creating", method, url, "request failed.")
	}

	return req
}
func MakeRequest(t *testing.T, method string, url string) *httptest.ResponseRecorder {
	mux := createServeMux()
	respRec := httptest.NewRecorder()
	req := createRequest(t, method, url)

	mux.ServeHTTP(respRec, req)

	return respRec
}
