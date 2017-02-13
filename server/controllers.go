package server

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func GETExcercise(writer http.ResponseWriter, reader *http.Request) {
	var excercise Excercise
	var err error

	vars := mux.Vars(reader)
	programName := vars["name"]

	excercise, err = getProgram(programName)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	if err = json.NewEncoder(writer).Encode(excercise); err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
	}
}

func GETExcerciseS(writer http.ResponseWriter, reader *http.Request) {
	var excercise []Excercise
	var err error

	excercise, err = getAllPrograms()
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	if err = json.NewEncoder(writer).Encode(excercise); err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
	}
}
