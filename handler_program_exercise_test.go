package main

import (
	"net/http"
	"testing"
)

func TestCreateProgramExercise_Validation(t *testing.T) {
	tests := []struct {
		name       string
		handler    http.HandlerFunc
		method     string
		path       string
		body       string
		wantStatus int
	}{
		// Malformed JSON---------------------------------
		{"malformed JSON", createProgramExercise, "POST", "/program-exercises",
			`{"program_id":}`, 400},

		// Program ID validation----------------------------
		{"missing program_id", createProgramExercise, "POST", "/program-exercises",
			`{"exercise_id":1}`, 400},

		{"wrong program_id type", createProgramExercise, "POST", "/program-exercises",
			`{"program_id":"hello","exercise_id":1}`, 400},

		{"zero program_id", createProgramExercise, "POST", "/program-exercises",
			`{"program_id":0,"exercise_id":1}`, 400},

		{"negative program_id", createProgramExercise, "POST", "/program-exercises",
			`{"program_id":-1,"exercise_id":1}`, 400},

		// Exercise_id validation-------------------------
		{"missing exercise_id", createProgramExercise, "POST", "/program-exercises",
			`{"program_id":1}`, 400},

		{"wrong exercise_id type", createProgramExercise, "POST", "/program-exercises",
			`{"program_id":1,"exercise_id":"hello"}`, 400},

		{"zero exercise_id", createProgramExercise, "POST", "/program-exercises",
			`{"program_id":1,"exercise_id":0}`, 400},

		{"negative exercise_id", createProgramExercise, "POST", "/program-exercises",
			`{"program_id":1,"exercise_id":-1}`, 400},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := executeRequest(
				t,
				tt.handler,
				tt.method,
				tt.path,
				tt.body,
			)

			if recorder.Code != tt.wantStatus {
				t.Errorf(
					"expected %d, got %d; response: %s",
					tt.wantStatus,
					recorder.Code,
					recorder.Body.String(),
				)
			}
		})
	}
}
