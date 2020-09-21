package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

	mdb := os.Getenv("MONGODB_URL")
	client, err := mongo.NewClient(options.Client().ApplyURI(mdb))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

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
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}
