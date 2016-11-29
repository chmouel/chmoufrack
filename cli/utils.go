package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/chmouel/chmoufrack"
)

// ListAllPrograms ...
func listAllPrograms() (err error) {
	programs, err := chmoufrack.GetPrograms()
	if err != nil {
		return
	}
	for p := range programs {
		var rounds []chmoufrack.Workout
		t := programs[p]
		rounds, err = chmoufrack.GetWorkoutsforProgram(t.Name)
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

func listAllWorkouts() (err error) {
	workouts, err := chmoufrack.GetWorkouts()
	if err != nil {
		return
	}
	for p := range workouts {
		t := workouts[p]
		//TODO: Add programName
		fmt.Printf("%sx%s@%s -- rest %s \n", t.Repetition, t.Meters, t.Percentage, t.Repos)
	}
	return
}

// CreateWorkout ...
func cliCreateWorkout(arg func(int) string) (err error) {
	programName := arg(0)
	p, err := chmoufrack.GetProgram(programName)
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
	_, err = chmoufrack.CreateWorkout(repetition, meters, percentage, repos, int(p.ID))
	if err != nil {
		log.Fatal(err)
	}
	return

}
