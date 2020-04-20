package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWebhookInvalidJSON(t *testing.T) {
	router := NewRouter()

	payload := `{}`

	const workflowName = "test"
	req, err := http.NewRequest("POST", APIURLJanusWebhook, strings.NewReader(payload))
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestWebhookValidJSON(t *testing.T) {
	router := NewRouter()

	payload := `[{
        "type": 1
	}, {
		"type": 2
	}]`

	const workflowName = "test"
	req, err := http.NewRequest("POST", APIURLJanusWebhook, strings.NewReader(payload))
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}
