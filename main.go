package main

import (
	"context"
	"os"

	handler "github.com/Sigdriv/paskelabyrint-api/handler"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetOutput(os.Stdout)
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	log.SetLevel(log.InfoLevel)
	log.Info("Starting Paskelabyrint API server")

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("err loading: %v", err)
	}

	dbPool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Unable to create connection pool: ", err)
		os.Exit(1)
	}
	defer dbPool.Close()

	log.Info("Connected to database successfully")

	handler.Handler()
}
