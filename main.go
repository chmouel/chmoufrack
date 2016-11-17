package main

import "log"

func main() {
	var rounds = []Workout{}

	db, err := createSchema()
	if err != nil {
		log.Fatal(err)
	}

	program_name := "Pyramidal"
	rounds, err = getWorkouts(program_name, db)
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
