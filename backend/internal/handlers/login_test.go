package handlers

import (
	"context"
	"encoding/json"
	"github.com/jackc/pgx/v5"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/anudeep22003/conman/internal"
	"github.com/anudeep22003/conman/internal/postgres"
)

func TestInvalidLogin(t *testing.T) {
	for _, body := range []string{`{}`, `{"username":"  "}`, `{"username":1}`, `{"username":"valid"} {}`} {
		rec := httptest.NewRecorder()
		New(nil).ServeHTTP(rec, httptest.NewRequest("POST", "/login", strings.NewReader(body)))
		if rec.Code != 400 {
			t.Errorf("login(%s) status = %d, want 400", body, rec.Code)
		}
	}
}

func TestLoginIntegration(t *testing.T) {
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("TEST_DATABASE_URL not set")
	}
	cfg := internal.Config{DatabaseURL: dsn, DBMaxConns: 1, DBConnectTimeout: 5 * time.Second}
	pool, err := postgres.NewPool(t.Context(), &cfg)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	migration, err := os.ReadFile("../../migrations/001_users.sql")
	if err != nil {
		t.Fatal(err)
	}
	// Use a temporary table scoped to the test transaction.
	tx, err := pool.Begin(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback(t.Context())
	if _, err := tx.Exec(t.Context(), strings.Replace(string(migration), "CREATE TABLE", "CREATE TEMP TABLE", 1)); err != nil {
		t.Fatal(err)
	}
	handler := New(loginTestDB{tx})
	var firstID string
	for _, wantNew := range []bool{true, false} {
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, httptest.NewRequest("POST", "/login", strings.NewReader(`{"username":"  alice "}`)))
		var got struct {
			UserID   string
			Username string
			IsNew    *bool
		}
		if rec.Code != 200 {
			t.Fatalf("status = %d: %s", rec.Code, rec.Body.String())
		}
		if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
			t.Fatal(err)
		}
		if got.UserID == "" || got.Username != "alice" || got.IsNew == nil || *got.IsNew != wantNew {
			t.Fatalf("unexpected login: %s", rec.Body.String())
		}
		if !wantNew && got.UserID != firstID {
			t.Fatalf("user ID changed from %s to %s", firstID, got.UserID)
		}
		firstID = got.UserID
	}
}

type loginTestDB struct{ pgx.Tx }

func (loginTestDB) Ping(context.Context) error { return nil }
