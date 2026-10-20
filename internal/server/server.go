// Package server builds the HTTP routes for Tiebreaker.
package server

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"

	"github.com/hbabb/tiebreaker-go/internal/views"
	"github.com/hbabb/tiebreaker-go/static"
)

// New returns the app's router with every route registered.
func New() http.Handler {
	r := chi.NewRouter()

	r.Method(http.MethodGet, "/", templ.Handler(views.Home()))
	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServerFS(static.Files)))

	return r
}
