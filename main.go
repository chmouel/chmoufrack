package main

import "log"

var VMA = []int{13, 14, 15, 16, 17, 18, 19}

func main() {
	var rounds = []Workout{}

	db, err := createTable()
	if err != nil {
		log.Fatal(err)
	}

	program_name := "Pyramidal"
	rounds, err = getWorkouts(program_name, db)
	if err != nil {
		log.Fatal(err)
	}

	err = generate_html(rounds)
	if err != nil {
		log.Fatal(err)
	}
}

// Local Variables:
// compile-command: "go run main.go generate_html.go calcul.go sql.go";
// End:
