package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/cpustejovsky/estuary/pkg/models"
	"github.com/cpustejovsky/estuary/pkg/models/psql"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "password"
	dbname   = "estuarydev"
)

type Config struct {
	Addr string
}

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	users    interface {
		Insert(string, string, string, string) error
		Authenticate(string, string) (int, error)
		Get(int) (*models.User, error)
	}
}

func main() {
	cfg := new(Config)
	flag.StringVar(&cfg.Addr, "addr", ":5000", "HTTP network address")
	flag.Parse()

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

	fmt.Println("Successfully connected!")

	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.LUTC|log.Llongfile)
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.LUTC)

	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
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
