package testutils

import (
	"github.com/jackc/pgx/v5"
	"net/http"
	"net/http/httptest"
	"context"
	"testing"
)

type TestServer struct {
	mux *http.ServeMux
	db  *pgx.Conn
	ctx context.Context
}

func (ts TestServer) MakeRequest(t *testing.T, method string, url string) *httptest.ResponseRecorder {
	respRec := httptest.NewRecorder()
	req := ts.createRequest(t, method, url)
	ts.mux.ServeHTTP(respRec, req)
	return respRec
}

func (ts TestServer)createRequest(t *testing.T, method string, url string) *http.Request {
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		t.Fatal("Creating", method, url, "request failed.")
	}

	return req
}

func (ts TestServer) CleanUp() {
	ts.db.Close(ts.ctx)
}
