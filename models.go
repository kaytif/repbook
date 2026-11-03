package main

import (
	"time"
)

type Workout struct {
	ID        int        `json:"id"`
	Name      string     `json:"name"`
	Date      string     `json:"date"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
	ProgramID *int       `json:"program_id"`
}

type Program struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Exercise struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type WorkoutExercise struct {
	ID         int `json:"id"`
	WorkoutID  int `json:"workout_id"`
	ExerciseID int `json:"exercise_id"`
}

type ProgramExercise struct {
	ID         int `json:"id"`
	ProgramID  int `json:"program_id"`
	ExerciseID int `json:"exercise_id"`
}

type WorkoutExerciseResponse struct {
	ID         int           `json:"id"`
	ExerciseID int           `json:"exercise_id"`
	Name       string        `json:"name"`
	Sets       []ExerciseSet `json:"sets"`
}

type WorkoutResponse struct {
	ID          int                       `json:"id"`
	Name        string                    `json:"name"`
	Date        string                    `json:"date"`
	ProgramName *string                   `json:"program_name"`
	Exercises   []WorkoutExerciseResponse `json:"exercises"`
}

type ProgramExerciseResponse struct {
	ExerciseID int    `json:"exercise_id"`
	Name       string `json:"name"`
}

type ProgramResponse struct {
	ID        int                       `json:"id"`
	Name      string                    `json:"name"`
	Exercises []ProgramExerciseResponse `json:"exercises"`
}

type ExerciseSet struct {
	ID                int      `json:"id"`
	WorkoutExerciseID int      `json:"workout_exercise_id"`
	SetNumber         int      `json:"set_number"`
	Reps              *int     `json:"reps"`
	Weight            *float64 `json:"weight"`
}
