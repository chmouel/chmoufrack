package server

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"io/ioutil"
	"strings"
)

// POSTExercise ...
func POSTExercise(writer http.ResponseWriter, reader *http.Request) {
	var exercise Exercise
	var err error

	if reader.Body == nil {
		http.Error(writer, "Please send a request body", http.StatusBadRequest)
		return
	}

	x, _ := ioutil.ReadAll(reader.Body)
	err = json.NewDecoder(strings.NewReader(string(x))).Decode(&exercise)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	if exercise.Name == "" { // TODO: proper validation
		http.Error(writer, "Name is invalid in json", http.StatusBadRequest)
		return
	}

	_, err = AddExercise(exercise)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	writer.Header().Set("Content-Type", "application/json; charset=UTF-8")
	writer.WriteHeader(http.StatusCreated)
	return
}

func DeleteExercise(writer http.ResponseWriter, reader *http.Request) {
	var err error
	var i, id int
	vars := mux.Vars(reader)
	exerciseID := vars["id"]
	if i, err = strconv.Atoi(exerciseID); err == nil {
		id = i
	} else {
		id, err = getIdOfExerciseName(exerciseID)

		if err != nil {
			if _, ok := err.(*error404); ok {
				http.Error(writer, err.Error(), http.StatusNotFound)
			} else {
				http.Error(writer, err.Error(), http.StatusBadRequest)
			}
			return
		}
	}

	err = deleteExercise(id)
	if err != nil {
		if _, ok := err.(*error404); ok {
			http.Error(writer, err.Error(), http.StatusNotFound)
		} else {
			http.Error(writer, err.Error(), http.StatusBadRequest)
		}
		return
	}

	writer.Header().Set("Content-Type", "application/json; charset=UTF-8")
	writer.WriteHeader(http.StatusNoContent)
}

func GETExercise(writer http.ResponseWriter, reader *http.Request) {
	var exercise Exercise
	var err error
	var i, id int
	vars := mux.Vars(reader)
	exerciseID := vars["id"]

	if i, err = strconv.Atoi(exerciseID); err == nil {
		id = i
	} else {
		id, err = getIdOfExerciseName(exerciseID)
		if err != nil {
			if _, ok := err.(*error404); ok {
				http.Error(writer, err.Error(), http.StatusNotFound)
			} else {
				http.Error(writer, err.Error(), http.StatusBadRequest)
			}
			return
		}
	}

	exercise, err = getExercise(id)
	if err != nil {
		if _, ok := err.(*error404); ok {
			http.Error(writer, err.Error(), http.StatusNotFound)
		} else {
			http.Error(writer, err.Error(), http.StatusBadRequest)
		}
		return
	}

	if err = json.NewEncoder(writer).Encode(exercise); err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
	}
}

func GETExercises(writer http.ResponseWriter, reader *http.Request) {
	var exercises []Exercise
	var err error

	exercises, err = getAllExercises()
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	if err = json.NewEncoder(writer).Encode(exercises); err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
	}
}
