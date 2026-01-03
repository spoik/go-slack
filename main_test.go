package main

import (
	"go-slack/channels"
	"go-slack/database"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var mux *http.ServeMux

func TestMain(m *testing.M) {
	db := database.Connect()
	defer db.Close()

	mux = createServeMux(db)

	exitCode := m.Run()

	os.Exit(exitCode)
}

func TestRoot(t *testing.T) {
	respRec := MakeRequest(t, mux, "GET", "/")

	assert.Equal(t, http.StatusNotFound, respRec.Code)
}

func TestChannels(t *testing.T) {
	respRec := MakeRequest(t, mux, "GET", "/channels")
	var channels []channels.Channel

	DecodeJsonResponse(t, respRec, &channels)

	assert.Equal(t, http.StatusOK, respRec.Code)
	assert.Equal(t, channels[0].Id, int64(1))
	assert.Equal(t, channels[0].Name, "Main")

	assert.Equal(t, channels[1].Id, int64(2))
	assert.Equal(t, channels[1].Name, "Help")
}
