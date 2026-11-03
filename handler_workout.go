package main

import "net/http"

// 1.Workout
func handlerWorkout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		if r.URL.Query().Has("id") {
			readWorkoutbyID(w, r)
		} else {
			readWorkout(w, r)
		}
	case http.MethodPost:
		createWorkout(w, r)

	case http.MethodPut:
		updateWorkout(w, r)

	case http.MethodDelete:
		deleteWorkout(w, r)

	default:
		writeJSONError(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// 2. Exercise
func handlerExercise(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		readExercise(w, r)
			
	case http.MethodPost:
		createExercise(w, r)

	case http.MethodPut:
		updateExercise(w, r)

	case http.MethodDelete:
		deleteExercise(w, r)

	default:
		writeJSONError(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// 3. Program
func handlerProgram(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		if r.URL.Query().Has("id") {
			readProgram(w, r)
		}
	case http.MethodPost:
		createProgram(w, r)

	case http.MethodPut:
		updateProgram(w, r)

	case http.MethodDelete:
		deleteProgram(w, r)

	default:
		writeJSONError(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// 4. Workout Exercise
func handlerWorkoutExercise(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {

	case http.MethodPost:
		createWorkoutExercise(w, r)

	case http.MethodDelete:
		deleteWorkoutExercise(w, r)

	default:
		writeJSONError(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// 5. Program Exercise
func handlerProgramExercise(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {

	case http.MethodPost:
		createProgramExercise(w, r)

	case http.MethodDelete:
		deleteProgramExercise(w, r)

	default:
		writeJSONError(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// 6. Exercise Set
func handlerExerciseSet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {

	case http.MethodPost:
		createExerciseSet(w, r)

	case http.MethodDelete:
		deleteExerciseSet(w, r)

	default:
		writeJSONError(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}