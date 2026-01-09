package testserver

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"go-slack/httpserver"
	"net/http"
	"net/http/httptest"
	"testing"
)

type TestServer struct {
	mux *http.ServeMux
}

func New(ctx context.Context, db *pgxpool.Pool) (*TestServer, error) {
	mux := httpserver.NewMux(ctx, db)

	ts := TestServer{mux: mux}

	return &ts, nil
}

func (ts TestServer) MakeRequest(t *testing.T, method string, url string) *httptest.ResponseRecorder {
	respRec := httptest.NewRecorder()
	req := createRequest(t, method, url)
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
