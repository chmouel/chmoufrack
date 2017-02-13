package server

import (
	"log"
	"os"

	"github.com/davecgh/go-spew/spew"
)

func Debug(a ...interface{}) {
	f, err := os.Create("/tmp/debug")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	spew.Fdump(f, a...)
}
