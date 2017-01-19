package main

import (
	"flag"
	"log"
	"os"

	"github.com/chmouel/chmoufrack/common"
	"github.com/chmouel/chmoufrack/server"
)

func main() {
	yamlFile := flag.String("y", common.Yaml_File, "Location for the yaml file")
	flag.Parse()

	common.Yaml_File = *yamlFile
	if _, err := os.Stat(*yamlFile); os.IsNotExist(err) {
		log.Fatal("Cannot found: " + *yamlFile)
	}
	server.Serve()
}
