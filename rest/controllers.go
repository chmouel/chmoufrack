package rest

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	c "github.com/chmouel/chmoufrack"
	"github.com/chmouel/chmoufrack/db"
	"github.com/gorilla/mux"
)

func GETPrograms(writer http.ResponseWriter, r *http.Request) {
	var programs []c.Program
	var err error

	if programs, err = db.GetPrograms(); err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
	}

	if err = json.NewEncoder(writer).Encode(programs); err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
	}
}

func GetProgramFull(writer http.ResponseWriter, reader *http.Request) {
	type RestRep struct {
		ProgramName string
		Workouts    []c.TemplateStruct
		TargetVMA   string
	}

	vars := mux.Vars(reader)
	programName := vars["name"]
	vma := vars["vma"]
	if vma == "" {
		vma = c.TargetVma
	}

	var restRep = RestRep{ProgramName: programName, TargetVMA: vma}

	workouts, err := db.GetWorkoutsforProgram(programName)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	for _, workout := range workouts {
		var ts c.TemplateStruct

		if ts, err = c.GenerateProgram(workout, vma); err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		ts.ProgramName = programName
		restRep.Workouts = append(restRep.Workouts, ts)
	}
	if err = json.NewEncoder(writer).Encode(restRep); err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
	}
}

func GetWorkoutsForProgram(writer http.ResponseWriter, reader *http.Request) {
	vars := mux.Vars(reader)
	programName := vars["name"]

	workouts, err := db.GetWorkoutsforProgram(programName)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	if len(workouts) == 0 {
		http.Error(writer, "Workout is empty", http.StatusNotFound)
		return
	}

	if err := json.NewEncoder(writer).Encode(workouts); err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
	}
}

func CreateProgram(writer http.ResponseWriter, reader *http.Request) {
	vars := mux.Vars(reader)
	programName := vars["name"]

	if reader.Body == nil {
		http.Error(writer, "Please send a request body", http.StatusBadRequest)
		return
	}

	if _, err := db.CreateProgram(programName, ""); err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	writer.Header().Set("Content-Type", "application/json; charset=UTF-8")
	writer.WriteHeader(http.StatusCreated)
}

func CreateMultipleWorkouts(writer http.ResponseWriter, reader *http.Request) {
	vars := mux.Vars(reader)
	programName := vars["name"]

	var workouts []c.Workout
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
		if err = convertAndCreateWorkout(programName, workout); err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
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

	p, err := db.GetProgram(programName)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = db.DeleteAllWorkoutProgram(p.ID)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	writer.Header().Set("Content-Type", "application/json; charset=UTF-8")
	writer.WriteHeader(http.StatusCreated)
}

func convertAndCreateWorkout(ProgramName string, w c.Workout) (err error) {
	var percentage, meters, repetition int

	p, err := db.GetProgram(ProgramName)
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

	_, err = db.CreateWorkout(repetition, meters, percentage, w.Repos, int(p.ID))
	return
}
