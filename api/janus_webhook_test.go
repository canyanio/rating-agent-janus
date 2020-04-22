package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWebhookInvalidJSON(t *testing.T) {
	const payload = `{}`

	req, err := http.NewRequest("POST", APIURLJanusWebhook, strings.NewReader(payload))
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	router := NewRouter(nil, nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestWebhookValid(t *testing.T) {
	const payload = `[{"type":1,"timestamp":1587492912377714,"session_id":3446611678650423,"handle_id":97153772170402,"opaque_id":"siptest-fHYbsADjXHzF"}]`

	req, err := http.NewRequest("POST", APIURLJanusWebhook, strings.NewReader(payload))
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	router := NewRouter(nil, nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestWebhookInvalid(t *testing.T) {
	const payload = `[{"type":"dummy"","timestamp":1587492912377714,"session_id":3446611678650423,"handle_id":97153772170402,"opaque_id":"siptest-fHYbsADjXHzF"}]`

	req, err := http.NewRequest("POST", APIURLJanusWebhook, strings.NewReader(payload))
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	router := NewRouter(nil, nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
