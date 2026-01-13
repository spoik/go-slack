package tests

import (
	"go-slack/channels/queries"
	"go-slack/testutils"
	"go-slack/testutils/testrunner"
	"go-slack/testutils/testserver"
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

func createMessage(t *testing.T, channelId int64, message string) *queries.Message {
	p := queries.CreateMessageParams{ChannelID: channelId, Message: message}
	m, err := q.CreateMessage(tr.Context(), p)

	if err != nil {
		t.Fatal(err)
		return nil
	}

	return &m
}

