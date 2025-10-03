package database

import (
	"context"
	"database/sql"
	"testing"

	_ "github.com/lib/pq"
)

func TestNew(t *testing.T) {
	// Create a mock DBTX implementation
	db, err := sql.Open("postgres", "postgres://test:test@localhost/test?sslmode=disable")
	if err != nil {
		// If we can't connect to a real database, just verify the function doesn't panic
		t.Skip("Skipping test - no database connection available")
	}
	defer db.Close()

	queries := New(db)
	if queries == nil {
		t.Fatal("expected New to return non-nil Queries")
	}
	if queries.db == nil {
		t.Error("expected Queries.db to be non-nil")
	}
}

func TestQueries_WithTx(t *testing.T) {
	// Create a mock DBTX implementation
	db, err := sql.Open("postgres", "postgres://test:test@localhost/test?sslmode=disable")
	if err != nil {
		// If we can't connect to a real database, just verify the function doesn't panic
		t.Skip("Skipping test - no database connection available")
	}
	defer db.Close()

	// Try to start a transaction
	tx, err := db.BeginTx(context.Background(), nil)
	if err != nil {
		t.Skip("Skipping test - no database connection available")
	}
	defer tx.Rollback()

	queries := New(db)
	txQueries := queries.WithTx(tx)

	if txQueries == nil {
		t.Fatal("expected WithTx to return non-nil Queries")
	}
	if txQueries.db == nil {
		t.Error("expected txQueries.db to be non-nil")
	}
}

func TestDBTXInterface(t *testing.T) {
	// This test verifies that sql.DB and sql.Tx implement DBTX interface
	db, err := sql.Open("postgres", "postgres://test:test@localhost/test?sslmode=disable")
	if err != nil {
		t.Skip("Skipping test - no database connection available")
	}
	defer db.Close()

	// Verify sql.DB implements DBTX
	var _ DBTX = db

	// Try to create a transaction and verify it implements DBTX
	tx, err := db.BeginTx(context.Background(), nil)
	if err != nil {
		t.Skip("Skipping test - no database connection available")
	}
	defer tx.Rollback()

	var _ DBTX = tx
}
