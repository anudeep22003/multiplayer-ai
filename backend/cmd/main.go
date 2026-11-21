package main

import (
	"errors"
	"io/fs"
	"log"
	"net/http"

	"github.com/anudeep22003/conman/internal"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env.local"); err != nil && !errors.Is(err, fs.ErrNotExist) {
		log.Fatalf("load .env.local: %v", err)
	}

	config, err := internal.Load()
	if err != nil {
		log.Fatal(err)
	}

	server := internal.NewServer(&config)
	log.Printf("Starting %s (%s) on %s", config.AppName, config.Environment, server.Addr)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("serve HTTP: %v", err)
	}
}
