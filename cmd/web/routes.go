package main

import (
	"net/http"

	"github.com/bmizerany/pat"
)

func (app *application) routes() http.Handler {
	mux := pat.New()

	mux.Get("/", http.HandlerFunc(app.home))
	mux.Get("/api/notes/:id", http.HandlerFunc(app.getNote))
	mux.Post("/api/notes", http.HandlerFunc(app.home))

	return mux
}
