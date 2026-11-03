package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/jackc/pgx/v5/pgconn"
)

func createExercise(w http.ResponseWriter, r *http.Request) {

	// Create Exercise var to store incoming exercise in
	var exercise Exercise

	// first we need to decode the exercise
	err := json.NewDecoder(r.Body).Decode(&exercise)

	// check for any json errors
	if err != nil {
		writeJSONError(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	exercise.Name = strings.TrimSpace(exercise.Name)

	// ------ name validation --------------

	// 1. validate that name is not empty
	err = validation.Validate(
		exercise.Name,
		validation.Required,
	)

	if err != nil {
		writeJSONError(w, "Name not given", http.StatusBadRequest)
		return
	}

	err = db.QueryRow(
		`INSERT 
		INTO exercises
		(name)
		VALUES ($1)
		RETURNING id`,
		exercise.Name,
	).Scan(&exercise.ID)

	// id validation not needed since postgresql generates it for us.
	if err != nil {

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			writeJSONError(w, "Exercise already exists", http.StatusConflict)
			return
		}
		log.Println("Database insert error:", err)
		writeJSONError(w, "Internal Server Error", http.StatusInternalServerError)
		return

	}

	// set the header
	w.Header().Set("content-Type", "application/json")

	// Workout successfully created
	w.WriteHeader(http.StatusCreated)

	// Send created workout back as JSON
	json.NewEncoder(w).Encode(exercise)

}
