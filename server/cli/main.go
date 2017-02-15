package main

import (
	"flag"
	"fmt"
	"log"

	"os"

	"github.com/chmouel/chmoufrack/server"
)

func main() {
	_location := os.Getenv("FRACK_DB")
	if _location == "" {
		_location = "./frack.db"
	}
	_initDB := false
	if os.Getenv("FRACK_INIT_DB") != "" {
		_initDB = true
	}
	dblocation := flag.String("dblocation", _location, "sqlite db location")
	initDBbool := flag.Bool("initDB", _initDB, "init DB with samples DATA")

	flag.Parse()

	fmt.Printf("Using DB from %s\n", *dblocation)
	err := server.DBConnect(*dblocation)
	if *initDBbool {
		fmt.Println("InitDB")
		err := server.InitFixturesDB()
		if err != nil {
			log.Fatal(err)
		}
	}

	if err != nil {
		fmt.Printf("repeat error: %s\n", err.Error())
		return
	}

	server.Serve()
}
