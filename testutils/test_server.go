package testutils

import (
	"context"
	"github.com/jackc/pgx/v5"
	"net/http"
	"net/http/httptest"
	"testing"
)

type TestServer struct {
	mux *http.ServeMux
	db  *pgx.Conn
}

func (ts TestServer) MakeRequest(t *testing.T, method string, url string) *httptest.ResponseRecorder {
	respRec := httptest.NewRecorder()
	req := ts.createRequest(t, method, url)
	ts.mux.ServeHTTP(respRec, req)
	return respRec
}

func (ts TestServer) createRequest(t *testing.T, method string, url string) *http.Request {
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		t.Fatal("Creating", method, url, "request failed.")
	}

	return req
}

// Should be called before TestMain returns.
func (ts TestServer) CleanUp(ctx context.Context) {
	ts.db.Close(ctx)
}
