package server

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// POSTExcercise ...
func POSTExcercise(writer http.ResponseWriter, reader *http.Request) {
	var excercise Excercise
	var err error

	if reader.Body == nil {
		http.Error(writer, "Please send a request body", http.StatusBadRequest)
		return
	}

	err = json.NewDecoder(reader.Body).Decode(&excercise)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	_ = addProgram(excercise)

	writer.Header().Set("Content-Type", "application/json; charset=UTF-8")
	writer.WriteHeader(http.StatusCreated)
	return
}

func GETExcercise(writer http.ResponseWriter, reader *http.Request) {
	var excercise Excercise
	var err error

	vars := mux.Vars(reader)
	excerciseName := vars["name"]

	excercise, err = getProgram(excerciseName)
	if err != nil {
		if _, ok := err.(*error404); ok {
			http.Error(writer, err.Error(), http.StatusNotFound)
		} else {
			http.Error(writer, err.Error(), http.StatusBadRequest)
		}
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
