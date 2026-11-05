package main

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"net/http"
	"strconv"
)

func deleteWorkout(w http.ResponseWriter, r *http.Request) {
	idString := r.URL.Query().Get("id")

	// convert from string to integer
	id, err := strconv.Atoi(idString)

	// validation checks on ID (same as the ones in update workout)
	// 1. no id given
	if err := validation.Validate(
		idString,
		validation.Required.Error("ID is required"),
	); err != nil {
		writeJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 2. id is invalid i.e. like non integer
	id, err = strconv.Atoi(idString)
	if err != nil {
		writeJSONError(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// 3. id is not 0 or negative
	if err := validation.Validate(
		id,
		validation.Min(1).Error("ID must be greater than 0"),
	); err != nil {
		writeJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID, ok := getUserID(r)
	if !ok {
		writeJSONError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// DELETE transaction from PostgreSQL
	result, err := db.Exec(
		`UPDATE workouts
		'SET deleted_at = NOW()
		WHERE id = $1 
		AND deleted_at is NULL
		AND user_id = $2`,
		id,
		userID,
	)

	if err != nil {
		writeJSONError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Check whether if anything was deleted
	rowsAffected, err := result.RowsAffected()

	if err != nil {
		writeJSONError(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// If not matching account not existed
	if rowsAffected == 0 {
		writeJSONError(w, "Transaction not found", http.StatusNotFound)
		return

	}

	// Successful delete no body return
	w.WriteHeader(http.StatusNoContent)

}
