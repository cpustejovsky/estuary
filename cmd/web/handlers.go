package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/google/uuid"

	"github.com/cpustejovsky/estuary/pkg/mailer"
	"github.com/cpustejovsky/estuary/pkg/models"
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
	Token        string
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

func (app *application) sendPasswordResetEmail(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var user FormUser
	err := decoder.Decode(&user)
	if err != nil {
		fmt.Println(http.StatusBadRequest)
		app.clientError(w, http.StatusBadRequest)
		return
	}
	//the handler always reurns true unless an unexpected error occurs in which case it returns "error"
	//this prevents user enumeration
	ok, err := app.users.CheckForEmail(user.EmailAddress)
	if errors.Is(err, models.ErrNoRecord) {
		fmt.Fprint(w, true)
		return
	} else if err != nil {
		fmt.Println(http.StatusBadRequest)
		app.clientError(w, http.StatusBadRequest)
		return
	}
	if ok == true {
		//create token
		id := uuid.New()
		//insert record into PasswordResetToken
		err := app.resetTokens.Insert(id, user.EmailAddress)
		if err != nil {
			fmt.Println(err)
			app.clientError(w, http.StatusBadRequest)
			return
		}
		//then send email
		mailer.SendPasswordResetEmail(user.EmailAddress, id.String(), app.mgInstance)
	}
	fmt.Fprint(w, true)
}

func (app *application) resetPassword(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var user FormUser
	err := decoder.Decode(&user)
	if err != nil {
		fmt.Println(http.StatusBadRequest)
		app.clientError(w, http.StatusBadRequest)
		return
	}
	resetToken, err := app.resetTokens.Get(user.Token, user.EmailAddress)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) || err.Error() == fmt.Sprintf("invalid UUID length: %d", len(user.Token)) {
			fmt.Fprint(w, "no token found")
			return
		} else if err.Error() == "expired token" {
			fmt.Fprint(w, err)
			return
		} else {
			fmt.Fprint(w, err)
			return
		}
	}
	uuid, err := uuid.Parse(user.Token)
	if err != nil {
		fmt.Fprint(w, err)
		return
	}
	if resetToken.ID == uuid {
		err = app.users.UpdatePassword(user.EmailAddress, user.Password)
		if err != nil {
			fmt.Fprint(w, err)
			return
		}
		fmt.Fprint(w, "success")
	}

	//check if pw req token CreatedAt < Now - (30 minutes)
	//if, run update with hashed password
	//else, return an error that either let's the user know it's expired or that there was never a UUID issued

	//TODO:"pipe that sort of info into fail2ban and if someone gets cheeky give them a 12h IP ban." - advice from a Discord
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
