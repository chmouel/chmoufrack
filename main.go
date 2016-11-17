package main

import (
	"fmt"
	"log"
)

var (
	TRACK_LENGTH = 400
)

var VMA = []int{13, 14, 15, 16, 17, 18, 19}

func main() {
	var rounds = []Workout{
		Workout{
			Repetition:  "1",
			Meters:      "1000",
			Percentage:  "90",
			TrackLength: TRACK_LENGTH,
			Repos:       "200m active",
		},
		Workout{
			Repetition:  "1",
			Meters:      "800",
			Percentage:  "95",
			TrackLength: TRACK_LENGTH,
			Repos:       "200m active",
		},
		Workout{
			Repetition:  "1",
			Meters:      "600",
			Percentage:  "100",
			TrackLength: TRACK_LENGTH,
			Repos:       "2mn arret",
		},
		Workout{
			Repetition:  "1",
			Meters:      "400",
			Percentage:  "105",
			TrackLength: TRACK_LENGTH,
			Repos:       "200m active",
		},
		Workout{
			Repetition:  "1",
			Meters:      "200",
			Percentage:  "110",
			TrackLength: TRACK_LENGTH,
			Repos:       "200m active",
		},
	}

	db, err := createTable()
	if err != nil {
		log.Fatal(err)
	}

	err = generate_html(rounds)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(db)
	fmt.Println(rounds)
}

// Local Variables:
// compile-command: "go run main.go generate_html.go calcul.go sql.go";
// End:
