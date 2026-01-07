package tests

import (
	"context"
	"go-slack/channels/queries"
	"go-slack/testutils"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
)

var mux *http.ServeMux

func TestMain(m *testing.M) {
	ctx := context.Background()

	var err error
	var db *pgx.Conn
	mux, db, err = testutils.TestInit(ctx, m)

	defer db.Close(ctx)

	if err != nil {
		log.Println("Failed to initialize tests:", err.Error())
		return
	}

	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestChannels(t *testing.T) {
	respRec := testutils.MakeRequest(t, mux, "GET", "/channels")
	var channels []queries.Channel

	testutils.DecodeJsonResponse(t, respRec, &channels)

	assert.Equal(t, http.StatusOK, respRec.Code)
	// assert.Equal(t, channels[0].ID, int64(1))
	// assert.Equal(t, channels[0].Name, "Main")
	//
	// assert.Equal(t, channels[1].ID, int64(2))
	// assert.Equal(t, channels[1].Name, "Help")
}
