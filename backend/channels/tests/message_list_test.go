package tests

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"go-slack/channels/queries"
	"go-slack/testutils"
	"net/http"
	"testing"
)

func TestMessageListUnkownChannel(t *testing.T) {
	tr.Test(func() {
		respRec := ts.MakeRequest(t, "GET", "/channels/1/messages")
		assert.Equal(t, http.StatusNotFound, respRec.Code)
		assert.Equal(t, "Channel does not exist\n", respRec.Body.String())
	})
}

func TestMessageListInvalidChannelId(t *testing.T) {
	tr.Test(func() {
		respRec := ts.MakeRequest(t, "GET", "/channels/invalid/messages")
		assert.Equal(t, http.StatusUnprocessableEntity, respRec.Code)
		assert.Equal(t, "Invalid channel id\n", respRec.Body.String())
	})
}

func TestMessageListKnownChannelWithMessages(t *testing.T) {
	tr.Test(func() {
		channel := createChannel(t, "Main")
		message1 := createMessage(t, channel.ID, "Message")
		message2 := createMessage(t, channel.ID, "Message 2")

		otherChannel := createChannel(t, "Secondary")
		createMessage(t, otherChannel.ID, "Other Message")

		route := fmt.Sprintf("/channels/%d/messages", channel.ID)
		respRec := ts.MakeRequest(t, "GET", route)

		var messages []queries.Message
		testutils.DecodeJsonResponse(t, respRec, &messages)

		assert.Equal(t, http.StatusOK, respRec.Code)

		assert.Equal(t, message1.ID, messages[0].ID)
		assert.Equal(t, message1.Message, messages[0].Message)

		assert.Equal(t, message2.ID, messages[1].ID)
		assert.Equal(t, message2.Message, messages[1].Message)
	})
}

func TestMessageListKnownChannelWithoutMessages(t *testing.T) {
	tr.Test(func() {
		channel := createChannel(t, "Main")

		route := fmt.Sprintf("/channels/%d/messages", channel.ID)
		respRec := ts.MakeRequest(t, "GET", route)

		var messages []queries.Message
		testutils.DecodeJsonResponse(t, respRec, &messages)

		assert.Equal(t, http.StatusOK, respRec.Code)

		assert.Equal(t, "[]\n", respRec.Body.String())
	})
}
