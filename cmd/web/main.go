package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/cpustejovsky/estuary/pkg/models"
	"github.com/cpustejovsky/estuary/pkg/models/psql"
	"github.com/golangcollege/sessions"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/mailgun/mailgun-go/v4"
)

const (
	host   = "localhost"
	port   = 5432
	user   = "postgres"
	dbname = "estuarydev"
)

var (
	errorLog *log.Logger
	infoLog  *log.Logger
)

type Config struct {
	Addr string
}

type contextKey string

const contextKeyIsAuthenticated = contextKey("isAuthenticated")

type application struct {
	errorLog    *log.Logger
	infoLog     *log.Logger
	mgInstance  *mailgun.MailgunImpl
	session     *sessions.Session
	resetTokens interface {
		Insert(id uuid.UUID, email string) error
		Get(string, string) (*models.ResetToken, error)
	}
	users interface {
		Authenticate(string, string) (string, error)
		CheckForEmail(string) (bool, error)
		Delete(string) error
		Get(string) (*models.Account, error)
		Insert(string, string, string, string) error
		Update(string, string, string, bool, bool) (*models.Account, error)
		UpdatePassword(string, string) error
	}
	notes interface {
		Insert(string, string) (*models.Note, error)
		GetByCategory(string, string) (*[]models.Note, error)
		Update(string, string, string) (*models.Note, error)
		Delete(string) error
	}
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}

	// Logging
	errorLog = log.New(os.Stderr, "ERROR\t", log.Ldate|log.LUTC|log.Llongfile)
	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.LUTC)
}

func main() {
	// Flag and Config Setup
	cfg := new(Config)
	flag.StringVar(&cfg.Addr, "addr", ":5000", "HTTP network address")
	flag.Parse()

	// Environemntal Variables
	var password = os.Getenv("TEST_PSQL_PW")
	var sessionSecret = []byte(os.Getenv("SESSION_SECRET"))
	fmt.Println("password: ", password)
	// DB Setup
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	infoLog.Println("Successfully connected to database!")

	// Session Setup
	session := sessions.New(sessionSecret)
	session.Lifetime = 12 * time.Hour

	//MailGun Set Up
	mg, err := mailgun.NewMailgunFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	mgInstance := mg

	// Application and Server Initialization
	app := &application{
		errorLog:    errorLog,
		infoLog:     infoLog,
		mgInstance:  mgInstance,
		session:     session,
		notes:       &psql.NoteModel{DB: db},
		resetTokens: &psql.ResetTokenModel{DB: db},
		users:       &psql.UserModel{DB: db},
	}

	srv := &http.Server{
		Addr:    cfg.Addr,
		Handler: app.routes(),
	}
	infoLog.Printf("Starting server on %s", cfg.Addr)

	// Server Start
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}
