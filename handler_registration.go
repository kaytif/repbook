package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func registerUser(w http.ResponseWriter, r *http.Request) {

	// create an empty registration
	var newRegistration RegisterRequest

	// first take the incoming request and store it as a registration

	err := json.NewDecoder(r.Body).Decode(&newRegistration)
	if err != nil {
		writeJSONError(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// now validate the registration details
	if newRegistration.Username == "" || newRegistration.Password == "" {
		writeJSONError(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	// now we gotta hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(newRegistration.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("failed to hash the password:", err) // Real error for developer
		writeJSONError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	var newUser User
	newUser.Username = newRegistration.Username
	newUser.PasswordHash = string(hash)

	// Search whether username already exists
	err = db.QueryRow(
		"SELECT username FROM users WHERE username = $1",
		newUser.Username,
	).Scan(&newUser.Username)

	// Catch errors except for no row errors
	if err != nil && errors.Is(err, sql.ErrNoRows) != true {
		log.Println("failed to select username from database:", err) // Real error for developer
		writeJSONError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if err == nil && errors.Is(err, sql.ErrNoRows) != true {
		writeJSONError(w, "Username is not unique", http.StatusConflict)
		return
	}

	err = db.QueryRow(
		"INSERT INTO users (username, password_hash) VALUES ($1, $2) RETURNING id",
		newUser.Username,
		newUser.PasswordHash,
	).Scan(&newUser.ID)

	if err != nil {
		log.Println("failed to insert user:", err) // Real error for developer
		writeJSONError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	// User succesfull created
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(newUser)
}
