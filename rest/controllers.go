package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/chmouel/chmoufrack"
	"github.com/gorilla/mux"
)

func GETPrograms(writer http.ResponseWriter, r *http.Request) {
	var programs []chmoufrack.Program
	var err error

	if programs, err = chmoufrack.GetPrograms(); err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
	}

	if err = json.NewEncoder(writer).Encode(programs); err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
	}
}

func GetWorkoutsForProgram(writer http.ResponseWriter, reader *http.Request) {
	vars := mux.Vars(reader)
	programName := vars["name"]

	workouts, err := chmoufrack.GetWorkoutsforProgram(programName)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	if err := json.NewEncoder(writer).Encode(workouts); err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
	}
}

func CreateProgram(writer http.ResponseWriter, reader *http.Request) {
	var program chmoufrack.Program
	if reader.Body == nil {
		http.Error(writer, "Please send a request body", http.StatusBadRequest)
		return
	}
	err := json.NewDecoder(reader.Body).Decode(&program)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	if _, err := chmoufrack.CreateProgram(program.Name, program.Comment); err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	writer.Header().Set("Content-Type", "application/json; charset=UTF-8")
	writer.WriteHeader(http.StatusCreated)
}

func CreateMultipleWorkouts(writer http.ResponseWriter, reader *http.Request) {
	vars := mux.Vars(reader)
	programName := vars["name"]

	var workouts []chmoufrack.Workout
	if reader.Body == nil {
		http.Error(writer, "Please send a request body", http.StatusBadRequest)
		fmt.Println("Please send a request body")
		return
	}

	err := json.NewDecoder(reader.Body).Decode(&workouts)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		fmt.Println("Cannot create Workout")
		return
	}

	for _, workout := range workouts {
		if err = convertAndCreateWorkout(programName, workout); err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			fmt.Println(err)
			return
		}
	}

	writer.Header().Set("Content-Type", "application/json; charset=UTF-8")
	writer.WriteHeader(http.StatusCreated)
}

//RESTProgramCleanup Cleanup all workout of a program
func CleanupProgram(writer http.ResponseWriter, reader *http.Request) {
	vars := mux.Vars(reader)
	programName := vars["name"]

	p, err := chmoufrack.GetProgram(programName)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = chmoufrack.DeleteAllWorkoutProgram(p.ID)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	writer.Header().Set("Content-Type", "application/json; charset=UTF-8")
	writer.WriteHeader(http.StatusCreated)
}

func convertAndCreateWorkout(ProgramName string, w chmoufrack.Workout) (err error) {
	var percentage, meters, repetition int

	p, err := chmoufrack.GetProgram(ProgramName)
	if p.ID == 0 {
		return errors.New("Cannot find programName " + ProgramName)
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

	_, err = chmoufrack.CreateWorkout(repetition, meters, percentage, w.Repos, int(p.ID))
	return
}
