package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
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
	createW := flag.Bool("createW", false, "Create workout for program: PROGRAM_NAME REPETITION METERS PERCENTAGE REPOS")
	deleteP := flag.Bool("deleteP", false, "Create Program: PROGRAM_NAME")
	deleteW := flag.Bool("deleteW", false, "Delete workout attached to program: PROGRAM_NAME WORKOUT_NAME")
	populateSample := flag.Bool("populateS", false, "Populate samples")
	outputFile := flag.String("o", "", "Output file for the generated HTML")
	configDir := flag.String("configdir", filepath.Join(user.HomeDir, ".config/frack"), "Config directory for database")
	vmas := flag.String("v", VMA, "Set VMAS with a colon as delimter in between")
	trackLength := flag.Int("trackLength", TRACK_LENGTH, "Track Length")
	yamlSource := flag.String("y", "", "Use a yaml file as source instead of the DB")
	rest := flag.Bool("rest", false, "Start the REST server")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of frack: PROGRAM\n\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	TRACK_LENGTH = *trackLength
	VMA = *vmas

	CONFIG_DIR = *configDir
	if _, err := os.Stat(*configDir); os.IsNotExist(err) {
		err := os.MkdirAll(*configDir, 0755)
		if err != nil {
			log.Fatal(err)
		}
	}

	err = createSchema()
	if err != nil {
		log.Fatal(err)
	}

	if *rest {
		servRest()
		return
	} else if *listP {
		err = ListAllPrograms()
		if err != nil {
			log.Fatal(err)
		}
		return
	} else if *listW {
		err = ListAllWorkouts()
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

		_, err := createProgram(program, comment)
		if err != nil {
			log.Fatal(err)
		}
		return
	} else if *createW {
		if flag.Arg(0) == "" {
			log.Fatal("createW take at least one argument")
		}
		err := CreateWorkout(flag.Arg)
		if err != nil {
			log.Fatal(err)
		}
		return
	} else if *populateSample {
		err = createSample()
		if err != nil {
			log.Fatal(err)
		}
		return
	} else if *deleteP {
		if flag.Arg(0) == "" {
			log.Fatal("deleteP take at least one argument")
		}
		_, err = deleteProgram(flag.Arg(0))
		if err != nil {
			log.Fatal(err)
		}
		return
	} else if *deleteW {
		// Delete Workout of Program
		if flag.Arg(0) == "" {
			log.Fatal("deleteW take at least one argument")
		}
		p, err := getProgram(flag.Arg(0))
		if p.ID == 0 {
			log.Fatal("Could not find " + flag.Arg(0))
		}

		w, err := getWorkoutByName(flag.Arg(1))
		if err != nil {
			log.Fatal(err)
		}
		_, err = deleteWorkout(p.ID, w.ID)
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
		rounds, err = getWorkoutsforProgram(program_name)
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
