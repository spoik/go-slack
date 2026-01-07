package tests

import (
	"context"
	"go-slack/channels/queries"
	"go-slack/testutils"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
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

func TestChannels(t *testing.T) {
	respRec := ts.MakeRequest(t, "GET", "/channels")
	var channels []queries.Channel

	testutils.DecodeJsonResponse(t, respRec, &channels)

	assert.Equal(t, http.StatusOK, respRec.Code)
	// assert.Equal(t, channels[0].ID, int64(1))
	// assert.Equal(t, channels[0].Name, "Main")
	//
	// assert.Equal(t, channels[1].ID, int64(2))
	// assert.Equal(t, channels[1].Name, "Help")
}
