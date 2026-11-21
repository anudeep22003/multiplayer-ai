package postgres

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/anudeep22003/conman/internal"
)

func TestNewPool_InvalidDSN(t *testing.T) {
	t.Parallel()

	cfg := internal.Config{
		DatabaseURL:      "://not-a-valid-dsn",
		DBMaxConns:       4,
		DBMinConns:       1,
		DBConnectTimeout: time.Second,
	}

	pool, err := NewPool(context.Background(), &cfg)
	if err == nil {
		pool.Close()
		t.Fatal("expected error for invalid DSN, got nil")
	}
}

func TestNewPool_Integration(t *testing.T) {
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("TEST_DATABASE_URL not set; skipping PostgreSQL integration test")
	}

	cfg := internal.Config{
		DatabaseURL:       dsn,
		DBMaxConns:        4,
		DBMinConns:        1,
		DBMaxConnLifetime: 5 * time.Minute,
		DBMaxConnIdleTime: time.Minute,
		DBConnectTimeout:  5 * time.Second,
	}

	pool, err := NewPool(context.Background(), &cfg)
	if err != nil {
		t.Fatalf("NewPool() unexpected error: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(context.Background()); err != nil {
		t.Fatalf("pool.Ping() unexpected error: %v", err)
	}
}
