package main

import (
	"io/ioutil"
	"log"

	"github.com/chmouel/chmoufrack/server"

	yaml "gopkg.in/yaml.v2"
)

func YAMLImport(filename string) (err error) {
	var exercises []server.Exercise

	source, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}

	err = yaml.Unmarshal(source, &exercises)
	if err != nil {
		return
	}

	for k, e := range exercises {
		e.ID = k
		_, err := server.AddExercise(e)
		if err != nil {
			log.Fatal(err)
		}
	}
	return
}
