package main

import (
	"github.com/stretchr/testify/assert"
	"go-slack/testutils"
	"go-slack/testutils/testserver"
	"net/http"
	"testing"
)

var ts *testserver.TestServer

func TestMain(m *testing.M) {
	tr := testutils.TestInit()
	ts = tr.TestServer()
	tr.Run(m)
}

func TestRoot(t *testing.T) {
	respRec := ts.MakeRequest(t, "GET", "/")

	assert.Equal(t, http.StatusNotFound, respRec.Code)
}
