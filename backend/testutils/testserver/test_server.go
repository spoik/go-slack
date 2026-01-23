package testserver

import (
	"bytes"
	"context"
	"encoding/json"
	"go-slack/httpserver"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TestServer struct {
	mux *http.ServeMux
}

func New(ctx context.Context, db *pgxpool.Pool) *TestServer {
	mux := httpserver.NewMux(ctx, db)
	ts := TestServer{mux: mux}
	return &ts
}

func (ts TestServer) MakeRequest(t *testing.T, method string, url string) *httptest.ResponseRecorder {
	respRec := httptest.NewRecorder()
	req := createRequest(t, method, url)
	ts.mux.ServeHTTP(respRec, req)
	return respRec
}

func (ts TestServer) MakeJsonRequest(t *testing.T, method string, url string, data any) *httptest.ResponseRecorder {
	respRec := httptest.NewRecorder()
	req := createJsonRequest(t, method, url, data)
	ts.mux.ServeHTTP(respRec, req)
	return respRec
}

func createRequest(t *testing.T, method string, url string) *http.Request {
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		t.Fatal("Creating", method, url, "request failed.")
	}

	return req
}

func createJsonRequest(t *testing.T, method string, url string, data any) *http.Request {
	jsonData, err := json.Marshal(data)
	if err != nil {
		t.Fatal("Unable to marshal data to JSON for", method, "to", url)
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonData))

	if err != nil {
		t.Fatal("Creating", method, url, "request failed.")
	}

	req.Header.Set("Content-Type", "application/json")

	return req
}
