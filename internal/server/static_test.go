package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHTMXScriptIsServed(t *testing.T) {
	handler := New()
	request := httptest.NewRequest(http.MethodGet, "/static/js/htmx.min.js", nil)
	recorder := httptest.NewRecorder()
	
	handler.ServeHTTP(recorder, request)
	
	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Contains(t, recorder.Header().Get("Content-Type"), "javascript")
}

func TestBasePageLoadsHTMX(t *testing.T) {
	handler := New()
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	recorder := httptest.NewRecorder()
	
	handler.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Contains(t, recorder.Body.String(), `<script src="/static/js/htmx.min.js"></script>`)
}