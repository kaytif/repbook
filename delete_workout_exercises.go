package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5/pgconn"
)

func deleteWorkoutExercise(w http.ResponseWriter, r *http.Request) {

	// first read ID from URL
	idString := r.URL.Query().Get("id")

	// convert ID from string to integer
	id, err := strconv.Atoi(idString)

	// check if id is valid
	if err != nil || id <= 0 {
		writeJSONError(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	userID, ok := getUserID(r)
	if !ok {
		writeJSONError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// delete workout exercise
	result, err := db.Exec(`
		DELETE FROM workout_exercises we
		USING workouts w
		WHERE we.workout_id = w.id
		AND w.user_id = $1
		AND we.id = $2`,
		userID,
		id,
	)

	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == "23503" {
			writeJSONError(w, "Cannot delete workout exercise with existing sets", http.StatusConflict)
			return
		}

		writeJSONError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// check whether anything was deleted
	rowsAffected, err := result.RowsAffected()

	if err != nil {
		writeJSONError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		writeJSONError(w, "Workout exercise not found", http.StatusNotFound)
		return
	}

	// successful deletion
	w.WriteHeader(http.StatusNoContent)
}
