package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func readProgram(w http.ResponseWriter, r *http.Request) {

	userID, ok := getUserID(r)
	if !ok {
		writeJSONError(w, "Internal server error", http.StatusInternalServerError)
		return
	}


	rows, err := db.Query(
		`SELECT programs.id, 
				programs.name
		FROM programs
		WHERE user_id = $1`,
		userID,
	)

	if err != nil {
		log.Printf("Failed to fetch programs: %v", err)
		writeJSONError(w, "Failed to fetch programs", http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	programs := []ProgramResponse{}

	for rows.Next() {
		var programResponse ProgramResponse

		err := rows.Scan(
			&programResponse.ID,
			&programResponse.Name,
		)

		if err != nil {
			log.Printf("Scan error: %v", err)
			writeJSONError(w, "Failed to process program data", 500)
			return
		}

		programResponse.Exercises = []ProgramExerciseResponse{}

		programs = append(programs, programResponse)

	}
	if err := rows.Err(); err != nil {
		log.Printf("Error iterating through programs: %v", err)
		writeJSONError(w, "Failed to read programs", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(programs); err != nil {
		log.Printf("Failed to encode programs response: %v", err)
	}

}
