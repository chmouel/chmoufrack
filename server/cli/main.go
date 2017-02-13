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

	server.Serve()
}
