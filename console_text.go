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
		fmt.Printf("%s -- %d workouts\n", t.Name, len(rounds))
	}
	return
}
