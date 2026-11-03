package main

import (
	"net/http"
	"net/http/httptest"
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

// func setupTestDB(t *testing.T) {
// 	t.Helper()

// 	databaseURL := os.Getenv("TEST_DATABASE_URL")
// 	if databaseURL == "" {
// 		t.Fatal("TEST_DATABASE_URL is not set")
// 	}

// 	testDB, err := sql.Open("pgx", databaseURL)
// 	if err != nil {
// 		t.Fatalf("failed to open test database: %v", err)
// 	}

// 	if err := testDB.Ping(); err != nil {
// 		testDB.Close()
// 		t.Fatalf("failed to connect to test database: %v", err)
// 	}

// 	originalDB := db
// 	db = testDB

// 	t.Cleanup(func() {
// 		db = originalDB
// 		testDB.Close()
// 	})
// }

// func TestCreateWorkout_DatabaseFailure(t *testing.T) {
// 	originalDB := db

// 	failedDB, err := sql.Open(
// 		"pgx",
// 		"postgres://workout:invalid@127.0.0.1:1/nonexistent?sslmode=disable&connect_timeout=1",
// 	)
// 	if err != nil {
// 		t.Fatalf("failed to create database handle: %v", err)
// 	}

// 	db = failedDB

// 	t.Cleanup(func() {
// 		db = originalDB
// 		failedDB.Close()
// 	})

// 	body := `{"name":"Chest Day","date":"2026-10-30"}`

// 	recorder := executeRequest(
// 		t,
// 		createWorkout,
// 		http.MethodPost,
// 		"/workouts",
// 		body,
// 	)

// 	if recorder.Code != http.StatusInternalServerError {
// 		t.Errorf(
// 			"expected status 500, got %d; response: %s",
// 			recorder.Code,
// 			recorder.Body.String(),
// 		)
// 	}
// }
