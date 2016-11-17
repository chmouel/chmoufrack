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
	rows, err := db.Query("SELECT W.repetition, W.meters, W.percentage, W.repos FROM Program P, Workout W, ProgramWorkout PW WHERE P.name = $1 AND PW.WorkoutID == W.ID AND PW.ProgramID == P.id", program_name)

	for rows.Next() {
		var w Workout
		err := rows.Scan(&w.Repetition, &w.Meters, &w.Percentage, &w.Repos)
		if err != nil {
			log.Fatal(err)
		}
		rounds = append(rounds, w)
	}
	err = generate_html(rounds)
	if err != nil {
		log.Fatal(err)
	}

}

// Local Variables:
// compile-command: "go run main.go generate_html.go calcul.go sql.go";
// End:
