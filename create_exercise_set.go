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

func createExerciseSet(w http.ResponseWriter, r *http.Request) {

	// create an empty workout
	var exerciseSet ExerciseSet

	// first decode the request
	err := json.NewDecoder(r.Body).Decode(&exerciseSet)

	if err != nil {
		writeJSONError(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Validate WorkoutExerciseID
	err = validation.Validate(exerciseSet.WorkoutExerciseID, validation.Required, validation.Min(1))
	if err != nil {
		writeJSONError(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// Validate reps
	err = validation.Validate(exerciseSet.Reps, validation.Min(0))
	if err != nil {
		writeJSONError(w, "Invalid Reps", http.StatusBadRequest)
		return
	}

	// Validate weight
	err = validation.Validate(exerciseSet.Weight, validation.Min(0))
	if err != nil {
		writeJSONError(w, "Invalid Weight", http.StatusBadRequest)
		return
	}

	// ---- insert, and then you need to validate the things inserted
	tx, err := db.BeginTx(r.Context(), nil)
	if err != nil {
		writeJSONError(w, "Internal server error", 500)
		return
	}

	defer tx.Rollback()

	// Lock the workout exercise row

	var id int
	err = tx.QueryRow(`
	SELECT id
	FROM workout_exercises
	WHERE id = $1
	FOR UPDATE`,
		exerciseSet.WorkoutExerciseID).Scan(&id)

	if errors.Is(err, sql.ErrNoRows) {
		writeJSONError(w, "Workout exercise not found", 404)
		return
	}

	if err != nil {
		writeJSONError(w, "Internal server error", 500)
		return
	}

	// Calculate next set number and insert
	err = tx.QueryRow(`
	INSERT INTO exercise_sets (
		workout_exercise_id,
		set_number,
		reps,
		weight_kg
	)
	SELECT 
	$1,
	COALESCE(MAX(set_number), 0) +1,
		$2,
		$3
	FROM exercise_sets
	WHERE workout_exercise_id = $1
	RETURNING id, set_number`,
		exerciseSet.WorkoutExerciseID,
		exerciseSet.Reps,
		exerciseSet.Weight,
	).Scan(&exerciseSet.ID, &exerciseSet.SetNumber)

	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) {
			// foreign key violation
			if pgErr.Code == "23503" {
				switch pgErr.ConstraintName {
				case "exercise_sets_workout_exercise_id_fkey":
					writeJSONError(w, "Workout exercise does not exist", http.StatusBadRequest)
					return
				}
			}
			// unique constraint violation
			if pgErr.Code == "23505" {
				writeJSONError(w, "Set number already exsits for this workout", http.StatusConflict)
				return
			}
		}

		log.Println("Database insert error:", err)
		writeJSONError(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = tx.Commit()
	if err != nil {
		log.Println("Transaction commit error:", err)
		writeJSONError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Setting header
	w.Header().Set("Content-Type", "application/json")

	// Returning HTTP 201
	w.WriteHeader(http.StatusCreated)

	// send updated workout_exercise back as a josn
	json.NewEncoder(w).Encode(exerciseSet)

}
