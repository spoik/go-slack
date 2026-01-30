package tests

import (
	"encoding/json"
	"go-slack/channels/handlers"
	"go-slack/channels/queries"
	"go-slack/testutils/testserver"
	"net/http"
	"net/http/httptest"
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

		assert.Equal(t, http.StatusUnprocessableEntity, r.Code)
		assert.Equal(t, "Name is a required field.\n", r.Body.String())

		assertNumChannels(t, 0)
	})
}

func TestCreateChannelWithDuplicateChannelName(t *testing.T) {
	tr.Test(func() {
		chanName := "Name"
		q.CreateChannel(tr.Context(), chanName)
		assertNumChannels(t, 1)

		data := handlers.CreateChannelRequest{Name: chanName}
		r := ts.MakeJsonRequest(t, "POST", "/channels", data)
		assert.Equal(t, http.StatusUnprocessableEntity, r.Code)
		assert.Equal(t, r.Body.String(), "Channel name already taken.\n")

		assertNumChannels(t, 1)
	})
}

func TestCreateChannelWithValidJSON(t *testing.T) {
	tr.Test(func() {
		assertNumChannels(t, 0)

		data := handlers.CreateChannelRequest{Name: "New channel"}
		r := ts.MakeJsonRequest(t, "POST", "/channels", data)
		assert.Equal(t, http.StatusCreated, r.Code)

		assertNumChannels(t, 1)

		channels, err := q.ListChannels(tr.Context())

		if err != nil {
			t.Fatal("Unable to fetch channels", err)
			return
		}

		dbChan := channels[0]

		assert.Equal(t, data.Name, dbChan.Name)

		var respChan *queries.Channel
		err = json.NewDecoder(r.Body).Decode(&respChan)

		if err != nil {
			t.Fatal("Unable to decode response Channel JSON", err)
			return
		}

		assert.Equal(t, dbChan.ID, respChan.ID)
		assert.Equal(t, dbChan.Name, respChan.Name)
	})
}

func TestCreateChannelWithEmptyStringName(t *testing.T) {
	tr.Test(func() {
		assertNumChannels(t, 0)
		data := handlers.CreateChannelRequest{Name: "   "}

		r := ts.MakeJsonRequest(t, "POST", "/channels", data)
		assert.Equal(t, http.StatusUnprocessableEntity, r.Code)
		assert.Equal(t, "Name must not be blank.\n", r.Body.String())

		assertNumChannels(t, 0)
	})
}

func TestCreateChannelWithEmptyStringNameAndGermanLanguage(t *testing.T) {
	tr.Test(func() {
		assertNumChannels(t, 0)
		data := handlers.CreateChannelRequest{Name: "   "}

		r := httptest.NewRecorder()

		req := testserver.CreateJsonRequest(t, "POST", "/channels", data)
		req.Header.Add("Accept-Language", "de")
		tr.TestServer().Mux().ServeHTTP(r, req)

		assert.Equal(t, http.StatusUnprocessableEntity, r.Code)
		assert.Equal(t, "Der Name darf nicht leer sein.\n", r.Body.String())

		assertNumChannels(t, 0)
	})
}
