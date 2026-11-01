package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCreateWorkout_Validation(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		wantStatus int
	}{
		{
			name:       "empty name",
			body:       `{"name":"","date":"2026-10-30"}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "invalid date",
			body:       `{"name":"Chest Day","date":"hello"}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "missing name",
			body:       `{"date":"2026-10-30"}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "malformed JSON",
			body:       `{"name":"Chest Day","date":}`,
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := executeRequest(
				t, 
				createWorkout,
				http.MethodPost,
				"/workouts",
				tt.body,
			)

			if recorder.Code != tt.wantStatus {
				t.Errorf(
					"expected status %d, got %d",
					tt.wantStatus,
					recorder.Code,
				)
			}
		})
	}
}

func TestCreateWorkout_DatabaseFailure(t *testing.T) {
    originalDB := db

    failedDB, err := sql.Open(
        "pgx",
        "postgres://workout:invalid@127.0.0.1:1/nonexistent?sslmode=disable&connect_timeout=1",
    )
    if err != nil {
        t.Fatalf("failed to create database handle: %v", err)
    }

    db = failedDB

    t.Cleanup(func() {
        db = originalDB
        failedDB.Close()
    })

    body := `{"name":"Chest Day","date":"2026-10-30"}`

    recorder := executeRequest(
        t,
        createWorkout,
        http.MethodPost,
        "/workouts",
        body,
    )

    if recorder.Code != http.StatusInternalServerError {
        t.Errorf(
            "expected status 500, got %d; response: %s",
            recorder.Code,
            recorder.Body.String(),
        )
    }
}