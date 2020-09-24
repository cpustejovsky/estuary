package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Addr string
}

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	cfg := new(Config)
	flag.StringVar(&cfg.Addr, "addr", ":5000", "HTTP network address")
	flag.Parse()

	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.LUTC|log.Llongfile)
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.LUTC)

	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
	}

	srv := &http.Server{
		Addr:    cfg.Addr,
		Handler: app.routes(),
	}

	infoLog.Printf("Starting server on %s", cfg.Addr)
	err := srv.ListenAndServe()
	errorLog.Fatal(err)
}
