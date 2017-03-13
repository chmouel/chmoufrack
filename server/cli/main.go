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
	_initDB := os.Getenv("FRACK_INIT_DB")

	yamlExport := flag.Bool("export", false, "Export a yaml file in DB")
	yamlImport := flag.String("import", "", "Import a yaml file in DB")
	db := flag.String("db", _db_location, "DB Connexion detail")
	initDBbool := flag.String("initDB", _initDB, "init DB with samples DATA, the arg is the facebook id")
	staticHTML := flag.String("staticHTML", _staticHTML, "client static html location")
	serverPort := flag.Int("port", 8080, "DB Port")
	debug := flag.Bool("debug", false, "DB Port")

	flag.Parse()

	if *db == "" {
		log.Fatal("You need to specify a MySQL DB cnx with -db")
	}

	reset := false
	if *initDBbool != "" {
		reset = true
	}

	err := server.DBConnect(*db, reset)
	if err != nil {
		log.Fatalf("Cannot connect to mysql: %s", err.Error())
	}

	if *initDBbool != "" {
		fmt.Println("Adding Fixtures to DB")
		err = server.InitFixturesDB(*initDBbool)
		if err != nil {
			log.Fatal(err)
		}
	}

	if *yamlExport {
		err := server.YAMLExport()
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	if *yamlImport != "" {
		err := server.YAMLImport(*yamlImport)
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	server.Serve(*staticHTML, *serverPort, *debug)
}
