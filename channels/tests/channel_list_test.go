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

func createChannel(t *testing.T, name string) {
	_, err := q.CreateChannel(tr.Context(), name)
	if err != nil {
		t.Fatal(err)
	}
}

func TestListChannels(t *testing.T) {
	tr.Test(func() {
		createChannel(t, "Main")
		createChannel(t, "Help")

		respRec := ts.MakeRequest(t, "GET", "/channels")
		var channels []queries.Channel

		testutils.DecodeJsonResponse(t, respRec, &channels)

		assert.Equal(t, http.StatusOK, respRec.Code)
		assert.Equal(t, channels[0].ID, int64(1))
		assert.Equal(t, channels[0].Name, "Main")

		assert.Equal(t, channels[1].ID, int64(2))
		assert.Equal(t, channels[1].Name, "Help")
	})
}
