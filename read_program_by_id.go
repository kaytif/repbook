package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
)

func readProgramByID(w http.ResponseWriter, r *http.Request) {
	// Find ID from the URL request, convert from string to integer
	id, err := strconv.Atoi(r.PathValue("id"))

	// Validate ID since it is now converted to integer
	if err != nil || id <= 0 {
		writeJSONError(w, "Invalid program ID", http.StatusBadRequest)
		return
	}

	userID, ok := getUserID(r)
	if !ok {
		writeJSONError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// initialize empty program exercise response in the program response
	var programResponse ProgramResponse

	err = db.QueryRow(
		`SELECT programs.id, 
				programs.name
		FROM programs
		WHERE programs.id = $1
		and programs.user_id = $2`,
		id,
		userID,
	).Scan(&programResponse.ID, &programResponse.Name)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeJSONError(w, "Program not found", http.StatusNotFound)
			return
		}

		log.Printf("Failed to fetch program %d: %v", id, err)
		writeJSONError(w, "Failed to fetch program", http.StatusInternalServerError)
		return
	}

	rows, err := db.Query(
		`SELECT exercises.id, exercises.name
		FROM exercises
		JOIN program_exercises on
			program_exercises.exercise_id = exercises.id
		WHERE program_exercises.program_id = $1`,
		id,
	)

	if err != nil {
		log.Printf("Failed to fetch program exercises: %v", err)
		writeJSONError(w, "Failed to fetch program exercises", http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	programResponse.Exercises = []ProgramExerciseResponse{}

	// we basically need to fill this up

	for rows.Next() {
		var programExercise ProgramExerciseResponse

		err := rows.Scan(
			&programExercise.ExerciseID,
			&programExercise.Name,
		)

		if err != nil {
			log.Printf("Scan error: %v", err)
			writeJSONError(w, "Failed to process program data", 500)
			return
		}

		programResponse.Exercises = append(programResponse.Exercises, programExercise)

	}

	if err := rows.Err(); err != nil {
		log.Printf("Error iterating through program response: %v", err)
		writeJSONError(w, "Failed to read program response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(programResponse); err != nil {
		log.Printf("Failed to encode programs response: %v", err)
	}

}

// we start with extracting the ID and name
// note that there is just one id and name so we use
// query.row and do it just once
// after that we do something like extract
// exercise information
// the we take the db created from the exercise
// information
// we scan the stuff into individual empty structs
// then keep appending to the main struct
