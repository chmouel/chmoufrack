package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	var rounds = []Workout{}

	listP := flag.Bool("listp", false, "List all programs")
	listW := flag.Bool("listw", false, "List all workouts")
	outputFile := flag.String("output", "", "Output file for the generated HTML")
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

	outputWriter := os.Stdout
	if *outputFile != "" {
		fo, err := os.Create(*outputFile)
		if err != nil {
			log.Fatal(err)
		}
		outputWriter = fo
	}

	if flag.Arg(0) == "" {
		fmt.Println("I need a workout program name to generate for use -listp to list them.")
		os.Exit(1)
	}

	program_name := flag.Arg(0)

	rounds, err = getWorkoutsforProgram(program_name, db)
	if err != nil {
		log.Fatal(err)
	}

	if len(rounds) == 0 {
		fmt.Println("No program or workouts associated with this program", program_name)
		os.Exit(1)
	}

	// fmt.Println(rounds)
	err = generate_html(rounds, outputWriter)
	if err != nil {
		log.Fatal(err)
	}
}

// Local Variables:
// compile-command: "go run main.go generate_html.go calcul.go sql.go";
// End:
