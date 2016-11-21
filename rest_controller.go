package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

func RESTProgramsIndex(w http.ResponseWriter, r *http.Request) {
	programs, err := getPrograms()
	if err != nil {
		panic(err)
	}
	if err := json.NewEncoder(w).Encode(programs); err != nil {
		panic(err)
	}
}

func RESTProgramCreate(writer http.ResponseWriter, reader *http.Request) {
	var program Program
	if reader.Body == nil {
		http.Error(writer, "Please send a request body", http.StatusBadRequest)
		return
	}
	err := json.NewDecoder(reader.Body).Decode(&program)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	if _, err := createProgram(program.Name, program.Comment); err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
	}
	writer.Header().Set("Content-Type", "application/json; charset=UTF-8")
	writer.WriteHeader(http.StatusCreated)
}

func RESTMultipleWorkoutsCreate(writer http.ResponseWriter, reader *http.Request) {
	var workouts []Workout
	if reader.Body == nil {
		http.Error(writer, "Please send a request body", http.StatusBadRequest)
		return
	}

	err := json.NewDecoder(reader.Body).Decode(&workouts)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	for _, workout := range workouts {
		if err = convertAndCreateWorkout(workout); err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
		}
	}

	writer.Header().Set("Content-Type", "application/json; charset=UTF-8")
	writer.WriteHeader(http.StatusCreated)
}

func convertAndCreateWorkout(w Workout) (err error) {
	var percentage, meters, repetition int

	p, err := getProgram(w.ProgramName)
	if p.ID == 0 {
		return errors.New("Cannot find programName " + w.ProgramName)
	}

	if repetition, err = strconv.Atoi(w.Repetition); err != nil {
		return
	}

	if meters, err = strconv.Atoi(w.Meters); err != nil {
		return
	}

	if percentage, err = strconv.Atoi(w.Percentage); err != nil {
		return
	}

	_, err = createWorkout(repetition, meters, percentage, w.Repos, int(p.ID))
	return
}
