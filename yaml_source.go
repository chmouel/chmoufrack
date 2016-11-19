package main

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// Nested struct and marshall is not working
type yamlProgram struct {
	Name     string
	Comment  string
	Workouts []Workout
}
type yamlStruct struct {
	Program []yamlProgram
}

func yamlImport(program_name, filename string) (rounds []Workout, err error) {
	t := yamlStruct{}

	source, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}

	err = yaml.Unmarshal(source, &t)
	if err != nil {
		return
	}

	for _, program := range t.Program {
		if program.Name != program_name {
			continue
		}
		for _, workout := range program.Workouts {
			rounds = append(rounds, workout)
		}
	}
	return
}
