package auth

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5"
)

type DB interface {
	QueryRow(context.Context, string, ...any) pgx.Row
}

type handler struct {
	db DB
}

func NewMux(db DB) *http.ServeMux {
	h := &handler{db: db}

	mux := http.NewServeMux()

	mux.HandleFunc("POST /login", h.login)
	mux.HandleFunc("OPTIONS /login", h.loginPreflight)

	return mux
}

func (h *handler) loginPreflight(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

// Finds an existing user or creates a new one in a single database call,
// and tells us whether the user was just created (is_new).
//
// How it works:
//
//  1. ON CONFLICT (username) DO UPDATE SET username = EXCLUDED.username
//     If the username already exists, we do a dummy update that sets the
//     username to the exact same value we tried to insert (EXCLUDED).
//     We do this instead of "DO NOTHING" because DO NOTHING returns 0 rows
//     on a conflict, whereas DO UPDATE forces RETURNING to hand back the existing row.
//
//  2. (xmax = 0) AS is_new
//     Postgres attaches a hidden "xmax" column to every row to track updates.
//     - Freshly inserted rows have xmax = 0        -> is_new is true
//     - Rows touched by DO UPDATE get a non-zero ID -> is_new is false
const upsertUser = `
INSERT INTO users (username) VALUES ($1)
ON CONFLICT (username) DO UPDATE SET username = EXCLUDED.username
RETURNING id, username, (xmax = 0) AS is_new
`

func (h *handler) login(w http.ResponseWriter, r *http.Request) {
	// get the req out into a json
	var input struct {
		Username string `json:"username"`
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&input); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	username := strings.TrimSpace(input.Username)
	if username == "" {
		http.Error(w, "username cannot be empty", http.StatusBadRequest)
		return
	}

	var result struct {
		UserID   string `json:"userId"`
		Username string `json:"username"`
		IsNew    bool   `json:"isNew"`
	}

	if err := h.db.QueryRow(r.Context(), upsertUser, username).Scan(&result.UserID, &result.Username, &result.IsNew); err != nil {
		http.Error(w, "login failed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(result)

}
