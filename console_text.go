package main

import (
	"database/sql"
	"fmt"
)

// ListAllPrograms ...
func ListAllPrograms(db *sql.DB) (err error) {
	programs, err := getPrograms(db)
	if err != nil {
		return
	}
	for p := range programs {
		var rounds []Workout
		t := programs[p]
		rounds, err = getWorkoutsforProgram(t.Name, db)
		if err != nil {
			return
		}
		fmt.Printf("%-15s | ", t.Name)
		for w := range rounds {
			tt := rounds[w]
			fmt.Printf("%sx%s@%s ", tt.Repetition, tt.Meters, tt.Percentage)
		}
		fmt.Printf("\n")
	}
	return
}

func ListAllWorkouts(db *sql.DB) (err error) {
	workouts, err := getWorkouts(db)
	if err != nil {
		return
	}
	for p := range workouts {
		t := workouts[p]
		fmt.Printf("%sx%s@%s -- rest %s \n", t.Repetition, t.Meters, t.Percentage, t.Repos)
	}
	return
}
