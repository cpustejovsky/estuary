package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/cpustejovsky/estuary/pkg/models"
)

func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *application) isAuthenticated(r *http.Request) bool {
	isAuthenticated, ok := r.Context().Value(contextKeyIsAuthenticated).(bool)
	if !ok {
		return false
	}
	return isAuthenticated
}

func (app *application) authenticateAndReturnID(w http.ResponseWriter, r *http.Request, email, password string) {
	id, err := app.users.Authenticate(email, password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			fmt.Println("email address or password was incorrect")
		} else {
			app.serverError(w, err)
		}
		return
	}
	app.session.Put(r, "authenticatedUserID", id)
	fmt.Fprint(w, app.session.GetString(r, "authenticatedUserID"))
}

type FormUser struct {
	FirstName    string
	LastName     string
	EmailAddress string
	Password     string
	EmailUpdates bool
	AdvancedView bool
	Token        string
}

func (app *application) decodeUserForm(r *http.Request) (FormUser, error) {
	decoder := json.NewDecoder(r.Body)

	var user FormUser
	err := decoder.Decode(&user)
	if err != nil {
		return user, err
	}
	return user, nil
}

type FormNote struct {
	Content       string
	Category      string
	Tags          []string
	DueDate       time.Time
	RemindDate    time.Time
	Completed     bool
	CompletedDate time.Time
	AccountID     string
}

func (app *application) decodeNoteForm(r *http.Request) (FormNote, error) {
	decoder := json.NewDecoder(r.Body)

	var note FormNote
	err := decoder.Decode(&note)
	if err != nil {
		return note, err
	}
	return note, nil
}
