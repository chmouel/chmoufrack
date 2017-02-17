package server

import (
	"flag"
	"fmt"
	"log"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	temp_db_file := TempFileName(".sqlite")
	db_file := flag.String("db", temp_db_file, "run database integration tests")
	flag.Parse()

	fmt.Printf("Using DB: %s\n", *db_file)

	err := DBConnect(*db_file)
	if err != nil {
		log.Fatal(err)
	}
	code := m.Run()

	if *db_file == temp_db_file {
		err = os.Remove(*db_file)
		if err != nil {
			log.Fatal(err)
		}
	}
	os.Exit(code)
}
