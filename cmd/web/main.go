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
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

const (
	host   = "localhost"
	port   = 5432
	user   = "postgres"
	dbname = "estuarydev"
)

type Config struct {
	Addr string
}

type contextKey string

const contextKeyIsAuthenticated = contextKey("isAuthenticated")

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	session  *sessions.Session
	users    interface {
		Insert(string, string, string, string) error
		Authenticate(string, string) (string, error)
		Get(int) (*models.User, error)
	}
}

func main() {
	var password = os.Getenv("TEST_PSQL_PW")
	cfg := new(Config)
	flag.StringVar(&cfg.Addr, "addr", ":5000", "HTTP network address")
	flag.Parse()

	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.LUTC|log.Llongfile)
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.LUTC)

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

	session := sessions.New([]byte(os.Getenv("SESSION_SECRET")))
	session.Lifetime = 12 * time.Hour

	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
		session:  session,
		users:    &psql.UserModel{DB: db},
	}

	srv := &http.Server{
		Addr:    cfg.Addr,
		Handler: app.routes(),
	}

	infoLog.Printf("Starting server on %s", cfg.Addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}
