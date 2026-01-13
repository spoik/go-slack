package tests

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
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

func TestMessageListKnownChannel(t *testing.T) {
	tr.Test(func() {
		channel := createChannel(t, "Main")
		route := fmt.Sprintf("/channels/%d/messages", channel.ID)
		respRec := ts.MakeRequest(t, "GET", route)
		assert.Equal(t, http.StatusOK, respRec.Code)
	})
}
