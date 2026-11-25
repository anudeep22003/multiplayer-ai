package internal

import (
	"testing"
	"time"
)

func TestLoad(t *testing.T) {
	t.Run("missing DATABASE_URL returns error", func(t *testing.T) {
		t.Setenv("DATABASE_URL", "")

		if _, err := Load(); err == nil {
			t.Fatal("expected error when DATABASE_URL is empty, got nil")
		}
	})

	t.Run("loads required DATABASE_URL and pool defaults", func(t *testing.T) {
		wantDSN := "postgres://multiplayer_ai:multiplayer_ai@localhost:5432/multiplayer_ai?sslmode=disable"
		t.Setenv("DATABASE_URL", wantDSN)

		cfg, err := Load()
		if err != nil {
			t.Fatalf("Load() unexpected error: %v", err)
		}

		if cfg.DatabaseURL != wantDSN {
			t.Errorf("DatabaseURL = %q, want %q", cfg.DatabaseURL, wantDSN)
		}
		if cfg.DBMaxConns != 25 {
			t.Errorf("DBMaxConns = %d, want 25", cfg.DBMaxConns)
		}
		if cfg.DBMinConns != 2 {
			t.Errorf("DBMinConns = %d, want 2", cfg.DBMinConns)
		}
		if cfg.DBMaxConnLifetime != time.Hour {
			t.Errorf("DBMaxConnLifetime = %v, want 1h", cfg.DBMaxConnLifetime)
		}
		if cfg.DBMaxConnIdleTime != 30*time.Minute {
			t.Errorf("DBMaxConnIdleTime = %v, want 30m", cfg.DBMaxConnIdleTime)
		}
		if cfg.DBConnectTimeout != 5*time.Second {
			t.Errorf("DBConnectTimeout = %v, want 5s", cfg.DBConnectTimeout)
		}
	})
}
