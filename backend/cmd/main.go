package main

import (
	"errors"
	"fmt"
	"io/fs"
	"log"

	"github.com/anudeep22003/conman/internal"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env.local"); err != nil && !errors.Is(err, fs.ErrNotExist) {
		log.Fatalf("load .env: %v", err)
	}

	config, err := internal.Load()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Hello from %s (%s)\n", config.AppName, config.Environment)

}
