package main

import (
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	dynamicMiddleware := alice.New(app.session.Enable, CSRFMiddleware, app.authenticate)

	mux := pat.New()

	mux.Get("/api/test-twitter-bot", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.runTwitterBot))

	//Auth Routes
	mux.Get("/api/token", dynamicMiddleware.ThenFunc(app.getCSRFToken))
	mux.Post("/api/signup", dynamicMiddleware.ThenFunc(app.signup))
	mux.Post("/api/login", dynamicMiddleware.ThenFunc(app.login))
	mux.Get("/api/logout", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.logout))
	mux.Post("/api/password-reset", dynamicMiddleware.ThenFunc(app.sendPasswordResetEmail))
	mux.Post("/api/new-password", dynamicMiddleware.ThenFunc(app.resetPassword))

	//TODO: add Google SSO
	// mux.Get("/auth/google", dynamicMiddleware.ThenFunc(app.placeholder))
	// mux.Get("/auth/google/callback", dynamicMiddleware.ThenFunc(app.placeholder))

	//User Routes
	mux.Get("/api/user", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.getUser))
	mux.Patch("/api/user", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.updateUser))
	mux.Del("/api/user", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.deleteUser))

	//Notes Routes
	mux.Get("/api/notes/category/:name", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.getNoteByCategory))
	mux.Get("/api/notes/:id", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.placeholder))
	mux.Post("/api/notes", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.createNote))
	mux.Patch("/api/notes/:id", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.updateNote))
	mux.Patch("/api/notes/project", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.placeholder))
	mux.Patch("/api/notes/:category", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.placeholder))
	mux.Del("/api/notes/:id", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.deleteNote))

	//Project Routes
	mux.Get("/api/projects/", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.placeholder))
	mux.Get("/api/projects/done", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.placeholder))
	mux.Get("/api/projects/show/:id", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.placeholder))
	mux.Post("/api/projects", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.placeholder))
	mux.Patch("/api/projects", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.placeholder))
	mux.Patch("/api/projects/done", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.placeholder))
	mux.Del("/api/projects/:id", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.placeholder))

	return standardMiddleware.Then(mux)
}
