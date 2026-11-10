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

	log.Println("Server running on port", port)
	log.Fatal(http.ListenAndServe(":"+port, newRouter()))
}

func newRouter() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /login", handlerLogin)
	mux.HandleFunc("POST /register", registerUser)
	mux.Handle("/workouts", authMiddleware(http.HandlerFunc(handlerWorkout)))
	mux.Handle("/exercises", authMiddleware(http.HandlerFunc(handlerExercise)))
	mux.Handle("/programs", authMiddleware(http.HandlerFunc(handlerProgram)))
	mux.Handle("/workout-exercises", authMiddleware(http.HandlerFunc(handlerWorkoutExercise)))
	mux.Handle("/program-exercises", authMiddleware(http.HandlerFunc(handlerProgramExercise)))
	mux.Handle("/exercise-sets", authMiddleware(http.HandlerFunc(handlerExerciseSet)))
	return mux
}
