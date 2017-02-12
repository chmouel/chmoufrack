package main

import (
	"log"

	"github.com/chmouel/chmoufrack/server"
)

func main() {
	err := server.DBConnect()
	if err != nil {
		log.Fatal(err)
	}
	// e, err := server.GetProgram("WU5k/3x1000/WD5k")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// f, err := os.Create("/tmp/debug")
	// defer f.Close()
	// spew.Fdump(f, e)

	server.Serve()
}
