package main

import "fmt"

// ListAllPrograms ...
func ListAllPrograms() (err error) {
	programs, err := getPrograms()
	if err != nil {
		return
	}
	for p := range programs {
		var rounds []Workout
		t := programs[p]
		rounds, err = getWorkoutsforProgram(t.Name)
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

func ListAllWorkouts() (err error) {
	workouts, err := getWorkouts()
	if err != nil {
		return
	}
	for p := range workouts {
		t := workouts[p]
		fmt.Printf("%sx%s@%s -- rest %s \n", t.Repetition, t.Meters, t.Percentage, t.Repos)
	}
	return
}
