package main

import (
	"context"
	"errors"
	"io/fs"
	"log"

	"github.com/anudeep22003/conman/internal"
	"github.com/anudeep22003/conman/internal/handlers"
	"github.com/anudeep22003/conman/internal/handlers/auth"
	"github.com/anudeep22003/conman/internal/handlers/health"
	"github.com/anudeep22003/conman/internal/postgres"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env.local"); err != nil && !errors.Is(err, fs.ErrNotExist) {
		log.Fatalf("load .env.local: %v", err)
	}

	config, err := internal.Load()
	if err != nil {
		log.Fatal("invalid application configuration")
	}

	pool, err := postgres.NewPool(context.Background(), &config)
	if err != nil {
		log.Fatal("could not initialize PostgreSQL pool")
	}
	defer pool.Close()

	handler := handlers.Assemble(
		handlers.RouteGroup{Prefix: "", Handler: health.NewMux(pool)},
		handlers.RouteGroup{Prefix: "/auth", Handler: auth.NewMux(pool)},
	)

	adapter := httpadapter.NewV2(handler)
	lambda.Start(adapter.ProxyWithContext)
}
