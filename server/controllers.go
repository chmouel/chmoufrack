package server

import (
	"encoding/json"
	"net/http"

	c "github.com/chmouel/chmoufrack/common"
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

	if err, programs = c.YAMLImport(c.Yaml_File); err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
	}

	if err = json.NewEncoder(writer).Encode(programs); err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
	}

}
