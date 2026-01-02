package main

import (
	"encoding/json"
	"net/http"
	"testing"
	"go-slack/channels"

	"github.com/stretchr/testify/assert"
)

func TestRoot(t *testing.T) {
	respRec := MakeRequest(t, "GET", "/")

	assert.Equal(t, http.StatusNotFound, respRec.Code)
}

func TestChannels(t *testing.T) {
	respRec := MakeRequest(t, "GET", "/channels")

	assert.Equal(t, http.StatusOK, respRec.Code)

	var channels []channels.Channel
	err := json.Unmarshal(respRec.Body.Bytes(), &channels)

	assert.NoError(t, err)
	assert.Equal(t, channels[0].Name, "Main")
	assert.Equal(t, channels[1].Name, "Help")
}
