package chmoufrack

import (
	"bytes"
	"fmt"
	"log"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
)

func generate_content(ts TemplateStruct, content *bytes.Buffer) (err error) {
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

func generate_template(program_name, content string, outputWriter *bytes.Buffer) (err error) {
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

func generate_all_vmas(meters, percentage float64) (vmas map[string]WorkoutVMA, err error) {
	var total_time, time_lap string
	vmas = map[string]WorkoutVMA{}

	for _, vmad := range getVMAS(VMA) {
		wt := WorkoutVMA{}
		vma := float64(vmad)
		total_time, err = calcul_vma_distance(vma, percentage, meters)
		if err != nil {
			return
		}
		wt.VMA = fmt.Sprintf("%.f", vma)
		wt.TotalTime = total_time
		if int(meters) >= TRACK_LENGTH {
			time_lap, err = calcul_vma_distance(vma, percentage, float64(TRACK_LENGTH))
			if err != nil {
				return
			}
			wt.TimeTrack = time_lap
		} else {
			wt.TimeTrack = "NA"
		}
		speed := calcul_vma_speed(vma, percentage)
		wt.Speed = fmt.Sprintf("%.2f", speed)
		wt.Pace = calcul_pace(speed)

		vmas[wt.VMA] = wt
	}
	return
}

func generate_workout(w Workout) (wr Workout, err error) {
	wr = w
	meters, err := strconv.ParseFloat(w.Meters, 64)
	if err != nil {
		return
	}

	track_length, err := strconv.ParseFloat(strconv.Itoa(TRACK_LENGTH), 64)
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
	wr.TrackLaps = laps
	wr.TrackLength = TRACK_LENGTH
	return
}

func HTMLGen(program_name string, rounds []Workout, outputWriter *bytes.Buffer) error {
	var content bytes.Buffer

	for _, workout := range rounds {
		genw, err := generate_workout(workout)
		if err != nil {
			return err
		}

		meters, _ := strconv.ParseFloat(genw.Meters, 64)
		percentage, _ := strconv.ParseFloat(genw.Percentage, 64)

		vmas, err := generate_all_vmas(meters, percentage)
		if err != nil {
			log.Fatal(err)
		}

		ts := TemplateStruct{
			VMAs: vmas,
			WP:   genw,
		}
		err = generate_content(ts, &content)
		if err != nil {
			return err
		}
	}

	err := generate_template(program_name, content.String(), outputWriter)
	if err != nil {
		return err
	}
	return nil
}
