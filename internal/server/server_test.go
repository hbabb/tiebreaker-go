package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHomePageReturnsHTML(t *testing.T) {
	handler := New()
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, recorder.Code)
	}

	contentType := recorder.Header().Get("Content-Type")
	if !strings.HasPrefix(contentType, "text/html") {
		t.Errorf("Expected Content-Type to start with 'text/html', got '%s'", contentType)
	}
}
