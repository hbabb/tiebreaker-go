package server

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAlpineScriptsAreServed(t *testing.T) {
	tests := []struct {
		name string
		path string
	}{
		{name: "Alpine AJAX", path: "/static/js/alpine-ajax.min.js"},
		{name: "Alpine.js", path: "/static/js/alpine.min.js"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := New()
			request := httptest.NewRequest(http.MethodGet, tt.path, nil)
			recorder := httptest.NewRecorder()

			handler.ServeHTTP(recorder, request)

			assert.Equal(t, http.StatusOK, recorder.Code)
			assert.Contains(t, recorder.Header().Get("Content-Type"), "javascript")
		})
	}
}

func TestBasePageLoadsAlpineAJAXBeforeAlpine(t *testing.T) {
	handler := New()
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	require.Equal(t, http.StatusOK, recorder.Code)
	body := recorder.Body.String()

	ajaxTag := scriptTag(t, body, "/static/js/alpine-ajax.min.js")
	alpineTag := scriptTag(t, body, "/static/js/alpine.min.js")

	assert.Contains(t, ajaxTag, "defer")
	assert.Contains(t, alpineTag, "defer")
	assert.Less(t, strings.Index(body, ajaxTag), strings.Index(body, alpineTag),
		"Alpine AJAX must load before Alpine.js")
}

// scriptTag returns the opening <script> tag in body whose src is src.
// It stops the test if no such tag exists.
func scriptTag(t *testing.T, body, src string) string {
	t.Helper()

	pattern := regexp.MustCompile(`<script[^>]*src="` + regexp.QuoteMeta(src) + `"[^>]*>`)
	tag := pattern.FindString(body)
	require.NotEmpty(t, tag, "no <script> tag with src %q", src)

	return tag
}
