package auth

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/anudeep22003/multiplayer-ai/internal/db"
)

type UserStore interface {
	UpsertUser(ctx context.Context, username string) (db.UpsertUserRow, error)
}

type handler struct {
	store UserStore
}

func NewMux(store UserStore) *http.ServeMux {
	h := &handler{store: store}

	mux := http.NewServeMux()

	mux.HandleFunc("POST /login", h.login)
	mux.HandleFunc("OPTIONS /login", h.loginPreflight)

	return mux
}

func (h *handler) loginPreflight(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

func (h *handler) login(w http.ResponseWriter, r *http.Request) {
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

	row, err := h.store.UpsertUser(r.Context(), username)
	if err != nil {
		http.Error(w, "login failed", http.StatusInternalServerError)
		return
	}

	result := struct {
		UserID   string `json:"userId"`
		Username string `json:"username"`
		IsNew    bool   `json:"isNew"`
	}{
		UserID:   row.ID,
		Username: row.Username,
		IsNew:    row.IsNew,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
