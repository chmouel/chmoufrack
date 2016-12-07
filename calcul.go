package chmoufrack

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
func calcul_vma_distance(vma, percent, distance float64) (result string, err error) {
	vma_ms := vma * 1000 / 3600
	vma_100 := 100 / vma_ms
	calcul := vma_100 / percent * distance

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
func calcul_vma_speed(vma, percent float64) float64 {
	return (vma * percent) / 100
}

// calcul_pace from a speed (kmh)
func calcul_pace(vitesse float64) (ret string) {
	var e = 1 / vitesse * 60
	var t = math.Floor(e / 60)
	var n = math.Floor(e - t*60)
	var r = round(60 * (e - t*60 - n))

	if r == 60 {
		n += 1
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

func get_vmas(value string) (vmas []int) {
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

func GenerateProgram(workout Workout, target_vma string) (ts TemplateStruct, err error) {
	var total_time, time_lap string
	vmas := map[string]WorkoutVMA{}

	meters, err := strconv.ParseFloat(workout.Meters, 64)
	if err != nil {
		return
	}

	track_length, err := strconv.ParseFloat(strconv.Itoa(TRACK_LENGTH), 64)
	if err != nil {
		return
	}

	percentage, err := strconv.ParseFloat(workout.Percentage, 64)
	if err != nil {
		return
	}

	track_laps := meters / track_length
	laps := fmt.Sprintf("%.1f", track_laps)
	if strings.HasSuffix(laps, ".0") {
		laps = strings.TrimSuffix(laps, ".0")
	} else if laps == "0.5" {
		laps = "½"
	} else if strings.HasSuffix(laps, ".5") {
		laps = strings.Replace(laps, ".5", "½", -1)
	}
	workout.TrackLaps = laps
	workout.TrackLength = TRACK_LENGTH

	for _, vmad := range get_vmas(target_vma) {
		workout_vma := WorkoutVMA{}
		vma := float64(vmad)
		total_time, err = calcul_vma_distance(vma, percentage, meters)
		if err != nil {
			return
		}
		workout_vma.VMA = fmt.Sprintf("%.f", vma)
		workout_vma.TotalTime = total_time
		if int(meters) >= TRACK_LENGTH {
			time_lap, err = calcul_vma_distance(vma, percentage, float64(TRACK_LENGTH))
			if err != nil {
				return
			}
			workout_vma.TimeTrack = time_lap
		} else {
			workout_vma.TimeTrack = "NA"
		}
		speed := calcul_vma_speed(vma, percentage)
		workout_vma.Speed = fmt.Sprintf("%.2f", speed)
		workout_vma.Pace = calcul_pace(speed)

		vmas[workout_vma.VMA] = workout_vma
	}

	ts = TemplateStruct{
		VMAs:    vmas,
		Workout: workout,
	}

	return
}
