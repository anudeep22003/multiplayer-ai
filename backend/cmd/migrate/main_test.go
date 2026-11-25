package main

import (
	"strings"
	"testing"
)

func TestValidateDirectDSN(t *testing.T) {
	tests := []struct {
		name    string
		dsn     string
		wantErr string
	}{
		{
			name: "local postgres",
			dsn:  "postgres://multiplayer_ai:multiplayer_ai@localhost:5432/multiplayer_ai?sslmode=disable",
		},
		{
			name: "supabase direct connection",
			dsn:  "postgresql://postgres:supersecret@db.example.supabase.co:5432/postgres",
		},
		{
			name: "supabase session pooler port 5432",
			dsn:  "postgresql://postgres.ref:supersecret@aws-0-us-east-1.pooler.supabase.com:5432/postgres",
		},
		{
			name: "pgbouncer false allowed",
			dsn:  "postgres://user:supersecret@localhost:5432/db?pgbouncer=false",
		},
		{
			name:    "supabase transaction pooler port 6543",
			dsn:     "postgresql://postgres.ref:supersecret@aws-0-us-east-1.pooler.supabase.com:6543/postgres",
			wantErr: "port 6543",
		},
		{
			name:    "pgbouncer true rejected",
			dsn:     "postgres://user:supersecret@aws-0-us-east-1.pooler.supabase.com:5432/postgres?pgbouncer=true",
			wantErr: "pgbouncer=true",
		},
		{
			name:    "pgbouncer TRUE rejected case-insensitively",
			dsn:     "postgres://user:supersecret@localhost:5432/db?pgbouncer=TRUE",
			wantErr: "pgbouncer=true",
		},
		{
			name:    "invalid url format does not leak credentials",
			dsn:     "://user:supersecret@localhost:5432/db",
			wantErr: "invalid DIRECT_DATABASE_URL format",
		},
		{
			name:    "non-postgres scheme rejected",
			dsn:     "mysql://user:supersecret@localhost:3306/db",
			wantErr: "postgres:// or postgresql://",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateDirectDSN(tt.dsn)
			if tt.wantErr == "" {
				if err != nil {
					t.Fatalf("validateDirectDSN() unexpected error: %v", err)
				}
				return
			}
			if err == nil {
				t.Fatalf("validateDirectDSN() = nil, want error containing %q", tt.wantErr)
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf("error = %q, want substring %q", err.Error(), tt.wantErr)
			}
			if strings.Contains(err.Error(), "supersecret") {
				t.Errorf("error leaked credentials: %q", err.Error())
			}
		})
	}
}
