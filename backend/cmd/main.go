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
	"github.com/anudeep22003/conman/internal/server"
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

	srv := server.New(&config, server.PassThrough)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	log.Printf("Starting %s (%s) on %s", config.AppName, config.Environment, srv.Addr)
	if err := server.Run(ctx, srv); err != nil {
		log.Fatal(err)
	}
}
