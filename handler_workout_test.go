package main

import (
	"net/http"
	"testing"
)

func TestCreateWorkout_Validation(t *testing.T) {
	tests := []struct {
		name       string
		handler    http.HandlerFunc
		method     string
		path       string
		body       string
		wantStatus int
	}{

		{"missing required fields", createWorkout, "POST", "/workouts", `{}`, 400},

		{"invalid date", createWorkout, "POST", "/workouts", `{"name":"Chest day", "date":"hello"}`, 400},

		{"empty workout name", createWorkout, "POST", "/workouts", `{"name":"","date":"2026-10-30"}`, 400},

		{"malformed JSON", createWorkout, "POST", "/workouts", `{"name":`, 400},

		{"wrong data type", createWorkout, "POST", "/workouts", `{"name":123, "date":"2026-10-30"}`, 400},

		{"whitespace name only", createWorkout, "POST", "/workouts", `{"name": "  ", "date":"2026-10-30"}`, 400},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := executeRequest(t, tt.handler, tt.method, tt.path, tt.body)

			if recorder.Code != tt.wantStatus {
				t.Errorf("expected %d, got %d; response: %s",
					tt.wantStatus,
					recorder.Code,
					recorder.Body.String())
			}
		})
	}
}

func TestUpdateWorkout_Validation(t *testing.T) {
	tests := []struct {
		name       string
		handler    http.HandlerFunc
		method     string
		path       string
		body       string
		wantStatus int
	}{

		// Checking ID
		{"missing id", updateWorkout, "PUT", "/workouts", `{}`, 400},

		{"wrong data type", updateWorkout, "PUT", "/workouts?id=abc", `{}`, 400},

		{"zero", updateWorkout, "PUT", "/workouts?id=0", `{}`, 400},

		{"negative ID", updateWorkout, "PUT", "/workouts?id=-1", `{}`, 400},

		// Checking JSON body
		{"malformed JSON", updateWorkout, "PUT", "/workouts?id=1", `{}`, 400},

		{"wrong JSON data type", updateWorkout, "PUT", "/workouts?id=1", `{"name":123}`, 400},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := executeRequest(t, tt.handler, tt.method, tt.path, tt.body)

			if recorder.Code != tt.wantStatus {
				t.Errorf("expected %d, got %d", tt.wantStatus, recorder.Code)
			}
		})
	}
}


