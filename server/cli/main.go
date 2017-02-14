package main

import (
	"flag"
	"fmt"

	"os"

	"github.com/chmouel/chmoufrack/server"
)

func main() {
	_location := os.Getenv("FRACK_DB")
	if _location == "" {
		_location = "./frack.db"
	}
	dblocation := flag.String("dblocation", _location, "sqlite db location")
	flag.Parse()

	fmt.Printf("Using DB from %s\n", *dblocation)
	err := server.DBConnect(*dblocation)
	if err != nil {
		fmt.Printf("repeat error: %s\n", err.Error())
		return
	}

	server.Serve()
}
