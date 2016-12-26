package chmoufrack

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// Nested struct and yaml is not working
type yamlProgram struct {
	Name     string
	Comment  string
	Workouts []Workout
}
type yamlStruct struct {
	Program []yamlProgram
}

func YAMLImport(program_name, filename string) (rounds []Workout, err error) {
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
		rounds = append(rounds, program.Workouts...)
	}
	return
}
