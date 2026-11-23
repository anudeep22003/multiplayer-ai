package main

import (
	"context"
	"errors"
	"io/fs"
	"log"

	"github.com/anudeep22003/conman/internal"
	"github.com/anudeep22003/conman/internal/postgres"
	"github.com/anudeep22003/conman/internal/server"
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

	srv := server.New(&config, pool, server.DevUserMiddleware)

	adapter := httpadapter.NewV2(srv.Handler)
	lambda.Start(adapter.ProxyWithContext)
}
