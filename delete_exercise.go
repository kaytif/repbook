package main

import (
	"net/http"
	"strconv"
)

func deleteExercise(w http.ResponseWriter, r *http.Request) {

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

	// delete exercise
	result, err := db.Exec(`
        DELETE FROM exercises
        WHERE id = $1
		AND user_id = $2 
    `, id, userID)

	if err != nil {
		writeJSONError(w, "Cannot delete exercise because it may be referenced by existing records", http.StatusConflict)
		return
	}

	// check whether anything was deleted
	rowsAffected, err := result.RowsAffected()

	if err != nil {
		writeJSONError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		writeJSONError(w, "Exercise not found", http.StatusNotFound)
		return
	}

	// successful deletion
	w.WriteHeader(http.StatusNoContent)
}
