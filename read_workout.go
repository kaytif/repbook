package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func readWorkout(w http.ResponseWriter, r *http.Request) {

	userID, ok := getUserID(r)
	if !ok {
		writeJSONError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	rows, err := db.Query(`
		SELECT w.id, w.name, w.date, p.name
		FROM workouts w
		LEFT JOIN programs p 
				ON p.id = w.program_id
				AND p.user_id = w.user_id 
		WHERE w.deleted_at IS NULL
		AND w.user_id = $1
		ORDER BY w.date DESC, w.id DESC`,
		userID)

	if err != nil {
		log.Printf("Failed to fetch workouts: %v", err)
		writeJSONError(w, "Failed to fetch workouts", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	workouts := []WorkoutResponse{}

	for rows.Next() {
		var workout WorkoutResponse
		var date time.Time

		err := rows.Scan(
			&workout.ID,
			&workout.Name,
			&date,
			&workout.ProgramName,
		)

		if err != nil {
			log.Printf("Failed to scan workout row: %v", err)
			writeJSONError(w, "Failed to process workout data", http.StatusInternalServerError)
			return
		}

		workout.Date = date.Format("2006-01-02")
		workout.Exercises = []WorkoutExerciseResponse{}

		workouts = append(workouts, workout)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Failed to iterate through workout rows: %v", err)
		writeJSONError(w, "Failed to read workouts", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(workouts); err != nil {
		log.Printf("Failed to encode workouts response: %v", err)
	}
}
