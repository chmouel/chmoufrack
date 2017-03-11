package server

import (
	"fmt"
	"io/ioutil"
	"log"

	yaml "gopkg.in/yaml.v2"
)

func YAMLExport() (err error) {
	exercises, err := getAllExercises()
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

func YAMLImport(filename string) (err error) {
	var exercises []Exercise

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
		_, err := addExercise(e)
		if err != nil {
			log.Fatal(err)
		}
	}
	return
}
