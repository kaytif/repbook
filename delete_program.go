package main

import (
	"net/http"
	"strconv"
)

func deleteProgram(w http.ResponseWriter, r *http.Request) {
	// you see what the incoming id is
	// validate the id
	// then run it through delete

	// first read ID from URL
	idString := r.URL.Query().Get("id")

	// convert ID from string to integer
	id, err := strconv.Atoi(idString)

	// check if id i valid
	if err != nil || id <= 0 {
		writeJSONError(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	userID, ok := getUserID(r)
	if !ok {
		writeJSONError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// now you straight up jump into deletion
	// if you get a database error because porgram does not exist then you send an error
	result, err := db.Exec(
		`DELETE FROM programs
		WHERE id = $1
		AND user_id = $2
		`,
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

	// If not matching program not existed
	if rowsAffected == 0 {
		writeJSONError(w, "Program not found", http.StatusNotFound)
		return

	}

	// Successful delete no body return
	w.WriteHeader(http.StatusNoContent)

}
