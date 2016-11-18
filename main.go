package main

import (
	"flag"
	"log"
)

func main() {
	var rounds = []Workout{}

	listP := flag.Bool("listp", false, "list all programs")
	listW := flag.Bool("listw", false, "list all workouts")
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
	} else if *listW {
		err := ListAllWorkouts(db)
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
