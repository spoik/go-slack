package tests

import (
	"github.com/stretchr/testify/assert"
	"go-slack/channels/queries"
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
