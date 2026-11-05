package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// executeRequest creates and executes an authenticated HTTP request
// against a handler directly.
//
// It simulates the authentication context normally created by
// authMiddleware, without performing a real login.
//
// Use this helper for handler validation and database tests.
// Do NOT use it to test the authentication middleware itself.
func executeRequest(
	t *testing.T,
	handler http.HandlerFunc,
	method, path, body string,
) *httptest.ResponseRecorder {

	// Marks this function as a test helper.
	// If a test fails, Go reports the calling test's location
	// rather than pointing unnecessarily into this helper.
	t.Helper()

	// STEP 1: Create a simulated HTTP request.
	//
	// strings.NewReader converts the JSON string into an io.Reader,
	// which httptest.NewRequest accepts as the request body.
	req := httptest.NewRequest(
		method,
		path,
		strings.NewReader(body),
	)

	// STEP 2: Set the request's content type.
	//
	// This tells the handler that the request body contains JSON.
	req.Header.Set("Content-Type", "application/json")

	// STEP 3: Simulate the authentication middleware.
	//
	// In production, authMiddleware validates the session and
	// stores the authenticated user's ID in the request context.
	//
	// Since we are calling the handler directly, we bypass
	// authMiddleware and manually provide the user ID instead.
	//
	// IMPORTANT:
	// userIDKey must be the exact key used by your middleware.
	// User ID 1 must exist in the test database for DB-backed tests.
	ctx := context.WithValue(
		req.Context(),
		userIDKey,
		1,
	)

	// STEP 4: Attach the new context to the HTTP request.
	//
	// context.WithValue does not modify the original request.
	// WithContext returns a copy of the request using our new context.
	req = req.WithContext(ctx)

	// STEP 5: Create a response recorder.
	//
	// It captures the HTTP status code, response headers,
	// and response body produced by the handler.
	recorder := httptest.NewRecorder()

	// STEP 6: Execute the handler.
	//
	// This calls the handler directly, without going through
	// the HTTP router or authentication middleware.
	handler(recorder, req)

	// STEP 7: Return the response so the test can inspect it.
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
