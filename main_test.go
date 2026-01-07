package main

import (
	"context"
	"github.com/stretchr/testify/assert"
	"go-slack/testutils"
	"log"
	"net/http"
	"os"
	"testing"
)

var ts *testutils.TestServer

func TestMain(m *testing.M) {
	ctx := context.Background()

	var err error
	ts, err = testutils.TestInit(ctx)

	if err != nil {
		log.Println("Failed to initialize tests:", err.Error())
		return
	}

	defer ts.CleanUp(ctx)

	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestRoot(t *testing.T) {
	respRec := ts.MakeRequest(t, "GET", "/")

	assert.Equal(t, http.StatusNotFound, respRec.Code)
}
