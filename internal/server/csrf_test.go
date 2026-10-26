package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCrossSiteRequestForgery(t *testing.T) {
	tests := []struct {
		name    string
		headers map[string]string
	}{
		{
			name: "Sec-Fetch-Site says cross-site",
			headers: map[string]string{
				"Sec-Fetch-Site": "cross-site",
			},
		},
		{
			name: "no Sec-Fetch-Site, Origin does not match Host",
			headers: map[string]string{
				"Origin": "https://attacker.example.com",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := New()
			request := httptest.NewRequest(http.MethodPost, "/", nil)
			for name, value := range tt.headers {
				request.Header.Set(name, value)
			}
			recorder := httptest.NewRecorder()

			handler.ServeHTTP(recorder, request)

			assert.Equal(t, http.StatusForbidden, recorder.Code)
		})
	}
}
