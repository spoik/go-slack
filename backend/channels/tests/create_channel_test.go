package tests

import (
	"go-slack/channels/handlers"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func assertNumChannels(t *testing.T, expectedCount int64) {
	numChannels, err := q.CountChannels(tr.Context())
	if err != nil {
		t.Fatal("Unable to count messages", err)
		return
	}

	assert.Equal(t, expectedCount, numChannels)
}

func TestCreateChannelWithInvalidJSON(t *testing.T) {
	type invalidRequest struct{ invalidField string }
	tr.Test(func() {
		assertNumChannels(t, 0)

		data := invalidRequest{invalidField: "value"}
		r := ts.MakeJsonRequest(t, "POST", "/channels", data)

		assert.Equal(t, http.StatusBadRequest, r.Code)

		assertNumChannels(t, 0)
	})
}

func TestCreateChannelWithValidJSon(t *testing.T) {
	tr.Test(func() {
		assertNumChannels(t, 0)

		data := handlers.CreateChannelRequest{Name: "New channel"}
		r := ts.MakeJsonRequest(t, "POST", "/channels", data)
		assert.Equal(t, http.StatusOK, r.Code)

		assertNumChannels(t, 1)

		channels, err := q.ListChannels(tr.Context())

		if err != nil {
			t.Fatal("Unable to fetch channels", err)
			return
		}

		assert.Equal(t, data.Name, channels[0].Name)
	})
}

func TestCreateChannelWithDuplicateChannelName(t *testing.T) {
}
