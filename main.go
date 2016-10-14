package main

import "log"

var (
	TRACK_LENGTH = 400
)

var VMA = []int{16, 17, 18}

func main() {
	// converting to float

	var rounds = [][]string{
		[]string{"1", "90", "1000", "1mn"},
		[]string{"1", "95", "800", "200m active"},
		[]string{"1", "100", "600", "3mn"},
		[]string{"1", "110", "400", "3mn"},
		[]string{"1", "120", "200", "3mn"},
	}

	err := generate_html(rounds)
	if err != nil {
		log.Fatal(err)
	}
}

// Local Variables:
// compile-command: "go run main.go generate_html.go calcul.go";
// End:
