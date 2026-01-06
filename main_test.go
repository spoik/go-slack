package main

import (
	"context"
	"fmt"
	"go-slack/channels/queries"
	"go-slack/database"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var mux *http.ServeMux

func TestMain(m *testing.M) {
	ctx := context.Background()

	db, err := database.Connect(ctx)

	if err != nil {
		fmt.Printf("Failed to connect to the database: %s", err)
		return
	}

	defer db.Close(ctx)

	mux = createServeMux(ctx, db)

	exitCode := m.Run()

	os.Exit(exitCode)
}

func TestRoot(t *testing.T) {
	respRec := MakeRequest(t, mux, "GET", "/")

	assert.Equal(t, http.StatusNotFound, respRec.Code)
}

func TestChannels(t *testing.T) {
	respRec := MakeRequest(t, mux, "GET", "/channels")
	var channels []queries.Channel

	DecodeJsonResponse(t, respRec, &channels)

	assert.Equal(t, http.StatusOK, respRec.Code)
	// assert.Equal(t, channels[0].ID, int64(1))
	// assert.Equal(t, channels[0].Name, "Main")
	//
	// assert.Equal(t, channels[1].ID, int64(2))
	// assert.Equal(t, channels[1].Name, "Help")
}
