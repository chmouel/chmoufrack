package chmoufrack

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
)

func html_content(ts TemplateStruct, content *bytes.Buffer) (err error) {
	t, err := template.ParseFiles(filepath.Join(STATIC_DIR, "templates", "content.tmpl"))
	if err != nil {
		return
	}
	err = t.Execute(content, ts)
	if err != nil {
		return
	}
	return
}

func html_main_template(program_name, content string, outputWriter *bytes.Buffer) (err error) {
	dico := map[string]string{
		"Content":     content,
		"ProgramName": program_name,
	}

	t, err := template.ParseFiles(filepath.Join(STATIC_DIR, "templates", "template.tmpl"))
	if err != nil {
		return
	}
	err = t.Execute(outputWriter, dico)
	if err != nil {
		return
	}
	return
}

func getVMAS(value string) (vmas []int) {
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

func generate_program(workout Workout) (ts TemplateStruct, err error) {
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

	for _, vmad := range getVMAS(TARGET_VMA) {
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
		VMAs: vmas,
		WP:   workout,
	}

	return
}

func HTMLGen(program_name string, rounds []Workout, outputWriter *bytes.Buffer) (err error) {
	var content bytes.Buffer
	var ts TemplateStruct

	for _, workout := range rounds {
		if ts, err = generate_program(workout); err != nil {
			return
		}
		if err = html_content(ts, &content); err != nil {
			return
		}
	}

	err = html_main_template(program_name, content.String(), outputWriter)
	return
}
