package main

import (
	"os"

	"github.com/Sigdriv/paskelabyrint-api/db"
	handler "github.com/Sigdriv/paskelabyrint-api/handler"
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
		log.Fatalf("err loading << %v", err)
	}

	conn, err := db.DBConnect()
	if err != nil {
		log.Fatal("Unable to connect to database << ", err)
		os.Exit(1)
	}
	defer conn.Close()

	log.Info("Connected to database successfully")

	handler.Handler()
}
