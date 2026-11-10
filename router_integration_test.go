package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
	"time"
)

// Run against PostgreSQL with TEST_DATABASE_URL. Each run uses a private schema;
// existing tables and records in that database are never used or modified.
func TestRouterIntegration(t *testing.T) {
	databaseURL := os.Getenv("TEST_DATABASE_URL")
	if databaseURL == "" {
		t.Skip("set TEST_DATABASE_URL to run PostgreSQL integration checks")
	}
	admin, err := sql.Open("pgx", databaseURL)
	if err != nil {
		t.Fatal(err)
	}
	defer admin.Close()
	schema := fmt.Sprintf("workout_test_%d", time.Now().UnixNano())
	if _, err := admin.Exec("CREATE SCHEMA " + schema); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if _, err := admin.Exec("DROP SCHEMA " + schema + " CASCADE"); err != nil {
			t.Error(err)
		}
	}()
	parsed, err := url.Parse(databaseURL)
	if err != nil {
		t.Fatal(err)
	}
	query := parsed.Query()
	query.Set("search_path", schema)
	parsed.RawQuery = query.Encode()
	testDB, err := sql.Open("pgx", parsed.String())
	if err != nil {
		t.Fatal(err)
	}
	defer testDB.Close()
	previousDB := db
	db = testDB
	defer func() { db = previousDB }()
	migration, err := os.ReadFile("migrations/001_initial_schema.up.sql")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(string(migration)); err != nil {
		t.Fatal(err)
	}
	server := httptest.NewServer(newRouter())
	defer server.Close()
	newClient := func() *http.Client { jar, _ := cookiejar.New(nil); return &http.Client{Jar: jar} }
	request := func(client *http.Client, method, path string, body any, status int, result any) {
		t.Helper()
		payload, err := json.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}
		req, err := http.NewRequest(method, server.URL+path, bytes.NewReader(payload))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		response, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer response.Body.Close()
		content, _ := io.ReadAll(response.Body)
		if response.StatusCode != status {
			t.Fatalf("%s %s: expected %d, got %d: %s", method, path, status, response.StatusCode, content)
		}
		if result != nil {
			if err := json.Unmarshal(content, result); err != nil {
				t.Fatalf("decode %s: %v (%s)", path, err, content)
			}
		}
	}
	alice, bob := newClient(), newClient()
	request(alice, "GET", "/workouts", nil, 401, nil)
	request(alice, "GET", "/login", nil, 405, nil)
	for name, client := range map[string]*http.Client{"alice": alice, "bob": bob} {
		credentials := map[string]string{"username": name, "password": "integration-password"}
		request(client, "POST", "/register", credentials, 201, nil)
		request(client, "POST", "/login", credentials, 200, nil)
	}
	request(newClient(), "POST", "/login", map[string]string{"username": "alice", "password": "wrong"}, 401, nil)
	var exercise Exercise
	var program Program
	var workout Workout
	request(alice, "POST", "/exercises", map[string]string{"name": "Bench press"}, 201, &exercise)
	request(alice, "PUT", fmt.Sprintf("/exercises?id=%d", exercise.ID), map[string]string{"name": "Barbell bench press"}, 200, &exercise)
	request(alice, "POST", "/programs", map[string]string{"name": "Push day"}, 201, &program)
	request(alice, "PUT", fmt.Sprintf("/programs?id=%d", program.ID), map[string]string{"name": "Upper body"}, 200, &program)
	var pe ProgramExercise
	request(alice, "POST", "/program-exercises", map[string]int{"program_id": program.ID, "exercise_id": exercise.ID}, 201, &pe)
	var programs []ProgramResponse
	request(alice, "GET", "/programs", nil, 200, &programs)
	if len(programs) != 1 {
		t.Fatalf("expected one program, got %d", len(programs))
	}
	var programDetail ProgramResponse
	request(alice, "GET", fmt.Sprintf("/programs?id=%d", program.ID), nil, 200, &programDetail)
	if len(programDetail.Exercises) != 1 || programDetail.Exercises[0].ID != pe.ID {
		t.Fatalf("missing program association ID: %+v", programDetail)
	}
	request(alice, "POST", "/workouts", map[string]any{"name": "Monday", "date": "invalid"}, 400, nil)
	request(bob, "POST", "/workouts", map[string]any{"name": "Forbidden", "date": "2026-11-09", "program_id": program.ID}, 404, nil)
	request(alice, "POST", "/workouts", map[string]any{"date": "2026-11-09", "program_id": program.ID}, 201, &workout)
	var we WorkoutExercise
	request(alice, "POST", "/workout-exercises", map[string]int{"workout_id": workout.ID, "exercise_id": exercise.ID}, 201, &we)
	var set ExerciseSet
	request(alice, "POST", "/exercise-sets", map[string]any{"workout_exercise_id": we.ID, "reps": 10, "weight": 42.5}, 201, &set)
	if set.SetNumber != 1 {
		t.Fatalf("expected set number 1, got %d", set.SetNumber)
	}
	var blankSet ExerciseSet
	request(alice, "POST", "/exercise-sets", map[string]any{"workout_exercise_id": we.ID, "reps": nil, "weight": nil}, 201, &blankSet)
	if blankSet.SetNumber != 2 {
		t.Fatalf("expected set number 2, got %d", blankSet.SetNumber)
	}
	var detail WorkoutResponse
	request(alice, "GET", fmt.Sprintf("/workouts?id=%d", workout.ID), nil, 200, &detail)
	if detail.ProgramID == nil || *detail.ProgramID != program.ID || len(detail.Exercises) != 1 || len(detail.Exercises[0].Sets) != 2 {
		t.Fatalf("incomplete workout: %+v", detail)
	}
	if *detail.Exercises[0].Sets[0].Weight != 42.5 || *detail.Exercises[0].Sets[0].Reps != 10 {
		t.Fatal("saved set did not persist")
	}
	request(bob, "GET", fmt.Sprintf("/workouts?id=%d", workout.ID), nil, 404, nil)
	request(bob, "GET", fmt.Sprintf("/programs?id=%d", program.ID), nil, 404, nil)
	request(bob, "POST", "/exercise-sets", map[string]int{"workout_exercise_id": we.ID, "reps": 1}, 404, nil)
	request(bob, "DELETE", fmt.Sprintf("/exercise-sets?id=%d", set.ID), nil, 404, nil)
	var otherWorkouts []WorkoutResponse
	request(bob, "GET", "/workouts", nil, 200, &otherWorkouts)
	if len(otherWorkouts) != 0 {
		t.Fatal("another user's workouts were exposed")
	}
	request(alice, "DELETE", fmt.Sprintf("/exercise-sets?id=%d", set.ID), nil, 204, nil)
	request(alice, "DELETE", fmt.Sprintf("/workout-exercises?id=%d", we.ID), nil, 204, nil)
	request(alice, "PUT", fmt.Sprintf("/workouts?id=%d", workout.ID), map[string]any{"date": "2026-11-10", "program_id": program.ID}, 200, nil)
	request(alice, "DELETE", fmt.Sprintf("/program-exercises?id=%d", pe.ID), nil, 204, nil)
	request(alice, "DELETE", fmt.Sprintf("/programs?id=%d", program.ID), nil, 204, nil)
	request(alice, "GET", fmt.Sprintf("/workouts?id=%d", workout.ID), nil, 200, &detail)
	if detail.ProgramID != nil || detail.ProgramName != nil || detail.Date != "2026-11-10" {
		t.Fatalf("program deletion or workout update failed: %+v", detail)
	}
	request(alice, "DELETE", fmt.Sprintf("/workouts?id=%d", workout.ID), nil, 204, nil)
	request(alice, "GET", fmt.Sprintf("/workouts?id=%d", workout.ID), nil, 404, nil)
	request(alice, "DELETE", fmt.Sprintf("/exercises?id=%d", exercise.ID), nil, 204, nil)
}

