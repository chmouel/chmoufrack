package main

import (
	"fmt"

	"github.com/chmouel/chmoufrack/server"
)

func main() {
	err := server.DBConnect()
	if err != nil {
		fmt.Printf("repeat error: %s\n", err.Error())
		return
	}

	server.Serve()
}
