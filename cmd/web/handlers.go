package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/mailgun/mailgun-go/v4"

	"github.com/cpustejovsky/estuary/pkg/mailer"
	"github.com/cpustejovsky/estuary/pkg/models"
	t "github.com/cpustejovsky/twitter_bot"
	"github.com/gorilla/csrf"
)

//TODO: remove once all handlers are in place
func (app *application) placeholder(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, World")
}

//Auth Routes
func (app *application) getCSRFToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("X-CSRF-Token", csrf.Token(r))
}

func (app *application) signup(w http.ResponseWriter, r *http.Request) {
	user, err := app.decodeUserForm(r)
	if err != nil {
		app.serverError(w, err)
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
	user, err := app.decodeUserForm(r)
	if err != nil {
		app.serverError(w, err)
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
	user, err := app.decodeUserForm(r)
	if err != nil {
		app.serverError(w, err)
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
		id := uuid.New()
		err := app.resetTokens.Insert(id, user.EmailAddress)
		if err != nil {
			fmt.Println(err)
			app.clientError(w, http.StatusBadRequest)
			return
		}
		mailer.SendPasswordResetEmail(user.EmailAddress, id.String(), app.mgInstance)
	}
	fmt.Fprint(w, true)
}

func (app *application) resetPassword(w http.ResponseWriter, r *http.Request) {
	user, err := app.decodeUserForm(r)
	if err != nil {
		app.serverError(w, err)
		return
	}
	resetToken, err := app.resetTokens.Get(user.Token, user.EmailAddress)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			//TODO:"pipe that sort of info into fail2ban and if someone gets cheeky give them a 12h IP ban." - advice from a Discord
			fmt.Fprint(w, "no token found")
		} else {
			fmt.Fprint(w, err)
		}
		return
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
	user, err := app.decodeUserForm(r)
	if err != nil {
		app.serverError(w, err)
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
func (app *application) createNote(w http.ResponseWriter, r *http.Request) {
	noteForm, err := app.decodeNoteForm(r)
	if err != nil {
		fmt.Println(err)
		app.serverError(w, err)
		return
	}
	uuid := app.session.GetString(r, "authenticatedUserID")
	note, err := app.notes.Insert(uuid, noteForm.Content)
	if err != nil {
		fmt.Println(err)
		app.serverError(w, err)
		return
	}
	b, err := json.Marshal(note)
	if err != nil {
		app.serverError(w, err)
	}
	fmt.Fprint(w, string(b))
}

func (app *application) getNoteByCategory(w http.ResponseWriter, r *http.Request) {
	uuid := app.session.GetString(r, "authenticatedUserID")
	category := r.URL.Query().Get(":name")
	if category == "" {
		app.notFound(w)
		return
	}
	n, err := app.notes.GetByCategory(uuid, category)
	if err != nil {
		fmt.Println(err)
		app.serverError(w, err)
		return
	}
	b, err := json.Marshal(n)
	if err != nil {
		app.serverError(w, err)
	}
	fmt.Fprint(w, string(b))
}

func (app *application) updateNote(w http.ResponseWriter, r *http.Request) {
	uuid := app.session.GetString(r, "authenticatedUserID")
	noteId := r.URL.Query().Get(":id")
	form, err := app.decodeNoteForm(r)
	if err != nil {
		fmt.Println(err)
		app.serverError(w, err)
		return
	}
	note, err := app.notes.Update(uuid, noteId, form.Content)
	if err != nil {
		fmt.Println(err)
		app.serverError(w, err)
		return
	}
	b, err := json.Marshal(note)
	if err != nil {
		app.serverError(w, err)
	}
	fmt.Fprint(w, string(b))
}

func (app *application) deleteNote(w http.ResponseWriter, r *http.Request) {
	noteId := r.URL.Query().Get(":id")
	err := app.notes.Delete(noteId)
	if err != nil {
		fmt.Println(err)
		app.serverError(w, err)
		return
	}
}

func (app *application) runTwitterBot(w http.ResponseWriter, r *http.Request) {
	n := []string{"FluffyHookers", "elpidophoros"}
	c := make(chan t.User)
	creds := t.TwitterCredentials{
		AccessToken:       os.Getenv("TWITTER_ACCESS_TOKEN"),
		AccessTokenSecret: os.Getenv("TWITTER_ACCESS_TOKEN_SECRET"),
		ConsumerKey:       os.Getenv("TWITTER_CONSUMER_KEY"),
		ConsumerSecret:    os.Getenv("TWITTER_CONSUMER_SECRET"),
	}
	tb, err := t.NewBot(creds)
	if err != nil {
		log.Fatal(err)
	}
	for _, name := range n {
		go tb.FindUserTweets(name, c, 5)
		tb.AddUsers(c)
	}
	mg, err := mailgun.NewMailgunFromEnv()
	if err != nil {
		errorLog.Println(err)
	}
	if err := tb.SendEmail(mg, "charles.pustejovsky@gmail.com"); err != nil {
		fmt.Fprintf(w, "No email was sent.\n%v", err)
	} else {
		fmt.Fprintf(w, "Email is being sent")
	}
}
