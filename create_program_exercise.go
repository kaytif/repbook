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

func createProgramExercise(w http.ResponseWriter, r *http.Request) {

	userID, ok := getUserID(r)
	if !ok {
		writeJSONError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	var programExercise ProgramExercise

	// decode
	err := json.NewDecoder(r.Body).Decode(&programExercise)
	if err != nil {
		writeJSONError(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// we need to check the program id that is given is accurate, etc
	err = validation.Validate(programExercise.ProgramID, validation.Required, validation.Min(1))
	if err != nil {
		writeJSONError(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// repeat for exercise ID
	err = validation.Validate(programExercise.ExerciseID, validation.Required, validation.Min(1))
	if err != nil {
		writeJSONError(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// then we need to insert, and while inserting we need
	// to check that it actually exists

	// now inserting into postgre sql. we need to scan the rows and get the id
	// this is because the frontend might need to use it
	err = db.QueryRow(
		`INSERT 
		INTO program_exercises
		(program_id, exercise_id)
		SELECT p.id, e.id
		FROM programs p
		CROSS JOIN exercises e
		WHERE p.id = $1
		AND e.id = $2
		AND p.user_id = $3
		AND e.user_id = $3
		RETURNING id`,
		programExercise.ProgramID,
		programExercise.ExerciseID,
		userID).Scan(&programExercise.ID)

	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) {
			// foreign key violation
			if pgErr.Code == "23503" {
				switch pgErr.ConstraintName {
				case "program_exercises_program_id_fkey":
					writeJSONError(w, "Program does not exist", http.StatusBadRequest)
					return

				case "program_exercises_exercise_id_fkey":
					writeJSONError(w, "Exercise does not exist", http.StatusBadRequest)
					return

				}
			}
			// unique constraint violation
			if pgErr.Code == "23505" {
				writeJSONError(w, "Exercise already exists in this program", http.StatusConflict)
				return
			}
		}

		if errors.Is(err, sql.ErrNoRows) {
			writeJSONError(
				w,
				"Program or exercise not found",
					http.StatusNotFound,
			)
			return
		}

		log.Println("Database insert error:", err)
		writeJSONError(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Set the header
	w.Header().Set("Content-Type", "application/json")

	// Workout succesfuly created
	w.WriteHeader(http.StatusCreated)

	// Send created workout back as JSON
	json.NewEncoder(w).Encode(programExercise)

}
