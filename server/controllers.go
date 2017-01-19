package server

import (
	"encoding/json"
	"net/http"

	c "github.com/chmouel/chmoufrack/common"
	"github.com/gorilla/mux"
)

func GETPrograms(writer http.ResponseWriter, reader *http.Request) {
	var programs []c.Program
	var err error

	if err, programs = c.YAMLImport(c.Yaml_File); err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
	}

	if err = json.NewEncoder(writer).Encode(programs); err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
	}
}

func GETProgram(writer http.ResponseWriter, reader *http.Request) {
	var err error
	var programs []c.Program

	type RestRep struct {
		Name      string
		Comment   string
		Workouts  []c.CalculatedProgram
		TargetVMA string
	}

	vars := mux.Vars(reader)
	programName := vars["name"]
	vma := vars["vma"]
	if vma == "" {
		vma = c.TargetVma
	}
	rr := RestRep{}
	rr.Name = programName
	rr.TargetVMA = vma

	if err, programs = c.YAMLImport(c.Yaml_File); err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
	}

	for _, value := range programs {
		if value.Name == programName {
			rr.Comment = value.Comment

			for _, workout := range value.Workouts {
				var ts c.CalculatedProgram

				if ts, err = c.GenerateProgram(workout, vma); err != nil {
					http.Error(writer, err.Error(), http.StatusBadRequest)
					return
				}
				rr.Workouts = append(rr.Workouts, ts)
			}

			if err = json.NewEncoder(writer).Encode(rr); err != nil {
				http.Error(writer, err.Error(), http.StatusBadRequest)
			}

			break
		}
	}
	//spew.Fdump(writer, rr)
}
