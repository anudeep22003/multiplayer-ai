package auth

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/anudeep22003/multiplayer-ai/internal"
	"github.com/anudeep22003/multiplayer-ai/internal/db"
	"github.com/anudeep22003/multiplayer-ai/internal/postgres"
)

type stubUserStore struct {
	fn func(ctx context.Context, username string) (db.UpsertUserRow, error)
}

func (s stubUserStore) UpsertUser(ctx context.Context, username string) (db.UpsertUserRow, error) {
	return s.fn(ctx, username)
}

func TestLogin(t *testing.T) {
	t.Run("trims username and returns upserted user", func(t *testing.T) {
		var gotUsername string
		store := stubUserStore{
			fn: func(_ context.Context, username string) (db.UpsertUserRow, error) {
				gotUsername = username
				return db.UpsertUserRow{ID: "user-123", Username: username, IsNew: true}, nil
			},
		}

		rec := httptest.NewRecorder()
		NewMux(store).ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(`{"username":"  alice  "}`)))

		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d: %s", rec.Code, http.StatusOK, rec.Body.String())
		}
		if gotUsername != "alice" {
			t.Fatalf("store username = %q, want %q", gotUsername, "alice")
		}
		var got struct {
			UserID   string `json:"userId"`
			Username string `json:"username"`
			IsNew    bool   `json:"isNew"`
		}
		if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
			t.Fatal(err)
		}
		if got.UserID != "user-123" || got.Username != "alice" || !got.IsNew {
			t.Fatalf("unexpected response: %+v", got)
		}
	})

	t.Run("store error returns 500 without leaking details", func(t *testing.T) {
		store := stubUserStore{
			fn: func(context.Context, string) (db.UpsertUserRow, error) {
				return db.UpsertUserRow{}, errors.New("db connection refused")
			},
		}

		rec := httptest.NewRecorder()
		NewMux(store).ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(`{"username":"alice"}`)))

		if rec.Code != http.StatusInternalServerError {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusInternalServerError)
		}
		if strings.Contains(rec.Body.String(), "connection refused") {
			t.Fatalf("response leaked internal error: %q", rec.Body.String())
		}
	})
}

func TestInvalidLogin(t *testing.T) {
	mux := NewMux(nil)
	for _, body := range []string{`{}`, `{"username":"  "}`, `{"username":1}`, `{"username":"valid"} {}`} {
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(body)))
		if rec.Code != http.StatusBadRequest {
			t.Errorf("login(%s) status = %d, want %d", body, rec.Code, http.StatusBadRequest)
		}
	}
}

func TestLoginMethodNotAllowed(t *testing.T) {
	rec := httptest.NewRecorder()
	NewMux(nil).ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/login", nil))
	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusMethodNotAllowed)
	}
}

func TestLoginPreflight(t *testing.T) {
	req := httptest.NewRequest(http.MethodOptions, "/login", nil)
	req.Header.Set("Origin", "http://localhost:5173")
	req.Header.Set("Access-Control-Request-Method", "POST")
	req.Header.Set("Access-Control-Request-Headers", "content-type")
	rec := httptest.NewRecorder()

	NewMux(nil).ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent || rec.Body.Len() != 0 {
		t.Fatalf("preflight returned %d with body %q, want 204 with empty body", rec.Code, rec.Body.String())
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

	migration, err := os.ReadFile("../../../migrations/001_users.sql")
	if err != nil {
		t.Fatal(err)
	}

	// Use a temporary table scoped to the test transaction.
	tx, err := pool.Begin(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback(t.Context())

	upSQL, _, _ := strings.Cut(string(migration), "-- +goose Down")
	if _, err := tx.Exec(t.Context(), strings.Replace(upSQL, "CREATE TABLE", "CREATE TEMP TABLE", 1)); err != nil {
		t.Fatal(err)
	}

	mux := NewMux(db.New(tx))
	var firstID string
	for _, wantNew := range []bool{true, false} {
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(`{"username":"  alice "}`)))
		var got struct {
			UserID   string
			Username string
			IsNew    *bool
		}
		if rec.Code != http.StatusOK {
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
