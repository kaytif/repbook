package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"net/http"
	"strconv"
	"time"
)

func updateWorkout(w http.ResponseWriter, r *http.Request) {

	/// ----- ID validation first ----

	idString := r.URL.Query().Get("id")

	// and does not have any issues
	// we do that by converting to string
	// then we convert from string to integer

	// use the ozzo package to conduct usual id validation checks:

	// 1. no id given
	if err := validation.Validate(
		idString,
		validation.Required.Error("ID is required"),
	); err != nil {
		writeJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 2. id is invalid i.e. like non integer
	id, err := strconv.Atoi(idString)
	if err != nil {
		writeJSONError(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// 3. id is 0 or negative
	if err := validation.Validate(
		id,
		validation.Min(1).Error("ID must be greater than 0"),
	); err != nil {
		writeJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// get the verified user ID
	userID, ok := getUserID(r)
	if !ok {
		writeJSONError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// ----- Decoding the input ------

	// if all the id checks pass, then create a struct and decode the id into the struct
	var workout Workout

	err = json.NewDecoder(r.Body).Decode(&workout)
	if err != nil {
		writeJSONError(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// ----- Checking if workout actually exists ------

	// now we check that the workout actually exists based on the workout id
	err = db.QueryRow(
		`SELECT id FROM workouts WHERE id = $1 and deleted_at IS NULL AND user_id = $2`,
		id,
		userID,
	).Scan(&workout.ID)

	if errors.Is(err, sql.ErrNoRows) {
		writeJSONError(w, "Workout does not exist", http.StatusNotFound)
		return
	}

	if err != nil {
		writeJSONError(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// ------ Checking the updated workout like we would a new workout entry ---
	// Keep a date label for storage compatibility; users do not name workouts.
	if workout.Name == "" {
		workout.Name = workout.Date
	}

	// ----------- date validation ----------------

	// make sure it is not empty
	if workout.Date == "" {
		writeJSONError(w, "Workout date is required", http.StatusBadRequest)
		return
	}
	if _, err := time.Parse("2006-01-02", workout.Date); err != nil {
		writeJSONError(w, "Invalid date format", http.StatusBadRequest)
		return
	}

	// Parse the date
	// workoutDate, err := time.Parse("2006-01-02", workout.Date)
	// if err != nil {
	//	writeJSONError(w, "Invalid date format", http.StatusBadRequest)
	//	return
	// }

	// Make sure the date is not greater than today
	// now := time.Now()
	// today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	// if workoutDate.After(today) {
	//	writeJSONError(w, "Workout date can not be in the future", http.StatusBadRequest)
	//	return
	// }

	// Now we check that if the program_id was updated that the program updated to actually exists
	if workout.ProgramID != nil {
		var programID int

		err = db.QueryRow(
			`SELECT programs.id
			FROM programs
			where id = $1
			AND programs.user_id = $2`,
			*workout.ProgramID,
			userID,
		).Scan(&programID)

		// give error if no program with this id even exists
		if errors.Is(err, sql.ErrNoRows) {
			writeJSONError(w, "Program not found", http.StatusNotFound)
			return
		}

		// give error if no program with this id even exists
		if err != nil {
			writeJSONError(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

	}

	// once all is said and done, let's update the workout
	result, err := db.Exec(
		`UPDATE workouts 
		SET name = $1,
			date = $2,
			program_id = $3,
			updated_at = NOW()
		WHERE id = $4 AND deleted_at IS NULL AND workouts.user_ID = $5`,
		workout.Name,
		workout.Date,
		workout.ProgramID,
		id,
		userID,
	)

	if err != nil {
		writeJSONError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Check how many rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		writeJSONError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// no matching workout existed
	if rowsAffected == 0 {
		writeJSONError(w, "Workout not found", http.StatusNotFound)
		return
	}

	// keep the ID from the url
	workout.ID = id

	w.Header().Set("Content-Type", "application/json")

	// send updated workout back as json
	json.NewEncoder(w).Encode(workout)
}
