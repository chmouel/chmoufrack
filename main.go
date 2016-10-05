package main

import (
	"fmt"
	"log"
	"strconv"
)

var (
	TRACK_LENGTH = 400
)

var VMA = []int{16, 17, 18}

func print_separator(s string) {
	for i := 0; i < len(s); i++ {
		fmt.Print("-")
	}
	fmt.Println("")
}

// print_repeat ...
func print_repeat(repetion, percentage, track_laps, meters float64, repos string) {
	header := fmt.Sprintf("%d * %dm at %d%% / Total laps %d\n",
		int(repetion), int(meters), int(percentage), int(track_laps))

	print_separator(header)
	fmt.Print(header)
	print_separator(header)

	for _, vmad := range VMA {
		vma := float64(vmad)
		total_time, err := calcul_vma_distance(vma, percentage, meters)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("VMA: %0.f => Total: %s", vma, total_time)
		if int(meters) >= TRACK_LENGTH {
			time_track, err := calcul_vma_distance(vma, percentage, float64(TRACK_LENGTH))
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf(" / %d: %s", TRACK_LENGTH, time_track)
		}

		speed := calcul_vma_speed(vma, percentage)
		fmt.Printf(" / Speed: %.1fkh / Pace: %s\n", speed, calcul_pace(speed))
	}
	fmt.Printf("\nRepos: %s\n", repos)
	fmt.Println("")
}

func main() {
	// converting to float

	var rounds = [][]string{
		[]string{"1", "90", "1000", "1mn de repos"},
		[]string{"1", "95", "800", "3mn de repos"},
		[]string{"1", "100", "600", "3mn de repos"},
		[]string{"1", "110", "400", "3mn de repos"},
		[]string{"1", "120", "200", "3mn de repos"},
	}
	for i := range rounds {
		repetition, err := strconv.ParseFloat(rounds[i][0], 64)
		if err != nil {
			log.Fatal(err)
		}

		percentage, err := strconv.ParseFloat(rounds[i][1], 64)
		if err != nil {
			log.Fatal(err)
		}

		meters, err := strconv.ParseFloat(rounds[i][2], 64)
		if err != nil {
			log.Fatal(err)
		}
		track_length, err := strconv.ParseFloat(strconv.Itoa(TRACK_LENGTH), 64)
		if err != nil {
			log.Fatal(err)
		}
		track_laps := meters / track_length

		repos := rounds[i][3]

		print_repeat(repetition, percentage, track_laps, meters, repos)
	}
}
