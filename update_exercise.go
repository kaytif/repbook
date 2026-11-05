package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/jackc/pgx/v5/pgconn"
)

func updateExercise(w http.ResponseWriter, r *http.Request) {

	// Read ID from URL
	idString := r.URL.Query().Get("id")

	// this is a check on id as well
	// make sure you use ozzo here as well, later, for consistency

	// Convert ID from string to integer
	id, err := strconv.Atoi(idString)
	if err != nil || id <= 0 {
		writeJSONError(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// get the verified user ID
	userID, ok := getUserID(r)
	if !ok {
		writeJSONError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Create an empty exercise
	var exercise Exercise

	// Decode incoming JSON
	err = json.NewDecoder(r.Body).Decode(&exercise)
	if err != nil {
		writeJSONError(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	exercise.Name = strings.TrimSpace(exercise.Name)

	// Validate that exercise name is not empty
	err = validation.Validate(
		exercise.Name,
		validation.Required,
	)

	if err != nil {
		writeJSONError(w, "exercise name is required", http.StatusBadRequest)
		return
	}

	// ------ Check whether exercise to update exists -----

	// Validate that exercise name is unique
	// The error will come through from postgresql and we need to
	// check the error code and give the required alert

	// --- Validate that exercise name is unique

	// PostgreSQL does it for you, here we show the error if PostgreSQL gives us an error
	// Note that you can only update the new exercise name and not id because id is already assigned itself by SQL
	result, err := db.Exec(
		`UPDATE exercises SET name = $1 WHERE id = $2 AND user_id = $3`,
		exercise.Name,
		id,
		userID,
	)

	if err != nil {
		var pgErr *pgconn.PgError

		// Duplicate exercise name
		// PostgreSQL code 23505 = UNIQUE constraint violation
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			writeJSONError(w, "Exercise name already exists", http.StatusConflict)
			return
		}

		// We check here for an error in case there was an error in updating exercises
		// Here the error checks for errors that happen while PostgreSQL is executing the UPDATE for e.g.;
		// database connection died, check constraint violated, foreign key violated, SQL statement itself is invalid

		writeJSONError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Check how many rows were updated
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		writeJSONError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// No matching exercise existed
	if rowsAffected == 0 {
		writeJSONError(w, "Exercise not found", http.StatusNotFound)
		return
	}

	// Keep the ID from the URL
	exercise.ID = id

	// Set the header
	w.Header().Set("Content-Type", "application/json")

	// Send the updated exercise back as JSON
	json.NewEncoder(w).Encode(exercise)
}
