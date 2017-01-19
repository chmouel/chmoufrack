package common

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// YAMLImport import yaml function from filename assuming we have a programname
func YAMLImport(filename string) (err error, programs []Program) {
	source, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}

	err = yaml.Unmarshal(source, &programs)
	if err != nil {
		return
	}

	return
}
