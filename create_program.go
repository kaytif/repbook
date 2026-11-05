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

func createProgram(w http.ResponseWriter, r *http.Request) {

	userID, ok := getUserID(r)
	if !ok {
		writeJSONError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Create Program var to store incoming program in
	var program Program

	// first we need to decode the program
	err := json.NewDecoder(r.Body).Decode(&program)

	// check for any json errors
	if err != nil {
		writeJSONError(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	program.Name = strings.TrimSpace(program.Name)

	// ------ name validation --------------

	// 1. validate that name is not empty
	err = validation.Validate(
		program.Name,
		validation.Required,
	)

	if err != nil {
		writeJSONError(w, "Name not given", http.StatusBadRequest)
		return
	}

	err = db.QueryRow(
		`INSERT 
		INTO programs
		(name, user_id)
		VALUES ($1, $2)
		RETURNING id`,
		program.Name,
		userID,
	).Scan(&program.ID)

	// id validation not needed since postgresql generates it for us.

	if err != nil {

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			writeJSONError(w, "Program already exists", http.StatusConflict)
			return
		}
		log.Println("Database insert error:", err)
		writeJSONError(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// set the header
	w.Header().Set("content-Type", "application/json")

	// Program successfully created
	w.WriteHeader(http.StatusCreated)

	// Send created Program back as JSON
	json.NewEncoder(w).Encode(program)
}
