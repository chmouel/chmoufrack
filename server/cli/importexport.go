package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/chmouel/chmoufrack/server"

	yaml "gopkg.in/yaml.v2"
)

func yAMLExport() (err error) {
	exercises, err := server.GetAllExercises()
	if err != nil {
		return
	}

	d, err := yaml.Marshal(exercises)
	if err != nil {
		return
	}
	fmt.Println(string(d))
	return
}

func yAMLImport(filename string) (err error) {
	var exercises []server.Exercise

	server.ACL = false

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
