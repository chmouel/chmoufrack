package common

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// round float64 - cause the stdlib doesnt have it :(
func round(val float64) float64 {
	var newVal float64

	if math.Abs(val-math.Ceil(val)) <= 0.5 {
		newVal = math.Ceil(val)
	} else {
		newVal = math.Floor(val)
	}
	return newVal
}

// calcul vma of a distance, you give a vma and a percent for a distance
func calculVmaDistance(vma, percent, distance float64) (result string, err error) {
	vmaMs := vma * 1000 / 3600
	vma100 := 100 / vmaMs
	calcul := vma100 / percent * distance

	stemps := int(calcul)
	minute := ((stemps - (stemps)%60) / 60)
	second := ((stemps % 60) * 10) / 10

	if minute > 0 {
		result += strconv.Itoa(minute) + "'"
	}
	result += fmt.Sprintf("%.2d", second)
	if minute == 0 {
		result += "s"
	}

	return
}

// calcul_vma_vitesse ...
func calculVmaSpeed(vma, percent float64) float64 {
	return (vma * percent) / 100
}

// calculPace from a speed (kmh)
func calculPace(vitesse float64) (ret string) {
	var e = 1 / vitesse * 60
	var t = math.Floor(e / 60)
	var n = math.Floor(e - t*60)
	var r = round(60 * (e - t*60 - n))

	if r == 60 {
		n++
		r = 0
	}

	if n == 0 && r != 0 {
		return fmt.Sprintf("%0.f\"", r)
	}

	ret += fmt.Sprintf("%0.f'", n)
	if r == 0 {
		return
	} else if r < 10 {
		ret += "0"
	}

	ret += fmt.Sprintf("%0.f", r)

	return
}

func getVmas(value string) (vmas []int) {
	var s, e int

	if strings.Index(value, ":") > 0 {
		s, _ = strconv.Atoi(strings.Split(value, ":")[0])
		e, _ = strconv.Atoi(strings.Split(value, ":")[1])
	} else {
		s, _ = strconv.Atoi(value)
		e = s
	}

	for i := s; i <= e; i++ {
		vmas = append(vmas, i)
	}

	return
}

// GenerateProgram Generate a program
func GenerateProgram(workout Workout, targetVma string) (ts CalculatedProgram, err error) {
	var totalTime, timeLap string
	vmas := map[string]WorkoutVMA{}

	meters, err := strconv.ParseFloat(workout.Meters, 64)
	if err != nil {
		return
	}

	trackLength, err := strconv.ParseFloat(strconv.Itoa(TrackLength), 64)
	if err != nil {
		return
	}

	percentage, err := strconv.ParseFloat(workout.Percentage, 64)
	if err != nil {
		return
	}

	trackLaps := meters / trackLength
	laps := fmt.Sprintf("%.1f", trackLaps)
	if strings.HasSuffix(laps, ".0") {
		laps = strings.TrimSuffix(laps, ".0")
	} else if laps == "0.5" {
		laps = "½"
	} else if strings.HasSuffix(laps, ".5") {
		laps = strings.Replace(laps, ".5", "½", -1)
	}
	workout.TrackLaps = laps
	workout.TrackLength = TrackLength

	for _, vmad := range getVmas(targetVma) {
		workoutVma := WorkoutVMA{}
		vma := float64(vmad)
		totalTime, err = calculVmaDistance(vma, percentage, meters)
		if err != nil {
			return
		}
		workoutVma.VMA = fmt.Sprintf("%.f", vma)
		workoutVma.TotalTime = totalTime
		if int(meters) >= TrackLength {
			timeLap, err = calculVmaDistance(vma, percentage, float64(TrackLength))
			if err != nil {
				return
			}
			workoutVma.TimeTrack = timeLap
		} else {
			workoutVma.TimeTrack = "NA"
		}
		speed := calculVmaSpeed(vma, percentage)
		workoutVma.Speed = fmt.Sprintf("%.2f", speed)
		workoutVma.Pace = calculPace(speed)

		vmas[workoutVma.VMA] = workoutVma
	}

	ts = CalculatedProgram{
		VMAs:    vmas,
		Workout: workout,
	}

	return
}
