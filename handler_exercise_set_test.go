package main

import (
	"net/http"
	"testing"
)

func TestCreateExerciseSet_Validation(t *testing.T) {
	tests := []struct {
		name       string
		handler    http.HandlerFunc
		method     string
		path       string
		body       string
		wantStatus int
	}{
		// Malformed JSON---------------------------------
		{"malformed JSON", createExerciseSet, "POST", "/exercise-sets",
			`{"workout_exercise_id":}`, 400},

		// Workout exercise ID validation----------------------------
		{"missing workout_exercise_id", createExerciseSet, "POST", "/exercise-sets",
			`{"set_number":1}`, 400},

		{"wrong workout_exercise_id type", createExerciseSet, "POST", "/exercise-sets",
			`{"workout_exercise_id":"hello","set_number":1}`, 400},

		{"zero workout_exercise_id", createExerciseSet, "POST", "/exercise-sets",
			`{"workout_exercise_id":0,"set_number":1}`, 400},

		{"negative workout_exercise_id", createExerciseSet, "POST", "/exercise-sets",
			`{"workout_exercise_id":-1,"set_number":1}`, 400},

		// Set number validation-------------------------
		// {"missing set_number", createExerciseSet, "POST", "/exercise-sets",
		// 	`{"workout_exercise_id":1}`, 400},

		// {"wrong set_number type", createExerciseSet, "POST", "/exercise-sets",
		// 	`{"workout_exercise_id":1,"set_number":"hello"}`, 400},

		// {"zero set_number", createExerciseSet, "POST", "/exercise-sets",
		// 	`{"workout_exercise_id":1,"set_number":0}`, 400},

		// {"negative set_number", createExerciseSet, "POST", "/exercise-sets",
		// 	`{"workout_exercise_id":1,"set_number":-1}`, 400},

		// Reps validation ----------------
		{"negative reps", createExerciseSet, "POST", "/exercise-sets",
			`{"workout_exercise_id":1,"set_number":-1, "reps":-1}`, 400},

		{"wrong reps type", createExerciseSet, "POST", "/exercise-sets",
			`{"workout_exercise_id":1,"set_number":-1, "reps":"hello"}`, 400},

		// Weight validation ---------------
		{"negative weight", createExerciseSet, "POST", "/exercise-sets",
			`{"workout_exercise_id":1,"set_number":-1, "weight":-1}`, 400},

		{"wrong weights type", createExerciseSet, "POST", "/exercise-sets",
			`{"workout_exercise_id":1,"set_number":-1, "weight":"hello"}`, 400},
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
