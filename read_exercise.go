package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func readExercise(w http.ResponseWriter, r *http.Request) {

	userID, ok := getUserID(r)
	if !ok {
		writeJSONError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	rows, err := db.Query(
		`SELECT exercises.id, 
				exercises.name
		FROM exercises
		WHERE exercises.user_id = $1`,
		userID,
	)

	if err != nil {
		log.Printf("Failed to fetch exercises: %v", err)
		writeJSONError(w, "Failed to fetch exercises", http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	exercises := []Exercise{}

	for rows.Next() {
		var exerciseResponse Exercise

		err := rows.Scan(
			&exerciseResponse.ID,
			&exerciseResponse.Name,
		)

		if err != nil {
			log.Printf("Scan error: %v", err)
			writeJSONError(w, "Failed to process exercise data", http.StatusInternalServerError)
			return
		}

		exercises = append(exercises, exerciseResponse)

	}
	if err := rows.Err(); err != nil {
		log.Printf("Error iterating through exercises: %v", err)
		writeJSONError(w, "Failed to read exercises", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(exercises); err != nil {
		log.Printf("Failed to encode exercises response: %v", err)
	}

}
