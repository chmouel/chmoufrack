package main

import (
	"fmt"
	"log"
	"strconv"
)

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
		fmt.Printf("%-15s - %sx%s@%s -- rest %s \n", t.ProgramName, t.Repetition, t.Meters, t.Percentage, t.Repos)
	}
	return
}

// CreateWorkout ...
func CreateWorkout(arg func(int) string) (err error) {
	programName := arg(0)
	p, err := getProgram(programName)
	//TODO: custom error
	if p.ID == 0 {
		log.Fatal("Cannot find workout " + programName)
	}

	repetition, err := strconv.Atoi(arg(1))
	if err != nil {
		return
	}

	meters, err := strconv.Atoi(arg(2))
	if err != nil {
		return
	}

	percentage, err := strconv.Atoi(arg(3))
	if err != nil {
		return
	}

	repos := arg(4)

	//TODO
	_, err = createWorkout(repetition, meters, percentage, repos, int(p.ID))
	if err != nil {
		log.Fatal(err)
	}
	return

}
