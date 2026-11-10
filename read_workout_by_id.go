package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"
)

// note that this one only applies when
// particular workout given id is needed

func readWorkoutbyID(w http.ResponseWriter, r *http.Request) {

	// Find ID from the URL request, convert from string to integer
	id, err := strconv.Atoi(r.URL.Query().Get("id"))

	// Validate ID since it is now converted to integer
	if err != nil || id <= 0 {
		writeJSONError(w, "Invalid workout ID", http.StatusBadRequest)
		return
	}

	userID, ok := getUserID(r)
	if !ok {
		writeJSONError(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	// Create a new variable called workout response which will store all the response items we retrieve
	var workout WorkoutResponse

	// first you want to extract that particular workout, the items that can be directly extracted and are not slice

	// Creating a date variable to handle date issues and make sure date is in the string format that we will need later
	var date time.Time

	// Now Query
	err = db.QueryRow(
		`SELECT workouts.name,
				workouts.date,
				programs.name,
				workouts.program_id
		FROM workouts
		LEFT JOIN programs 
			ON programs.id = workouts.program_id
			AND workouts.user_id = programs.user_id
		WHERE workouts.id = $1
		AND workouts.deleted_at IS NULL
		AND workouts.user_id = $2`,
		id,
		userID,
	).Scan(&workout.Name, &date, &workout.ProgramName, &workout.ProgramID)

	// We check of no workout found error and database error (internal server error)
	if err != nil {

		// StatusNotFound error in case the associate ID does not exist
		if errors.Is(err, sql.ErrNoRows) {
			writeJSONError(w, "Workout not found", http.StatusNotFound)
			return
		}

		// Database error, i.e., internal server error for any database retrieval errors
		log.Printf("Failed to fetch workout %d: %v", id, err)
		writeJSONError(w, "Failed to fetch workout", http.StatusInternalServerError)
		return
	}

	// Fix date format

	workout.Date = date.Format("2006-01-02")
	workout.ID = id

	// Now we need to check the workout exercises associated
	exercise_rows, err := db.Query(
		`SELECT 
		workout_exercises.id, 
		workout_exercises.exercise_id, 
		exercises.name
		FROM workout_exercises
		JOIN exercises
			ON workout_exercises.exercise_id = exercises.id
		WHERE workout_exercises.workout_id = $1`,
		id)

	if err != nil {
		log.Printf("Failed to fetch workout exercises: %v", err)
		writeJSONError(w, "Failed to fetch workout exercises", http.StatusInternalServerError)
		return
	}

	defer exercise_rows.Close()

	// okay you have everything you need for workout exercise response
	// except for sets

	// now we scan everything from
	// workout_exercise rows
	// to the relevant struct

	// We will just have one workout exercise response because it is one ID, so we don't need to repeatedly keep
	// creating empty slices, hence we just declare one at the beginning
	// this is a empty slice of time workout exercise responses created here
	workout.Exercises = []WorkoutExerciseResponse{}

	// Now iterating through the rows
	for exercise_rows.Next() {

		// i create one empty response that will contain the details for each exercise in the workout
		// for every exercise done and retrievied, I want to scan each one by one
		// and save the id, exercise id, and name
		var workoutExercise WorkoutExerciseResponse
		workoutExercise.Sets = []ExerciseSet{}

		err := exercise_rows.Scan(
			&workoutExercise.ID,
			&workoutExercise.ExerciseID,
			&workoutExercise.Name)

		if err != nil {
			log.Printf("Scan error: %v", err)
			writeJSONError(w, "Failed to process workout data", 500)
			return
		}

		// now we have everything we need for that particular workout exercise EXCEPT
		// EXCEPT for the sets slices.
		// We need a separate query and iteration for the SETs since sets is a slice itself in the
		// exercise reponse

		sets, err := db.Query(
			`SELECT id, workout_exercise_id, set_number, reps, weight 
			FROM exercise_sets
			WHERE workout_exercise_id = $1
			ORDER BY set_number ASC`,
			workoutExercise.ID)

		if err != nil {
			log.Printf("Failed to fetch workout sets: %v", err)
			writeJSONError(w, "Failed to fetch workout sets", http.StatusInternalServerError)
			return
		}

		for sets.Next() {
			var exerciseSet ExerciseSet
			err := sets.Scan(
				&exerciseSet.ID,
				&exerciseSet.WorkoutExerciseID,
				&exerciseSet.SetNumber,
				&exerciseSet.Reps,
				&exerciseSet.Weight,
			)

			if err != nil {
				sets.Close()
				log.Printf("Scan error: %v", err)
				writeJSONError(w, "Failed to process sets data", 500)
				return
			}

			workoutExercise.Sets = append(workoutExercise.Sets, exerciseSet)

		}

		if err := sets.Err(); err != nil {
			sets.Close()
			log.Printf("Error iterating through sets: %v", err)
			writeJSONError(w, "Failed to read sets", http.StatusInternalServerError)
			return
		}

		sets.Close()

		workout.Exercises = append(workout.Exercises, workoutExercise)
	}

	if err := exercise_rows.Err(); err != nil {
		log.Printf("Error iterating through workout exercises: %v", err)
		writeJSONError(w, "Failed to read workout exercises", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(workout); err != nil {
		log.Printf("Failed to encode workout %d: %v", id, err)
	}
}
