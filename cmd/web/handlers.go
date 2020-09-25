package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/cpustejovsky/estuary/pkg/models"
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

func (app *application) getUser(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "")
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
			// form.Errors.Add("email", "Address is already in use")
			// app.render(w, r, "signup.page.tmpl", &templateData{Form: form})
			fmt.Println("email address is already in use")
		} else {
			app.serverError(w, err)
		}
		return
	}
}
