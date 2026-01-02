package main

import (
	"go-slack/channels"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRoot(t *testing.T) {
	respRec := MakeRequest(t, "GET", "/")

	assert.Equal(t, http.StatusNotFound, respRec.Code)
}

func TestChannels(t *testing.T) {
	respRec := MakeRequest(t, "GET", "/channels")
	var channels []channels.Channel

	DecodeJsonResponse(t, respRec, &channels)

	assert.Equal(t, http.StatusOK, respRec.Code)
	assert.Equal(t, channels[0].Name, "Main")
	assert.Equal(t, channels[1].Name, "Help")
}
