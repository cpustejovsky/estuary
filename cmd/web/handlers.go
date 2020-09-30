package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/cpustejovsky/estuary/pkg/models"
	"github.com/cpustejovsky/estuary/pkg/mailer"
	"github.com/gorilla/csrf"
)

//TODO: remove once all handlers are in place
func (app *application) placeholder(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, World")
}

type FormUser struct {
	FirstName    string
	LastName     string
	EmailAddress string
	Password     string
	EmailUpdates bool
	AdvancedView bool
}

//Auth Routes
func (app *application) getCSRFToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("X-CSRF-Token", csrf.Token(r))
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
	app.authenticateAndReturnID(w, r, user.EmailAddress, user.Password)
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
	app.authenticateAndReturnID(w, r, user.EmailAddress, user.Password)
}

func (app *application) logout(w http.ResponseWriter, r *http.Request) {
	app.session.Remove(r, "authenticatedUserID")
	app.session.Put(r, "flash", "You've been logged out successfully")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) passwordReset(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var user FormUser
	err := decoder.Decode(&user)
	if err != nil {
		fmt.Println(http.StatusBadRequest)
		app.clientError(w, http.StatusBadRequest)
		return
	}

	ok, err := app.users.CheckForEmail(user.EmailAddress)
	if err != nil {
		fmt.Println(http.StatusBadRequest)
		app.clientError(w, http.StatusBadRequest)
		return
	}
	if ok == true {
		token := "testy-mctestface"
		mailer.SendPasswordResetEmail(user.EmailAddress, token, app.)
	}
	fmt.Fprint(w, ok)
}

//User Routes

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

func (app *application) updateUser(w http.ResponseWriter, r *http.Request) {
	uuid := app.session.GetString(r, "authenticatedUserID")

	decoder := json.NewDecoder(r.Body)

	var user FormUser
	err := decoder.Decode(&user)
	if err != nil {
		fmt.Println(http.StatusBadRequest)
		app.clientError(w, http.StatusBadRequest)
		return
	}
	u, err := app.users.Update(uuid, user.FirstName, user.LastName, user.EmailUpdates, user.AdvancedView)
	if err != nil {
		fmt.Println(err)
		app.clientError(w, http.StatusBadRequest)
		return
	}
	b, err := json.Marshal(u)
	if err != nil {
		app.serverError(w, err)
	}
	fmt.Fprint(w, string(b))
}

func (app *application) deleteUser(w http.ResponseWriter, r *http.Request) {
	uuid := app.session.GetString(r, "authenticatedUserID")
	err := app.users.Delete(uuid)
	if err != nil {
		fmt.Println(err)
		app.clientError(w, http.StatusBadRequest)
		return
	}
}

//Note Routes
func (app *application) getNote(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}
	fmt.Fprint(w, id)
}
