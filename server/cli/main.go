package main

import (
	"flag"
	"fmt"
	"log"

	"os"

	"github.com/chmouel/chmoufrack/server"
)

func main() {
	_staticHTML := os.Getenv("FRACK_STATIC_HTML")
	if _staticHTML == "" {
		_staticHTML = "client"
	}

	_db_location := os.Getenv("FRACK_DB")
	_initDB := false
	if os.Getenv("FRACK_INIT_DB") != "" {
		_initDB = true
	}
	yamlImport := flag.String("yamlImport", "", "Import a yaml file in DB")
	db := flag.String("db", _db_location, "DB Connexion detail")
	initDBbool := flag.Bool("initDB", _initDB, "init DB with samples DATA")
	staticHTML := flag.String("staticHTML", _staticHTML, "client static html location")
	serverPort := flag.Int("port", 8080, "DB Port")

	flag.Parse()

	if *db == "" {
		log.Fatal("You need to specify a MySQL DB cnx with -db")
	}
	err := server.DBConnect(*db, *initDBbool)
	if err != nil {
		log.Fatalf("Cannot connect to mysql: %s", err.Error())
	}

	if *initDBbool {
		fmt.Println("Adding Fixtures to DB")
		err := server.InitFixturesDB()
		if err != nil {
			log.Fatal(err)
		}
	}

	if *yamlImport != "" {
		err := YAMLImport(*yamlImport)
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	server.Serve(*staticHTML, *serverPort)
}
