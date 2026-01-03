package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createRequest(t *testing.T, method string, url string) *http.Request {
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		t.Fatal("Creating", method, url, "request failed.")
	}

	return req
}

func MakeRequest(t *testing.T, mux *http.ServeMux, method string, url string) *httptest.ResponseRecorder {
	respRec := httptest.NewRecorder()
	req := createRequest(t, method, url)
	mux.ServeHTTP(respRec, req)
	return respRec
}

func DecodeJsonResponse(t *testing.T, respRec *httptest.ResponseRecorder, data any) {
	assert.Equal(t, "application/json", respRec.Header().Get("Content-Type"))

	err := json.Unmarshal(respRec.Body.Bytes(), data)

	assert.NoError(t, err)
}
