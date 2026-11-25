package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/anudeep22003/multiplayer-ai/migrations"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/pressly/goose/v3"
)

const connectTimeout = 10 * time.Second

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("usage: %s <up|up-by-one|down|status|version> [args...]", os.Args[0])
	}
	command := os.Args[1]
	args := os.Args[2:]

	if err := godotenv.Load(".env.local"); err != nil && !errors.Is(err, fs.ErrNotExist) {
		log.Fatalf("load .env.local: %v", err)
	}

	directURL := os.Getenv("DIRECT_DATABASE_URL")
	if directURL == "" {
		log.Fatal("DIRECT_DATABASE_URL is required (use direct/session Postgres URL, not transaction pooler)")
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := run(ctx, directURL, command, args...); err != nil {
		log.Fatal(err)
	}
}

func validateDirectDSN(dsn string) error {
	u, err := url.Parse(dsn)
	if err != nil {
		return errors.New("invalid DIRECT_DATABASE_URL format")
	}
	if u.Scheme != "postgres" && u.Scheme != "postgresql" {
		return errors.New("DIRECT_DATABASE_URL must use postgres:// or postgresql:// scheme")
	}
	if u.Port() == "6543" {
		return errors.New("DIRECT_DATABASE_URL cannot use transaction pooler port 6543; use direct or session pooler port 5432")
	}
	if strings.EqualFold(u.Query().Get("pgbouncer"), "true") {
		return errors.New("DIRECT_DATABASE_URL cannot use pgbouncer=true transaction pooler; use direct or session pooler URL")
	}
	return nil
}

func run(ctx context.Context, dsn, command string, args ...string) error {
	if err := validateDirectDSN(dsn); err != nil {
		return err
	}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return fmt.Errorf("open db: %w", err)
	}
	defer db.Close()

	pingCtx, cancel := context.WithTimeout(ctx, connectTimeout)
	defer cancel()
	if err := db.PingContext(pingCtx); err != nil {
		return fmt.Errorf("ping db: %w", err)
	}

	goose.SetBaseFS(migrations.FS)
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("set goose dialect: %w", err)
	}

	// "." because migrations/embed.go embeds *.sql at the root of migrations.FS.
	if err := goose.RunContext(ctx, command, db, ".", args...); err != nil {
		return fmt.Errorf("goose %s: %w", command, err)
	}
	return nil
}
