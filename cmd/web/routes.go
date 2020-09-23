package main

import (
	"net/http"

	"github.com/bmizerany/pat"
)

func (app *application) routes() http.Handler {
	mux := pat.New()

	//Auth Routes
	mux.Post("/api/signup", http.HandlerFunc(app.signup))
	mux.Post("/api/login", http.HandlerFunc(app.placeholder))
	mux.Get("/api/logout", http.HandlerFunc(app.placeholder))
	mux.Get("/auth/google", http.HandlerFunc(app.placeholder))
	mux.Get("/auth/google/callback", http.HandlerFunc(app.placeholder))

	//User Routes
	mux.Get("/api/user", http.HandlerFunc(app.getUser))
	mux.Patch("/api/user", http.HandlerFunc(app.placeholder))
	mux.Del("/api/user", http.HandlerFunc(app.placeholder))

	//Notes Routes
	mux.Get("/api/notes/category/:name", http.HandlerFunc(app.getNote))
	mux.Get("/api/notes/:id", http.HandlerFunc(app.getNote))
	mux.Post("/api/notes", http.HandlerFunc(app.placeholder))
	mux.Patch("/api/notes/:id", http.HandlerFunc(app.placeholder))
	mux.Patch("/api/notes/project", http.HandlerFunc(app.placeholder))
	mux.Patch("/api/notes/:category", http.HandlerFunc(app.placeholder))
	mux.Del("/api/notes/:id", http.HandlerFunc(app.placeholder))

	//Project Routes
	mux.Get("/api/projects/", http.HandlerFunc(app.placeholder))
	mux.Get("/api/projects/done", http.HandlerFunc(app.placeholder))
	mux.Get("/api/projects/show/:id", http.HandlerFunc(app.placeholder))
	mux.Post("/api/projects", http.HandlerFunc(app.placeholder))
	mux.Patch("/api/projects", http.HandlerFunc(app.placeholder))
	mux.Patch("/api/projects/done", http.HandlerFunc(app.placeholder))
	mux.Del("/api/projects/:id", http.HandlerFunc(app.placeholder))

	return mux
}
