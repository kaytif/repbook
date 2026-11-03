package main

import (
	"net/http"
	"testing"
)

func TestWorkoutExercise_Validation(t *testing.T) {
	tests := []struct {
		name       string
		handler    http.HandlerFunc
		method     string
		path       string
		body       string
		wantStatus int
	}{
		// Malformed JSON---------------------------------
		{"malformed JSON", createWorkoutExercise, "POST", "/workout-exercises",
			`{"workout_id":}`, 400},

		// Workout exercise ID validation----------------------------
		{"missing workout_id", createWorkoutExercise, "POST", "/workout-exercises",
			`{"exercise_id":1}`, 400},

		{"wrong workout_id type", createWorkoutExercise, "POST", "/workout-exercises",
			`{"workout_id":"hello","exercise_id":1}`, 400},

		{"zero workout_id", createWorkoutExercise, "POST", "/workout-exercises",
			`{"workout_id":0,"exercise_id":1}`, 400},

		{"negative workout_id", createWorkoutExercise, "POST", "/workout-exercises",
			`{"workout_id":-1,"exercise_id":1}`, 400},

		// Set number validation-------------------------
		{"missing exercise_id", createWorkoutExercise, "POST", "/workout-exercises",
			`{"workout_id":1}`, 400},

		{"wrong exercise_id type", createWorkoutExercise, "POST", "/workout-exercises",
			`{"workout_id":1,"exercise_id":"hello"}`, 400},

		{"zero exercise_id", createWorkoutExercise, "POST", "/workout-exercises",
			`{"workout_id":1,"exercise_id":0}`, 400},

		{"negative exercise_id", createWorkoutExercise, "POST", "/workout-exercises",
			`{"workout_id":1,"exercise_id":-1}`, 400},
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
