package server

import (
	"log"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	dblocation := os.Getenv("FRACK_TEST_DB")
	if dblocation == "" {
		log.Fatal("You need to specify a FRACK_TEST_DB variable")
	}
	err := DBConnect(dblocation, true)
	if err != nil {
		log.Fatal(err)
	}
	code := m.Run()
	os.Exit(code)
}
