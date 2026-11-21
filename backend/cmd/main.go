package main

import (
	"context"
	"errors"
	"io/fs"
	"log"
	"os"
	"os/signal"
	"syscall"

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
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	stop()

	err = internal.Run(ctx, server)
	log.Printf("Starting %s (%s) on %s", config.AppName, config.Environment, server.Addr)

	if err != nil {
		log.Fatal(err)
	}
}
