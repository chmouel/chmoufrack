package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
)

func main() {
	var rounds = []Workout{}

	listP := flag.Bool("listP", false, "List all programs")
	listW := flag.Bool("listW", false, "List all workouts")
	createP := flag.Bool("createP", false, "Create Program: PROGRAM_NAME [COMMENT]")
	createW := flag.Bool("createW", false, "Create workout: REPETITION METERS PERCENTAGE REPOS")
	outputFile := flag.String("o", "", "Output file for the generated HTML")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of frack: PROGRAM\n\n")
		flag.PrintDefaults()
	}

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
	} else if *createP {
		if flag.Arg(0) == "" {
			log.Fatal("createP take at least one argument")
		}
		program := flag.Arg(0)
		comment := flag.Arg(1)

		_, err := createProgram(program, comment, db)
		if err != nil {
			log.Fatal(err)
		}
		return
	} else if *createW {
		if flag.Arg(0) == "" {
			log.Fatal("createW take at least one argument")
		}
		repetition, err := strconv.Atoi(flag.Arg(0))
		if err != nil {
			log.Fatal(err)
		}
		meters, err := strconv.Atoi(flag.Arg(1))
		if err != nil {
			log.Fatal(err)
		}
		percentage, err := strconv.Atoi(flag.Arg(2))
		if err != nil {
			log.Fatal(err)
		}
		repos := flag.Arg(3)

		_, err = createWorkout(repetition, meters, percentage, repos, db)
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	outputWriter := os.Stdout
	if *outputFile != "" {
		fileOutput, err := os.Create(*outputFile)
		if err != nil {
			log.Fatal(err)
		}
		outputWriter = fileOutput
	}

	if flag.Arg(0) == "" {
		fmt.Println("I need a workout program name to generate for use -listP to list them.")
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

	err = generate_html(rounds, outputWriter)
	if err != nil {
		log.Fatal(err)
	}
}
