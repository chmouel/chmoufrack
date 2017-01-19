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
	staticDir := flag.String("s", common.StaticDir, "Location for the static html files")
	flag.Parse()

	if _, err := os.Stat(*yamlFile); os.IsNotExist(err) {
		log.Fatal("Cannot found: " + *yamlFile)
	}
	common.Yaml_File = *yamlFile

	if _, err := os.Stat(*staticDir); os.IsNotExist(err) {
		log.Fatal("Cannot found: " + *staticDir)
	}
	common.StaticDir = *staticDir

	server.Serve()
}
