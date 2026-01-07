package testutils

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"go-slack/database"
	"go-slack/httpserver"
	"net/http/httptest"
	"testing"
)

func TestInit(ctx context.Context) (*TestServer, error) {
	db, err := database.Connect(ctx)

	if err != nil {
		fmt.Printf("Failed to connect to the database: %s", err)
		return nil, err
	}

	mux := httpserver.New(ctx, db)
	ts := TestServer{
		mux: mux,
		db:  db,
		ctx: ctx,
	}
	return &ts, nil
}

func DecodeJsonResponse(t *testing.T, respRec *httptest.ResponseRecorder, data any) {
	assert.Equal(t, "application/json", respRec.Header().Get("Content-Type"))

	err := json.Unmarshal(respRec.Body.Bytes(), data)

	assert.NoError(t, err)
}
