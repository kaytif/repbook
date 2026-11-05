package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/jackc/pgx/v5/pgconn"
)

func createWorkoutExercise(w http.ResponseWriter, r *http.Request) {
	var workoutExercise WorkoutExercise

	// get the verified user ID
	userID, ok := getUserID(r)
	if !ok {
		writeJSONError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// decode
	err := json.NewDecoder(r.Body).Decode(&workoutExercise)
	if err != nil {
		writeJSONError(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// we need to check the workout id that is given is accurate, etc
	err = validation.Validate(workoutExercise.WorkoutID, validation.Required, validation.Min(1))
	if err != nil {
		writeJSONError(w, "Invalid workout ID", http.StatusBadRequest)
		return
	}

	// repeat for exercise ID
	err = validation.Validate(workoutExercise.ExerciseID, validation.Required, validation.Min(1))
	if err != nil {
		writeJSONError(w, "Invalid exercise ID", http.StatusBadRequest)
		return
	}

	// then we need to insert, and while inserting we need
	// to check that it actually exists

	// now inserting into postgre sql. we need to scan the rows and get the id
	// this is because the frontend might need to use it
	err = db.QueryRow(
		`INSERT 
		INTO workout_exercises
		(workout_id, exercise_id)
		SELECT w.id, e.id
		FROM workouts w
		CROSS JOIN exercises e
		WHERE w.id = $1
		AND e.id = $2
		AND w.user_id = $3
		AND e.user_id = $3
		AND w.deleted_at is NULL
		RETURNING id`,
		workoutExercise.WorkoutID,
		workoutExercise.ExerciseID,
		userID).Scan(&workoutExercise.ID)

	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) {
			// foreign key violation
			if pgErr.Code == "23503" {
				switch pgErr.ConstraintName {
				case "workout_exercises_workout_id_fkey":
					writeJSONError(w, "Workout does not exist", http.StatusBadRequest)
					return

				case "workout_exercises_exercise_id_fkey":
					writeJSONError(w, "Exercise does not exist", http.StatusBadRequest)
					return

				}
			}

		}

		if errors.Is(err, sql.ErrNoRows) {
			writeJSONError(w, "Workout or exercise not found", http.StatusNotFound)
			return
		}

		log.Println("Database insert error:", err)
		writeJSONError(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Set the header
	w.Header().Set("Content-Type", "application/json")

	// Workoutexercise succesfuly created
	w.WriteHeader(http.StatusCreated)

	// Send created workoutexercise back as JSON
	json.NewEncoder(w).Encode(workoutExercise)

}
