package tests

import (
	"fmt"
	"go-slack/channels/handlers"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func assertNumMessages(t *testing.T, expectedNumMessages int64) {
	numMessages, err := q.CountMessages(tr.Context())

	if err != nil {
		t.Fatal("Failed to count messages", err)
		return
	}

	assert.Equal(t, expectedNumMessages, numMessages)
}

func TestCreateMessageWithInvalidChannelId(t *testing.T) {
	tr.Test(func() {
		respRec := ts.MakeRequest(t, "POST", "/channels/invalid/messages")
		assert.Equal(t, http.StatusUnprocessableEntity, respRec.Code)
		assert.Equal(t, "Invalid channel id\n", respRec.Body.String())
	})
}

func TestCreateMessageWithUnknownChannelId(t *testing.T) {
	tr.Test(func() {
		respRec := ts.MakeRequest(t, "POST", "/channels/1/messages")
		assert.Equal(t, http.StatusNotFound, respRec.Code)
		assert.Equal(t, "Channel does not exist\n", respRec.Body.String())
	})
}

func TestCreateMessageWithInValidJSON(t *testing.T) {
	type invalidRequest struct{ name string }

	tr.Test(func() {
		assertNumMessages(t, 0)
		channel := createChannel(t, "Channel")

		data := invalidRequest{name: "test"}

		url := fmt.Sprintf("/channels/%d/messages", channel.ID)
		r := ts.MakeJsonRequest(t, "POST", url, data)

		assert.Equal(t, http.StatusBadRequest, r.Code)

		assertNumMessages(t, 0)
	})
}

func TestSuccessfulCreateMessage(t *testing.T) {
	tr.Test(func() {
		// Count the number of messages currently in the database.
		// We'll count again later to test that a new message was created as a result of the request.
		assertNumMessages(t, 0)

		// Create the channel we'll be add a message to in the request.
		channel := createChannel(t, "Channel")

		// Create data for the new message.
		newMessage := handlers.CreateMessageRequest{Message: "New Message"}

		// Make the request to create the message.
		url := fmt.Sprintf("/channels/%d/messages", channel.ID)
		respRec := ts.MakeJsonRequest(t, "POST", url, newMessage)

		// Assert the http status code is a 201.
		assert.Equal(t, http.StatusCreated, respRec.Code)

		assertNumMessages(t, 1)

		// Get the most recent message in the channel and assert it's
		// data matches the data posted in the request.
		messages, error := q.MessagesInChannel(tr.Context(), channel.ID)

		if error != nil {
			t.Fatal("Failed to fetch channel messages")
			return
		}

		message := messages[len(messages)-1]
		assert.Equal(t, newMessage.Message, message.Message)
	})
}
