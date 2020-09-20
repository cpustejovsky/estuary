package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

type Config struct {
	Addr string
}

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

func main() {
	cfg := new(Config)
	flag.StringVar(&cfg.Addr, "addr", ":5000", "HTTP network address")
	flag.Parse()

	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.LUTC|log.Llongfile)
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.LUTC)

	handler := http.HandlerFunc(Server)

	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
	}

	srv := &http.Server{
		Addr:    cfg.Addr,
		Handler: handler,
	}

	infoLog.Printf("Starting server on %s", cfg.Addr)
	err := srv.ListenAndServe()
	errorLog.Fatal(err)
}
