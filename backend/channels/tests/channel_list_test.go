package tests

import (
	"github.com/stretchr/testify/assert"
	"go-slack/channels/queries"
	"go-slack/testutils"
	"go-slack/testutils/testrunner"
	"go-slack/testutils/testserver"
	"net/http"
	"testing"
)

var tr *testrunner.TestRunner
var ts *testserver.TestServer
var q *queries.Queries

func TestMain(m *testing.M) {
	tr = testutils.TestInit()
	ts = tr.TestServer()
	q = queries.New(tr.DB())
	tr.Run(m)
}

func createChannel(t *testing.T, name string) *queries.Channel {
	channel, err := q.CreateChannel(tr.Context(), name)

	if err != nil {
		t.Fatal(err)
		return nil
	}

	return &channel
}

func TestListChannels(t *testing.T) {
	tr.Test(func() {
		channel1 := createChannel(t, "Calls")
		channel2 := createChannel(t, "Apples")
		channel3 := createChannel(t, "Bananas")

		respRec := ts.MakeRequest(t, "GET", "/channels")
		var channels []queries.Channel

		testutils.DecodeJsonResponse(t, respRec, &channels)

		assert.Equal(t, http.StatusOK, respRec.Code)

		assert.Equal(t, channel2.ID, channels[0].ID)
		assert.Equal(t, channel2.Name, channels[0].Name)

		assert.Equal(t, channel3.ID, channels[1].ID)
		assert.Equal(t, channel3.Name, channels[1].Name)

		assert.Equal(t, channel1.ID, channels[2].ID)
		assert.Equal(t, channel1.Name, channels[2].Name)
	})
}
