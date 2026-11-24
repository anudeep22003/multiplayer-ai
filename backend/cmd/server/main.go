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
	"github.com/anudeep22003/conman/internal/handlers"
	"github.com/anudeep22003/conman/internal/handlers/auth"
	"github.com/anudeep22003/conman/internal/handlers/health"
	"github.com/anudeep22003/conman/internal/postgres"
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

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	pool, err := postgres.NewPool(ctx, &config)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	mux := handlers.Assemble(
		handlers.RouteGroup{Prefix: "", Handler: health.NewMux(pool)},
		handlers.RouteGroup{Prefix: "/auth", Handler: auth.NewMux(pool)},
	)

	srv := server.New(&config, mux)

	log.Printf("Starting %s (%s) on %s", config.AppName, config.Environment, srv.Addr)
	if err := server.Run(ctx, srv); err != nil {
		log.Fatal(err)
	}
}
