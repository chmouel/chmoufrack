package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
)

func main() {
	var rounds = []Workout{}
	user, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	listP := flag.Bool("listP", false, "List all programs")
	listW := flag.Bool("listW", false, "List all workouts")
	createP := flag.Bool("createP", false, "Create Program: PROGRAM_NAME [COMMENT]")
	createW := flag.Bool("createW", false, "Create workout: REPETITION METERS PERCENTAGE REPOS")
	assignWP := flag.Bool("assignWP", false, "Assign: WORKOUT PROGRAM")
	deleteP := flag.Bool("deleteP", false, "Create Program: PROGRAM_NAME")
	deleteW := flag.Bool("deleteW", false, "Delete workout: WORKOUT_NAME")
	populateSample := flag.Bool("populateS", false, "Populate samples")
	outputFile := flag.String("o", "", "Output file for the generated HTML")
	configDir := flag.String("configdir", filepath.Join(user.HomeDir, ".config/frack"), "Config directory for database")
	trackLength := flag.Int("trackLength", TRACK_LENGTH, "Track Length")
	yamlSource := flag.String("y", "", "Use a yaml file as source instead of the DB")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of frack: PROGRAM\n\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	TRACK_LENGTH = *trackLength

	CONFIG_DIR = *configDir
	if _, err := os.Stat(*configDir); os.IsNotExist(err) {
		err := os.MkdirAll(*configDir, 0755)
		if err != nil {
			log.Fatal(err)
		}
	}

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
	} else if *assignWP {
		if flag.Arg(0) == "" || flag.Arg(1) == "" {
			log.Fatal("assignWP take at least two arguments")
		}
		_, err = associateWorkoutProgramByName(flag.Arg(0), flag.Arg(1), db)
		if err != nil {
			log.Fatal(err)
		}
		return
	} else if *populateSample {
		err = createSample(db)
		if err != nil {
			log.Fatal(err)
		}
		return
	} else if *deleteP {
		if flag.Arg(0) == "" {
			log.Fatal("deleteP take at least one argument")
		}
		_, err = deleteProgram(flag.Arg(0), db)
		if err != nil {
			log.Fatal(err)
		}
		return
	} else if *deleteW {
		if flag.Arg(0) == "" {
			log.Fatal("deleteW take at least one argument")
		}
		w, err := getWorkoutByName(flag.Arg(0), db)
		if err != nil {
			log.Fatal(err)
		}
		_, err = deleteWorkout(w.ID, db)
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

	if *yamlSource != "" {
		rounds, err = yamlImport(program_name, *yamlSource)
	} else {
		rounds, err = getWorkoutsforProgram(program_name, db)
	}
	if err != nil {
		log.Fatal(err)
	}
	if len(rounds) == 0 {
		fmt.Println("No program or workouts associated with this program", program_name)
		os.Exit(1)
	}

	err = generate_html(program_name, rounds, outputWriter)
	if err != nil {
		log.Fatal(err)
	}
}
