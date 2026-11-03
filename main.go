package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	connectDB()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/workouts", handlerWorkout)
	http.HandleFunc("/exercises", handlerExercise)
	http.HandleFunc("/programs", handlerProgram)
	http.HandleFunc("/workout-exercises", handlerWorkoutExercise)
	http.HandleFunc("/program-exercises", handlerProgramExercise)
	http.HandleFunc("/exercise-sets", handlerExerciseSet)

	log.Println("Server running on port", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
