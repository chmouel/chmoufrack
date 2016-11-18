package main

import (
	"flag"
	"log"
)

func main() {
	var rounds = []Workout{}

	listP := flag.Bool("list", false, "list all programs")
	flag.Parse()

	db, err := createSchema()
	if err != nil {
		log.Fatal(err)
	}

	if *listP {
		err := ListAllPrograms(db)
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	program_name := "Pyramidal"
	rounds, err = getWorkoutsforProgram(program_name, db)
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Println(rounds)
	err = generate_html(rounds)
	if err != nil {
		log.Fatal(err)
	}
}

// Local Variables:
// compile-command: "go run main.go generate_html.go calcul.go sql.go";
// End:
