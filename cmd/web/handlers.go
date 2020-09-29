package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/cpustejovsky/estuary/pkg/models"
	"github.com/gorilla/csrf"
)

func (app *application) placeholder(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, World")
}

func (app *application) getNote(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}
	fmt.Fprint(w, id)
}

type User struct {
	FirstName    string
	LastName     string
	EmailAddress string
}

func (app *application) getUser(w http.ResponseWriter, r *http.Request) {
	u, err := app.users.Get(app.session.GetString(r, "authenticatedUserID"))
	if errors.Is(err, models.ErrNoRecord) || !u.Active {
		app.session.Remove(r, "authenticatedUserID")
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}
	b, err := json.Marshal(u)
	if err != nil {
		app.serverError(w, err)
	}
	fmt.Fprint(w, string(b))
}

func (app *application) getCSRFToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("X-CSRF-Token", csrf.Token(r))
}

type FormUser struct {
	FirstName    string
	LastName     string
	EmailAddress string
	Password     string
}

func (app *application) signup(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var user FormUser
	err := decoder.Decode(&user)
	if err != nil {
		fmt.Println(http.StatusBadRequest)
		app.clientError(w, http.StatusBadRequest)
		return
	}
	fmt.Println(user)
	err = app.users.Insert(user.FirstName, user.LastName, user.EmailAddress, user.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			//TODO: communicate serverside errors to front-end from golang to react
			fmt.Println("email address is already in use")
		} else {
			app.serverError(w, err)
		}
		return
	}
	app.authenticateAndRedirect(w, r, user.EmailAddress, user.Password)
}

func (app *application) login(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var user FormUser
	err := decoder.Decode(&user)
	if err != nil {
		fmt.Println(http.StatusBadRequest)
		app.clientError(w, http.StatusBadRequest)
		return
	}
	app.authenticateAndRedirect(w, r, user.EmailAddress, user.Password)
}

func (app *application) logout(w http.ResponseWriter, r *http.Request) {
	app.session.Remove(r, "authenticatedUserID")
	app.session.Put(r, "flash", "You've been logged out successfully")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
