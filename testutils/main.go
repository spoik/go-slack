package testutils

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"go-slack/database"
	"go-slack/httpserver"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestInit(ctx context.Context, m *testing.M) (*http.ServeMux, *pgx.Conn, error) {
	db, err := database.Connect(ctx)

	if err != nil {
		fmt.Printf("Failed to connect to the database: %s", err)
		return nil, nil, err
	}

	return httpserver.New(ctx, db), db, nil
}

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
