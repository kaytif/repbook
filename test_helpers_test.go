package main

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func executeRequest(
	t *testing.T,
	handler http.HandlerFunc,
	method, path, body string,
) *httptest.ResponseRecorder {
	t.Helper()

	req := httptest.NewRequest(
		method,
		path,
		strings.NewReader(body),
	)

	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	handler(recorder, req)

	return recorder
}

func setupTestDB(t *testing.T) {
	t.Helper()

	databaseURL := os.Getenv("TEST_DATABASE_URL")
	if databaseURL == "" {
		t.Fatal("TEST_DATABASE_URL is not set")
	}

	testDB, err := sql.Open("pgx", databaseURL)
	if err != nil {
		t.Fatalf("failed to open test database: %v", err)
	}

	if err := testDB.Ping(); err != nil {
		testDB.Close()
		t.Fatalf("failed to connect to test database: %v", err)
	}

	originalDB := db
	db = testDB

	t.Cleanup(func() {
		db = originalDB
		testDB.Close()
	})
}

