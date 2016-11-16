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
	// converting to float

	var rounds = [][]string{
		[]string{"1", "90", "1000", "200m active"},
		[]string{"1", "95", "800", "200m active"},
		[]string{"1", "100", "600", "2mn"},
		[]string{"1", "110", "400", "2mn"},
		[]string{"1", "120", "200", "2mn"},
	}

	// rounds = [][]string{
	// 	[]string{"8", "110", "100", "1mn30"},
	// 	[]string{"6", "100", "200", "1mn30"},
	// 	[]string{"4", "90", "400", "1mn30"},
	// }
	//err := generate_html(rounds)

	db, err := createTable()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(db)
	fmt.Println(rounds)
}

// Local Variables:
// compile-command: "go run main.go generate_html.go calcul.go sql.go";
// End:
