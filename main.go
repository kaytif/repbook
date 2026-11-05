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

	http.Handle("/workouts", authMiddleware(http.HandlerFunc(handlerWorkout)))
	http.Handle("/exercises", authMiddleware(http.HandlerFunc(handlerExercise)))
	http.Handle("/programs", authMiddleware(http.HandlerFunc(handlerExercise)))
	http.Handle("/workout-exercises", authMiddleware(http.HandlerFunc(handlerExercise)))
	http.Handle("/program-exercises", authMiddleware(http.HandlerFunc(handlerExercise)))
	http.Handle("/exercise-sets", authMiddleware(http.HandlerFunc(handlerExercise)))

	log.Println("Server running on port", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
