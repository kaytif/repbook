package main

import (
	"encoding/json"
	"errors"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/jackc/pgx/v5/pgconn"
	"log"
	"net/http"
	"strings"
)

func createWorkout(w http.ResponseWriter, r *http.Request) {

	userID, ok := getUserID(r)
	if !ok {
		writeJSONError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// create an empty workout
	var newWorkout Workout

	// first decode the request
	err := json.NewDecoder(r.Body).Decode(&newWorkout)

	if err != nil {
		writeJSONError(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// ----------- name validation ----------------

	// validate that name is not empty
	err = validation.Validate(
		strings.TrimSpace(newWorkout.Name),
		validation.Required,
	)
	if err != nil {
		writeJSONError(w, "Name not given", http.StatusBadRequest)
		return
	}

	// ----------- date validation ----------------

	// make sure it is not empty
	if newWorkout.Date == "" {
		writeJSONError(w, "Workout date is required", http.StatusBadRequest)
		return
	}

	// Parse the date
	// workoutDate, err := time.Parse("2006-01-02", newWorkout.Date)
	// if err != nil {
	// 	writeJSONError(w, "Invalid date format", http.StatusBadRequest)
	//	return
	// }

	// Make sure the date is not greater than today
	// now := time.Now()
	// today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	// if workoutDate.After(today) {
	//	writeJSONError(w, "Workout date can not be in the future", http.StatusBadRequest)
	//	return
	//}

	// now inserting into postgre sql. we need to scan the rows and get the id
	// this is because the frontend might need to use it
	err = db.QueryRow(
		`INSERT 
		INTO workouts 
		(name, 
		workout_date, 
		program_id,
		user_id) 
		VALUES ($1, $2, $3, $4) 
		RETURNING id`,
		newWorkout.Name,
		newWorkout.Date,
		newWorkout.ProgramID,
		userID,
	).Scan(&newWorkout.ID)

	// catch error when program id does not exist
	var pgErr *pgconn.PgError

	if errors.As(err, &pgErr) && pgErr.Code == "23503" {
		writeJSONError(w, "Program does not exist", http.StatusBadRequest)
		return
	}

	if err != nil {
		log.Println("Database insert error:", err)
		writeJSONError(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Set the header
	w.Header().Set("Content-Type", "application/json")

	// Workout succesfuly created
	w.WriteHeader(http.StatusCreated)

	// Send created workout back as JSON
	json.NewEncoder(w).Encode(newWorkout)

}
