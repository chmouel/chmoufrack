package main

import (
	"flag"
	"fmt"
	"log"
	"os"
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

	if flag.Arg(0) == "" {
		fmt.Println("I need a program name to generate for use -listp to list them.")
		os.Exit(1)
	}

	program_name := flag.Arg(0)

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
