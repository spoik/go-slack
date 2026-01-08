package testutils

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"go-slack/testutils/testrunner"
	"log"
	"net/http/httptest"
	"os"
	"testing"
)

func TestInit() *testrunner.TestRunner {
	tr, err := testrunner.New()

	if err != nil {
		log.Println("Failed to initialize tests:", err.Error())
		os.Exit(1)
	}

	return tr
}

func DecodeJsonResponse(t *testing.T, respRec *httptest.ResponseRecorder, data any) {
	assert.Equal(t, "application/json", respRec.Header().Get("Content-Type"))

	err := json.Unmarshal(respRec.Body.Bytes(), data)

	assert.NoError(t, err)
}
