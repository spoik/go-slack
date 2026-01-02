package main

import (
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

	assert.Equal(t, http.StatusOK, respRec.Code)
}
